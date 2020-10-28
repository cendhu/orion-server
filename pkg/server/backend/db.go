package backend

import (
	"encoding/pem"
	"io/ioutil"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.ibm.com/blockchaindb/library/pkg/logger"
	"github.ibm.com/blockchaindb/protos/types"
	"github.ibm.com/blockchaindb/server/config"
	"github.ibm.com/blockchaindb/server/pkg/blockstore"
	"github.ibm.com/blockchaindb/server/pkg/fileops"
	"github.ibm.com/blockchaindb/server/pkg/worldstate"
	"github.ibm.com/blockchaindb/server/pkg/worldstate/leveldb"
)

//go:generate mockery --dir . --name DB --case underscore --output mocks/

// DB encapsulates functionality required to operate with database state
type DB interface {
	// LedgerHeight returns current height of the ledger
	LedgerHeight() (uint64, error)

	// DoesUserExist checks whenever user with given userID exists
	DoesUserExist(userID string) (bool, error)

	// GetUser retrieves user' record
	GetUser(querierUserID, targetUserID string) (*types.GetUserResponseEnvelope, error)

	// GetConfig returns database configuration
	GetConfig() (*types.GetConfigResponseEnvelope, error)

	// GetDBStatus returns status for database, checks whenever database was created
	GetDBStatus(dbName string) (*types.GetDBStatusResponseEnvelope, error)

	// GetData retrieves values for given key
	GetData(dbName, querierUserID, key string) (*types.GetDataResponseEnvelope, error)

	// SubmitTransaction submits transaction to the database
	SubmitTransaction(tx interface{}) error

	// BootstrapDB given bootstrap configuration initialize database by
	// creating required system tables to include database meta data
	BootstrapDB(conf *config.Configurations) error

	// IsReady returns true once instance of the DB is properly initiated, meaning
	// all system tables was created successfully
	IsReady() (bool, error)

	// Close frees and closes resources allocated by database instance
	Close() error
}

type db struct {
	queryProcessor *queryProcessor
	txProcessor    *transactionProcessor
	logger         *logger.SugarLogger
}

func NewDB(conf *config.Configurations, logger *logger.SugarLogger) (*db, error) {
	if conf.Node.Database.Name != "leveldb" {
		return nil, errors.New("only leveldb is supported as the state database")
	}

	ledgerDir := conf.Node.Database.LedgerDirectory
	if err := createLedgerDir(ledgerDir); err != nil {
		return nil, err
	}

	levelDB, err := leveldb.Open(
		&leveldb.Config{
			DBRootDir: constructWorldStatePath(ledgerDir),
			Logger:    logger,
		},
	)
	if err != nil {
		return nil, errors.WithMessage(err, "error while creating the world state database")
	}

	blockStore, err := blockstore.Open(
		&blockstore.Config{
			StoreDir: constructBlockStorePath(ledgerDir),
			Logger:   logger,
		},
	)
	if err != nil {
		return nil, errors.WithMessage(err, "error while creating the block store")
	}

	blkHeight, err := blockStore.Height()
	if err != nil {
		return nil, errors.WithMessage(err, "error while fetching the block store height")
	}

	qProcConfig := &queryProcessorConfig{
		nodeID:     conf.Node.Identity.ID,
		db:         levelDB,
		blockStore: blockStore,
		logger:     logger,
	}

	txProcConf := &txProcessorConfig{
		db:                 levelDB,
		blockStore:         blockStore,
		blockHeight:        blkHeight,
		txQueueLength:      conf.Node.QueueLength.Transaction,
		txBatchQueueLength: conf.Node.QueueLength.ReorderedTransactionBatch,
		blockQueueLength:   conf.Node.QueueLength.Block,
		maxTxCountPerBatch: conf.Consensus.MaxTransactionCountPerBlock,
		batchTimeout:       conf.Consensus.BlockTimeout,
		logger:             logger,
	}

	txProcessor, err := newTransactionProcessor(txProcConf)
	if err != nil {
		return nil, errors.WithMessage(err, "can't initiate tx processor")
	}

	return &db{
		queryProcessor: newQueryProcessor(qProcConfig),
		txProcessor:    txProcessor,
		logger:         logger,
	}, nil
}

// LedgerHeight returns ledger height
func (d *db) LedgerHeight() (uint64, error) {
	return d.queryProcessor.blockStore.Height()
}

// DoesUserExist checks whenever userID exists
func (d *db) DoesUserExist(userID string) (bool, error) {
	return d.queryProcessor.identityQuerier.DoesUserExist(userID)
}

