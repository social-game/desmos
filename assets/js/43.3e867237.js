(window.webpackJsonp=window.webpackJsonp||[]).push([[43],{390:function(e,t,s){"use strict";s.r(t);var a=s(42),r=Object(a.a)({},(function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("ContentSlotsDistributor",{attrs:{"slot-key":e.$parent.slotKey}},[s("h1",{attrs:{id:"automatic-full-node-setup"}},[s("a",{staticClass:"header-anchor",attrs:{href:"#automatic-full-node-setup"}},[e._v("#")]),e._v(" Automatic full node setup")]),e._v(" "),s("p",[e._v("Following you will find how to download and execute the script that allows you to run a Desmos full node in minutes.")]),e._v(" "),s("div",{staticClass:"custom-block warning"},[s("p",{staticClass:"custom-block-title"},[e._v("Requirements")]),e._v(" "),s("p",[e._v("Before proceeding, make sure you have read the "),s("RouterLink",{attrs:{to:"/fullnode/setup/overview.html"}},[e._v("overview page")]),e._v(" and your system satisfied the needed requirements.")],1)]),e._v(" "),s("h2",{attrs:{id:"_1-download-the-script"}},[s("a",{staticClass:"header-anchor",attrs:{href:"#_1-download-the-script"}},[e._v("#")]),e._v(" 1. Download the script")]),e._v(" "),s("p",[e._v("You can get the script by executing")]),e._v(" "),s("div",{staticClass:"language-shell line-numbers-mode"},[s("pre",{pre:!0,attrs:{class:"language-shell"}},[s("code",[s("span",{pre:!0,attrs:{class:"token function"}},[e._v("wget")]),e._v(" -O install-desmos-fullnode https://raw.githubusercontent.com/desmos-labs/desmos/master/contrib/validators/automatic-fullnode-installer.sh \n")])]),e._v(" "),s("div",{staticClass:"line-numbers-wrapper"},[s("span",{staticClass:"line-number"},[e._v("1")]),s("br")])]),s("p",[e._v("Once you downloaded it properly, you need to change its permissions making it executable:")]),e._v(" "),s("div",{staticClass:"language-shell line-numbers-mode"},[s("pre",{pre:!0,attrs:{class:"language-shell"}},[s("code",[s("span",{pre:!0,attrs:{class:"token function"}},[e._v("chmod")]),e._v(" +x ./install-desmos-fullnode\n")])]),e._v(" "),s("div",{staticClass:"line-numbers-wrapper"},[s("span",{staticClass:"line-number"},[e._v("1")]),s("br")])]),s("h2",{attrs:{id:"_2-execute-the-script"}},[s("a",{staticClass:"header-anchor",attrs:{href:"#_2-execute-the-script"}},[e._v("#")]),e._v(" 2. Execute the script")]),e._v(" "),s("p",[e._v("Once you got the script, you are now ready to use it.")]),e._v(" "),s("h3",{attrs:{id:"parameters"}},[s("a",{staticClass:"header-anchor",attrs:{href:"#parameters"}},[e._v("#")]),e._v(" Parameters")]),e._v(" "),s("p",[e._v("In order to work, it needs the following parameters:")]),e._v(" "),s("ol",[s("li",[e._v("The "),s("code",[e._v("moniker")]),e._v(" of the validator you are creating."),s("br"),e._v("\nThis is just a string that will allow you to identify the validator you are running locally. It can be anything you want.")])]),e._v(" "),s("h3",{attrs:{id:"running-the-script"}},[s("a",{staticClass:"header-anchor",attrs:{href:"#running-the-script"}},[e._v("#")]),e._v(" Running the script")]),e._v(" "),s("p",[e._v("Once you are ready to run the script, just execute:")]),e._v(" "),s("div",{staticClass:"language-shell line-numbers-mode"},[s("pre",{pre:!0,attrs:{class:"language-shell"}},[s("code",[e._v("./install-desmos-fullnode "),s("span",{pre:!0,attrs:{class:"token operator"}},[e._v("<")]),e._v("PARAMETERS"),s("span",{pre:!0,attrs:{class:"token operator"}},[e._v(">")]),e._v("\n")])]),e._v(" "),s("div",{staticClass:"line-numbers-wrapper"},[s("span",{staticClass:"line-number"},[e._v("1")]),s("br")])]),s("p",[e._v("E.g:")]),e._v(" "),s("div",{staticClass:"language- line-numbers-mode"},[s("pre",{pre:!0,attrs:{class:"language-text"}},[s("code",[e._v("./install-desmos-fullnode my-validator\n")])]),e._v(" "),s("div",{staticClass:"line-numbers-wrapper"},[s("span",{staticClass:"line-number"},[e._v("1")]),s("br")])]),s("h2",{attrs:{id:"how-it-works"}},[s("a",{staticClass:"header-anchor",attrs:{href:"#how-it-works"}},[e._v("#")]),e._v(" How it works")]),e._v(" "),s("p",[e._v("The script will perform the following operations.")]),e._v(" "),s("ol",[s("li",[s("p",[s("strong",[e._v("Environment setup")]),s("br"),e._v("\nIt will create all the necessary environmental variables.")])]),e._v(" "),s("li",[s("p",[s("strong",[e._v("Cosmovisor setup")]),s("br"),e._v("\nIt will download Cosmovisor and set it up so that your node is able to automatically update based on on-chain upgrades.")])]),e._v(" "),s("li",[s("p",[s("strong",[e._v("Desmos setup")]),s("br"),e._v("\nIt will download and install Desmos properly so that your node is able to start syncing and also update itself based on all the on-chain upgrades that have been done until now.")])]),e._v(" "),s("li",[s("p",[s("strong",[e._v("Service setup")]),s("br"),e._v("\nIt will setup a system service to make sure Desmos runs properly in the background.")])]),e._v(" "),s("li",[s("p",[s("strong",[e._v("Service start and log output")]),s("br"),e._v("\nFinally, it will start the system service and output the logs from it. You will see it syncing the blocks properly and catching up with the rest of the chain.")])])])])}),[],!1,null,null,null);t.default=r.exports}}]);