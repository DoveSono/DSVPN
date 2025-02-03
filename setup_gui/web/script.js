// Обработчик формы авторизации
document.getElementById("loginForm").addEventListener("submit", function(e) {
    e.preventDefault(); // предотвращаем отправку формы по умолчанию

    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;

    // Скрываем сообщение об ошибке перед отправкой запроса
    document.getElementById("error-message").classList.add("hidden");

    // Отправка запроса с логином и паролем
    fetch("/login", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ username: username, password: password })
    })
        .then(response => {
            // Лог ответа сервера
            console.log("Ответ от сервера: ", response);

            if (response.ok) {
                return response.json(); // Если запрос успешный
            } else {
                // Если статус не 200 OK
                return response.json().then(data => {
                    console.log("Ошибка от сервера:", data);
                    if (response.status === 401) {
                        document.getElementById("error-message").textContent = "Неверный логин или пароль";
                    } else {
                        document.getElementById("error-message").textContent = "Ошибка при авторизации";
                    }
                    document.getElementById("error-message").classList.remove("hidden");
                });
            }
        })
        .then(data => {
            if (data && data.message) {
                console.log("Авторизация успешна:", data);

                // Показываем всплывающее сообщение об успешной авторизации
                const successMessage = document.createElement("div");
                successMessage.classList.add("success-message");
                successMessage.innerHTML = `
                    <p>Авторизация успешна!</p>
                    <button id="ok-button">ОК</button>
                `;
                document.body.appendChild(successMessage);

                // Обработчик для кнопки "ОК"
                document.getElementById("ok-button").addEventListener("click", function() {
                    // Убираем сообщение
                    successMessage.remove();
                    // Перенаправляем на страницу генерации
                    window.location.href = "/generate";
                });
            }
        })
        .catch(error => {
            console.error("Ошибка при отправке данных:", error);
            document.getElementById("error-message").textContent = "Ошибка при отправке данных";
            document.getElementById("error-message").classList.remove("hidden");
        });
});

// Скрытие ошибки при повторном вводе
document.getElementById("username").addEventListener("input", function() {
    document.getElementById("error-message").classList.add("hidden");
});

document.getElementById("password").addEventListener("input", function() {
    document.getElementById("error-message").classList.add("hidden");
});
