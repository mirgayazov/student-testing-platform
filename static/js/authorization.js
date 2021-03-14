"use strict";

function authorization() {
  return React.createElement(
    "div",
    null,
    React.createElement(
      "a",
      { type: "button", "data-toggle": "modal", "data-target": "#exampleModalCenter" },
      "\u0410\u0432\u0442\u043E\u0440\u0438\u0437\u0430\u0446\u0438\u044F"
    ),
    React.createElement(
      "div",
      {
        "class": "modal fade",
        id: "exampleModalCenter",
        tabindex: "-1",
        role: "dialog",
        "aria-labelledby": "exampleModalCenterTitle",
        "aria-hidden": "false"
      },
      React.createElement(
        "div",
        { "class": "modal-dialog modal-dialog-centered", role: "document" },
        React.createElement(
          "div",
          { "class": "modal-content", style: { backgroundColor: "grey" } },
          React.createElement(
            "div",
            { "class": "modal-header" },
            React.createElement(
              "h5",
              { "class": "modal-title", id: "exampleModalLongTitle" },
              "\u0410\u0432\u0442\u043E\u0440\u0438\u0437\u0430\u0446\u0438\u044F"
            ),
            React.createElement(
              "button",
              {
                type: "button",
                "class": "close",
                "data-dismiss": "modal",
                "aria-label": "Close"
              },
              React.createElement(
                "span",
                { "aria-hidden": "false" },
                "\xD7"
              )
            )
          ),
          React.createElement(
            "div",
            { "class": "modal-body" },
            React.createElement(
              "form",
              { action: "/login", method: "POST", id: "authorization" },
              React.createElement("input", {
                required: true,
                type: "text",
                name: "user_name",
                id: "user_name",
                placeholder: "\u0412\u0432\u0435\u0434\u0438\u0442\u0435 \u0438\u043C\u044F \u043F\u043E\u043B\u044C\u0437\u043E\u0432\u0430\u0442\u0435\u043B\u044F",
                "class": "form-control"
              }),
              React.createElement("br", null),
              React.createElement("input", {
                required: true,
                type: "text",
                name: "password",
                id: "password",
                placeholder: "\u0412\u0432\u0435\u0434\u0438\u0442\u0435 \u043F\u0430\u0440\u043E\u043B\u044C",
                "class": "form-control"
              }),
              React.createElement("br", null),
              React.createElement(
                "button",
                { "class": "btn btn-outline-warning" },
                "\u0412\u043E\u0439\u0442\u0438"
              ),
              React.createElement("hr", null),
              "\u0423 \u0412\u0430\u0441 \u0435\u0449\u0435 \u043D\u0435\u0442 \u0430\u043A\u043A\u0430\u0443\u043D\u0442\u0430?",
              React.createElement(
                "a",
                { href: "/registration" },
                "\u0420\u0435\u0433\u0438\u0441\u0442\u0440\u0430\u0446\u0438\u044F"
              )
            )
          )
        )
      )
    )
  );
}

var domContainer = document.querySelector("#authorization");
ReactDOM.render(authorization(), domContainer);