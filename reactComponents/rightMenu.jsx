"use strict";

function back() {
  history.back();
}

function rightMenu() {
  return (
    <header class="header" role="banner">
      <div class="nav-wrap">
        <nav class="main-nav" role="navigation">
          <ul class="unstyled list-hover-slide">
            <li>
              <a style={{ cursor:"pointer" }} onClick={() => back()}>Назад</a>
            </li>
          </ul>
        </nav>
      </div>
    </header>
  );
}

let domContainer = document.querySelector("#rightMenu");
ReactDOM.render(rightMenu(), domContainer);
