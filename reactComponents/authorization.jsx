"use strict";

function authorization() {
  return (
    <div>
      <a type="button" data-toggle="modal" data-target="#exampleModalCenter">
        Авторизация
      </a>
      <div
        class="modal fade"
        id="exampleModalCenter"
        tabindex="-1"
        role="dialog"
        aria-labelledby="exampleModalCenterTitle"
        aria-hidden="false"
      >
        <div class="modal-dialog modal-dialog-centered" role="document">
          <div class="modal-content" style={{ backgroundColor: "grey" }}>
            <div class="modal-header">
              <h5 class="modal-title" id="exampleModalLongTitle">
                Авторизация
              </h5>
              <button
                type="button"
                class="close"
                data-dismiss="modal"
                aria-label="Close"
              >
                <span aria-hidden="false">&times;</span>
              </button>
            </div>
            <div class="modal-body">
              <form action="/login" method="POST" id="authorization">
                <input
                  required
                  type="text"
                  name="user_name"
                  id="user_name"
                  placeholder="Введите имя пользователя"
                  class="form-control"
                />
                <br />
                <input
                  required
                  type="text"
                  name="password"
                  id="password"
                  placeholder="Введите пароль"
                  class="form-control"
                />
                <br />
                <button class="btn btn-outline-warning">Войти</button>
                <hr />У Вас еще нет аккаунта?
                <a href="/registration">Регистрация</a>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

let domContainer = document.querySelector("#authorization");
ReactDOM.render(authorization(), domContainer);
