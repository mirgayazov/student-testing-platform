"use strict";

function back() {
  history.back();
}

function rightMenu() {
  return React.createElement(
    "header",
    { "class": "header", role: "banner" },
    React.createElement(
      "div",
      { "class": "nav-wrap" },
      React.createElement(
        "nav",
        { "class": "main-nav", role: "navigation" },
        React.createElement(
          "ul",
          { "class": "unstyled list-hover-slide" },
          React.createElement(
            "li",
            null,
            React.createElement(
              "a",
              { style: { cursor: "pointer" }, onClick: function onClick() {
                  return back();
                } },
              "\u041D\u0430\u0437\u0430\u0434"
            )
          )
        )
      )
    )
  );
}

var domContainer = document.querySelector("#rightMenu");
ReactDOM.render(rightMenu(), domContainer);