"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[4877],{3905:function(e,t,n){n.d(t,{Zo:function(){return l},kt:function(){return m}});var r=n(7294);function a(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function i(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function o(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?i(Object(n),!0).forEach((function(t){a(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):i(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function s(e,t){if(null==e)return{};var n,r,a=function(e,t){if(null==e)return{};var n,r,a={},i=Object.keys(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var c=r.createContext({}),u=function(e){var t=r.useContext(c),n=t;return e&&(n="function"==typeof e?e(t):o(o({},t),e)),n},l=function(e){var t=u(e.components);return r.createElement(c.Provider,{value:t},e.children)},d={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},p=r.forwardRef((function(e,t){var n=e.components,a=e.mdxType,i=e.originalType,c=e.parentName,l=s(e,["components","mdxType","originalType","parentName"]),p=u(n),m=a,b=p["".concat(c,".").concat(m)]||p[m]||d[m]||i;return n?r.createElement(b,o(o({ref:t},l),{},{components:n})):r.createElement(b,o({ref:t},l))}));function m(e,t){var n=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var i=n.length,o=new Array(i);o[0]=p;var s={};for(var c in t)hasOwnProperty.call(t,c)&&(s[c]=t[c]);s.originalType=e,s.mdxType="string"==typeof e?e:a,o[1]=s;for(var u=2;u<i;u++)o[u]=n[u];return r.createElement.apply(null,o)}return r.createElement.apply(null,n)}p.displayName="MDXCreateElement"},8532:function(e,t,n){n.r(t),n.d(t,{frontMatter:function(){return s},contentTitle:function(){return c},metadata:function(){return u},toc:function(){return l},default:function(){return p}});var r=n(7462),a=n(3366),i=(n(7294),n(3905)),o=["components"],s={id:"db",title:"Check the Existance of a Database"},c=void 0,u={unversionedId:"external/getting-started/queries/curl/db",id:"external/getting-started/queries/curl/db",isDocsHomePage:!1,title:"Check the Existance of a Database",description:"Checking the Database Existance",source:"@site/docs/external/getting-started/queries/curl/db.md",sourceDirName:"external/getting-started/queries/curl",slug:"/external/getting-started/queries/curl/db",permalink:"/orion-server/docs/external/getting-started/queries/curl/db",tags:[],version:"current",frontMatter:{id:"db",title:"Check the Existance of a Database"},sidebar:"Documentation",previous:{title:"Query an User Information",permalink:"/orion-server/docs/external/getting-started/queries/curl/user"},next:{title:"Query Data Using Keys",permalink:"/orion-server/docs/external/getting-started/queries/curl/simple-data-query"}},l=[{value:"Checking the Database Existance",id:"checking-the-database-existance",children:[],level:2}],d={toc:l};function p(e){var t=e.components,n=(0,a.Z)(e,o);return(0,i.kt)("wrapper",(0,r.Z)({},d,n,{components:t,mdxType:"MDXLayout"}),(0,i.kt)("h2",{id:"checking-the-database-existance"},"Checking the Database Existance"),(0,i.kt)("p",null,"To check whether a database exist/created, the user can issue a GET request on ",(0,i.kt)("inlineCode",{parentName:"p"},"/db/{dbname}")," endpoint where ",(0,i.kt)("inlineCode",{parentName:"p"},"{dbname}")," should be replaced with\nthe ",(0,i.kt)("inlineCode",{parentName:"p"},"dbname")," for which the user needs to perform this check."),(0,i.kt)("p",null,"For this query, the submitting user needs to sign ",(0,i.kt)("inlineCode",{parentName:"p"},'{"user_id":"<userid","db_name":"<dbname>"}')," where ",(0,i.kt)("inlineCode",{parentName:"p"},"userid")," denotes the submitting user and the\n",(0,i.kt)("inlineCode",{parentName:"p"},"<dbname>")," denotes the name of the database for which the user performs the existance check. "),(0,i.kt)("p",null,"When the BDB server bootups, it creates a default database called ",(0,i.kt)("inlineCode",{parentName:"p"},"bdb")," in the cluster. Hence, we can check its existance. For this case, the\nsubmitting user ",(0,i.kt)("inlineCode",{parentName:"p"},"admin")," needs to sign ",(0,i.kt)("inlineCode",{parentName:"p"},'{"user_id":"admin","db_name":"bdb"}')," as shown below:"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-sh"},'./bin/signer -privatekey=deployment/sample/crypto/admin/admin.key -data=\'{"user_id":"admin","db_name":"bdb"}\'\n')),(0,i.kt)("p",null,"The above command would produce a digital signature and prints it as base64 encoded string as shown below:"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-sh"},"MEUCIBzH0qIz88jKdHsJvmQsNNuK3Cf0G+7LDWSiwv6yjba0AiEAgb/hBFZrr3w64M0Q6LmZjQ0i/sjYr27K1DJSlXHWfRU=\n")),(0,i.kt)("p",null,"Once the signature is computed, we can issue a ",(0,i.kt)("inlineCode",{parentName:"p"},"GET")," request using the following ",(0,i.kt)("inlineCode",{parentName:"p"},"cURL")," command\nby setting the above signature string in the ",(0,i.kt)("inlineCode",{parentName:"p"},"Signature")," header."),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-sh"},'curl \\\n   -H "Content-Type: application/json" \\\n   -H "UserID: admin" \\\n   -H "Signature: abcd" \\\n   -X GET http://127.0.0.1:6001/db/bdb | jq .\n')),(0,i.kt)("p",null,"The above command results in the following output:"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-json"},'{\n  "response": {\n    "header": {\n      "node_id": "bdb-node-1"\n. \n\n},\n    "exist": true\n  },\n  "signature": "MEQCIAhSD5eQ+lBCaN7C/fILXcHADekGi+1RteDLmBbgHS4sAiAU+h/uwp/CrKmRdgHeAN7wOArRj5BdPC4Qp8Mzw4uIaQ=="\n}\n')))}p.isMDXComponent=!0}}]);