// GetUser returns user's record
func (d *db) GetUser(querierUserID, targetUserID string) (*types.GetUserResponseEnvelope, error) {
	return d.queryProcessor.getUser(querierUserID, targetUserID)
}

// GetConfig returns database configuration
func (d *db) GetConfig() (*types.GetConfigResponseEnvelope, error) {
	return d.queryProcessor.getConfig()
}

// GetDBStatus returns database status
func (d *db) GetDBStatus(dbName string) (*types.GetDBStatusResponseEnvelope, error) {
	return d.queryProcessor.getDBStatus(dbName)
}

// SubmitTransaction submits transaction
func (d *db) SubmitTransaction(tx interface{}) error {
	return d.txProcessor.submitTransaction(tx)
}

// GetData returns value for provided key
func (d *db) GetData(dbName, querierUserID, key string) (*types.GetDataResponseEnvelope, error) {
	return d.queryProcessor.getData(dbName, querierUserID, key)
}

// Close closes and release resources used by db
func (d *db) Close() error {
	if err := d.queryProcessor.close(); err != nil {
		return err
	}

	return d.txProcessor.close()
}

// BootstrapDB bootstraps DB with system tables
func (d *db) BootstrapDB(conf *config.Configurations) error {
	configTx, err := prepareConfigTx(conf)
	if err != nil {
		return errors.Wrap(err, "failed to prepare and commit a configuration transaction")
	}

	if err := d.txProcessor.submitTransaction(configTx); err != nil {
		return errors.Wrap(err, "error while committing configuration transaction")
	}
	return nil
}

// IsReady returns true once instance of the DB is properly initiated, meaning
// all system tables was created successfully
func (d *db) IsReady() (bool, error) {
	for _, dbName := range worldstate.SystemDBs() {
		status, err := d.queryProcessor.getDBStatus(dbName)
		if err != nil {
			return false, err
		}

		if status.GetPayload().GetExist() {
			continue
		}
		return false, nil
	}

	return true, nil
}

type certsInGenesisConfig struct {
	nodeCert   []byte
	adminCert  []byte
	rootCACert []byte
}

func prepareConfigTx(conf *config.Configurations) (*types.ConfigTxEnvelope, error) {
	certs, err := readCerts(conf)
	if err != nil {
		return nil, err
	}

	clusterConfig := &types.ClusterConfig{
		Nodes: []*types.NodeConfig{
			{
				ID:          conf.Node.Identity.ID,
				Certificate: certs.nodeCert,
				Address:     conf.Node.Network.Address,
				Port:        conf.Node.Network.Port,
			},
		},
		Admins: []*types.Admin{
			{
				ID:          conf.Admin.ID,
				Certificate: certs.adminCert,
			},
		},
		RootCACertificate: certs.rootCACert,
	}

	return &types.ConfigTxEnvelope{
		Payload: &types.ConfigTx{
			TxID:      uuid.New().String(), // TODO: we need to change TxID to string
			NewConfig: clusterConfig,
		},
		// TODO: we can make the node itself sign the transaction
	}, nil
}

func readCerts(conf *config.Configurations) (*certsInGenesisConfig, error) {
	nodeCert, err := ioutil.ReadFile(conf.Node.Identity.CertificatePath)
	if err != nil {
		return nil, errors.Wrapf(err, "error while reading node certificate %s", conf.Node.Identity.CertificatePath)
	}
	nodePemCert, _ := pem.Decode(nodeCert)

	adminCert, err := ioutil.ReadFile(conf.Admin.CertificatePath)
	if err != nil {
		return nil, errors.Wrapf(err, "error while reading admin certificate %s", conf.Admin.CertificatePath)
	}
	adminPemCert, _ := pem.Decode(adminCert)

	rootCACert, err := ioutil.ReadFile(conf.RootCA.CertificatePath)
	if err != nil {
		return nil, errors.Wrapf(err, "error while reading rootCA certificate %s", conf.RootCA.CertificatePath)
	}
	rootCAPemCert, _ := pem.Decode(rootCACert)

	return &certsInGenesisConfig{
		nodeCert:   nodePemCert.Bytes,
		adminCert:  adminPemCert.Bytes,
		rootCACert: rootCAPemCert.Bytes,
	}, nil
}

func createLedgerDir(dir string) error {
	exist, err := fileops.Exists(dir)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}

	return fileops.CreateDir(dir)
}
