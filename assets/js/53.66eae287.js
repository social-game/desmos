(window.webpackJsonp=window.webpackJsonp||[]).push([[53],{406:function(e,t,s){"use strict";s.r(t);var a=s(9),n=Object(a.a)({},(function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("ContentSlotsDistributor",{attrs:{"slot-key":e.$parent.slotKey}},[s("h2",{attrs:{id:"validator-security"}},[s("a",{staticClass:"header-anchor",attrs:{href:"#validator-security"}},[e._v("#")]),e._v(" Validator Security")]),e._v(" "),s("p",[e._v("Each validator candidate is encouraged to run its operations independently, as diverse setups increase the resilience of the network.")]),e._v(" "),s("h2",{attrs:{id:"key-management-hsm"}},[s("a",{staticClass:"header-anchor",attrs:{href:"#key-management-hsm"}},[e._v("#")]),e._v(" Key Management - HSM")]),e._v(" "),s("p",[e._v("It is mission critical that an attacker cannot steal a validator's key. If this is possible, it puts the entire stake delegated to the compromised validator at risk. Hardware security modules are an important strategy for mitigating this risk.")]),e._v(" "),s("p",[e._v("HSM modules must support "),s("code",[e._v("ed25519")]),e._v(" signatures for the hub. The YubiHSM2 supports "),s("code",[e._v("ed25519")]),e._v(" and can protect a private key but cannot ensure in a secure setting that it won't sign the same block twice.")]),e._v(" "),s("p",[e._v("We will update this page when more key storage solutions become available.")]),e._v(" "),s("h2",{attrs:{id:"sentry-nodes-ddos-protection"}},[s("a",{staticClass:"header-anchor",attrs:{href:"#sentry-nodes-ddos-protection"}},[e._v("#")]),e._v(" Sentry Nodes (DDOS Protection)")]),e._v(" "),s("p",[e._v("Validators are responsible for ensuring that the network can sustain denial of service attacks.")]),e._v(" "),s("p",[e._v("One recommended way to mitigate these risks is for validators to carefully structure their network topology in a so-called sentry node architecture.")]),e._v(" "),s("p",[e._v("Validator nodes should only connect to full-nodes they trust because they operate them themselves or are run by other validators they know socially. A validator node will typically run in a data center. Most data centers provide direct links to the networks of major cloud providers. The validator can use those links to connect to sentry nodes in the cloud. This shifts the burden of denial-of-service from the validator's node directly to its sentry nodes, and may require new sentry nodes be spun up or activated to mitigate attacks on existing ones.")]),e._v(" "),s("p",[e._v("Sentry nodes can be quickly spun up or change their IP addresses. Because the links to the sentry nodes are in private IP space, an internet based attacked cannot disturb them directly. This will ensure validator block proposals and votes always make it to the rest of the network.")]),e._v(" "),s("p",[e._v("We suggest sentry nodes to be set up on multiple cloud providers across different regions. A validator may be offline if the connected sentry nodes are all offline due to the outage of a cloud provider in a specific region.")]),e._v(" "),s("p",[e._v("To setup your sentry node architecture you can follow the instructions below:")]),e._v(" "),s("p",[e._v("Validator Nodes should edit their config.toml:")]),e._v(" "),s("div",{staticClass:"language-bash line-numbers-mode"},[s("pre",{pre:!0,attrs:{class:"language-bash"}},[s("code",[s("span",{pre:!0,attrs:{class:"token comment"}},[e._v("# Comma separated list of nodes to keep persistent connections to")]),e._v("\n"),s("span",{pre:!0,attrs:{class:"token comment"}},[e._v("# Do not add private peers to this list if you don't want them advertised")]),e._v("\npersistent_peers "),s("span",{pre:!0,attrs:{class:"token operator"}},[e._v("=")]),s("span",{pre:!0,attrs:{class:"token punctuation"}},[e._v("[")]),e._v("list of sentry nodes"),s("span",{pre:!0,attrs:{class:"token punctuation"}},[e._v("]")]),e._v("\n\n"),s("span",{pre:!0,attrs:{class:"token comment"}},[e._v("# Set true to enable the peer-exchange reactor")]),e._v("\npex "),s("span",{pre:!0,attrs:{class:"token operator"}},[e._v("=")]),e._v(" "),s("span",{pre:!0,attrs:{class:"token boolean"}},[e._v("false")]),e._v("\n")])]),e._v(" "),s("div",{staticClass:"line-numbers-wrapper"},[s("span",{staticClass:"line-number"},[e._v("1")]),s("br"),s("span",{staticClass:"line-number"},[e._v("2")]),s("br"),s("span",{staticClass:"line-number"},[e._v("3")]),s("br"),s("span",{staticClass:"line-number"},[e._v("4")]),s("br"),s("span",{staticClass:"line-number"},[e._v("5")]),s("br"),s("span",{staticClass:"line-number"},[e._v("6")]),s("br")])]),s("p",[e._v("Sentry Nodes should edit their config.toml:")]),e._v(" "),s("div",{staticClass:"language-bash line-numbers-mode"},[s("pre",{pre:!0,attrs:{class:"language-bash"}},[s("code",[s("span",{pre:!0,attrs:{class:"token comment"}},[e._v("# Comma separated list of peer IDs to keep private (will not be gossiped to other peers)")]),e._v("\n"),s("span",{pre:!0,attrs:{class:"token comment"}},[e._v("# Example ID: 3e16af0cead27979e1fc3dac57d03df3c7a77acc@3.87.179.235:26656")]),e._v("\n\nprivate_peer_ids "),s("span",{pre:!0,attrs:{class:"token operator"}},[e._v("=")]),e._v(" "),s("span",{pre:!0,attrs:{class:"token string"}},[e._v('"node_ids_of_private_peers"')]),e._v("\n")])]),e._v(" "),s("div",{staticClass:"line-numbers-wrapper"},[s("span",{staticClass:"line-number"},[e._v("1")]),s("br"),s("span",{staticClass:"line-number"},[e._v("2")]),s("br"),s("span",{staticClass:"line-number"},[e._v("3")]),s("br"),s("span",{staticClass:"line-number"},[e._v("4")]),s("br")])])])}),[],!1,null,null,null);t.default=n.exports}}]);