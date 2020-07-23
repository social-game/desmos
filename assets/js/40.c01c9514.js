(window.webpackJsonp=window.webpackJsonp||[]).push([[40],{393:function(e,t,s){"use strict";s.r(t);var a=s(9),o=Object(a.a)({},(function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("ContentSlotsDistributor",{attrs:{"slot-key":e.$parent.slotKey}},[s("h1",{attrs:{id:"polldata"}},[s("a",{staticClass:"header-anchor",attrs:{href:"#polldata"}},[e._v("#")]),e._v(" PollData")]),e._v(" "),s("p",[e._v("The "),s("code",[e._v("PollData")]),e._v(" object is used to specify the details of a poll that should be associated to a post. Please note that it is "),s("strong",[e._v("not")]),e._v(" necessary to associate a poll to each post. Instead, if you want to create a "),s("RouterLink",{attrs:{to:"/types/posts/post.html"}},[s("code",[e._v("Post")])]),e._v(" without any poll associated to it, you simply have to use the "),s("code",[e._v("nil")]),e._v(" value for this field.")],1),e._v(" "),s("p",[e._v("Following you will find a description for all the contained field of the "),s("code",[e._v("PollData")]),e._v(" object.")]),e._v(" "),s("h2",{attrs:{id:"question"}},[s("a",{staticClass:"header-anchor",attrs:{href:"#question"}},[e._v("#")]),e._v(" "),s("code",[e._v("Question")])]),e._v(" "),s("p",[e._v("This field contains the question that should be associated with the poll. It currently has no checks associated to it a part from the non-empty check.")]),e._v(" "),s("h2",{attrs:{id:"providedanswers"}},[s("a",{staticClass:"header-anchor",attrs:{href:"#providedanswers"}},[e._v("#")]),e._v(" "),s("code",[e._v("ProvidedAnswers")])]),e._v(" "),s("p",[e._v("This field allows to specify a list of answers that are provided to the users willing to answer the  poll.")]),e._v(" "),s("p",[e._v("Each answer should be composed of two attributes:")]),e._v(" "),s("ul",[s("li",[s("code",[e._v("ID")]),e._v(", which identifies uniquely inside the answers' list that answer.")]),e._v(" "),s("li",[s("code",[e._v("Text")]),e._v(", which contains the text of the answer itself.")])]),e._v(" "),s("p",[e._v("The minimum number of answers that a poll must have is 2.")]),e._v(" "),s("h2",{attrs:{id:"enddate"}},[s("a",{staticClass:"header-anchor",attrs:{href:"#enddate"}},[e._v("#")]),e._v(" "),s("code",[e._v("EndDate")])]),e._v(" "),s("p",[e._v("The "),s("code",[e._v("EndDate")]),e._v(" field allows you to specify the date after which the poll should be considered closed and no longer accepting the answers.")]),e._v(" "),s("h2",{attrs:{id:"open"}},[s("a",{staticClass:"header-anchor",attrs:{href:"#open"}},[e._v("#")]),e._v(" "),s("code",[e._v("Open")])]),e._v(" "),s("p",[e._v("This field tells whether the poll is still open and accepting new answers from users or not. Please note that the default value for this field is "),s("code",[e._v("true")]),e._v(" and trying to create a poll with it set to "),s("code",[e._v("false")]),e._v(" will result in an error.")]),e._v(" "),s("p",[e._v("During the chain execution, this field will be automatically changed to "),s("code",[e._v("false")]),e._v(" when the "),s("code",[e._v("EndDate")]),e._v(" is passed.")]),e._v(" "),s("h2",{attrs:{id:"allowsmultipleanswers"}},[s("a",{staticClass:"header-anchor",attrs:{href:"#allowsmultipleanswers"}},[e._v("#")]),e._v(" "),s("code",[e._v("AllowsMultipleAnswers")])]),e._v(" "),s("p",[e._v("This field allows to specify whether or not the poll allows multiple answers from the same user. If set to "),s("code",[e._v("true")]),e._v(", the users will be able to specify more than one answer to the same poll. Otherwise, if set to "),s("code",[e._v("false")]),e._v(", each user will be allowed to answer with only one option to the poll.")]),e._v(" "),s("h2",{attrs:{id:"allowsansweredits"}},[s("a",{staticClass:"header-anchor",attrs:{href:"#allowsansweredits"}},[e._v("#")]),e._v(" "),s("code",[e._v("AllowsAnswerEdits")])]),e._v(" "),s("p",[e._v("By setting this field to "),s("code",[e._v("true")]),e._v(", you will allow users to change their mind while the poll is still open, allowing them to change their answer(s). If set to "),s("code",[e._v("false")]),e._v(", they will not be able to do so and their final answer(s) will be the first one they give.")])])}),[],!1,null,null,null);t.default=o.exports}}]);