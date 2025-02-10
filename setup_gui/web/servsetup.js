document.addEventListener("DOMContentLoaded", function () {
    const form = document.getElementById("vpnSetupForm");

    form.addEventListener("submit", function (event) {
        event.preventDefault();

        const serverIp = document.getElementById("serverIp").value.trim();
        const protocol = document.getElementById("protocol").value;
        const port = document.getElementById("port").value.trim();
        const interfaceName = document.getElementById("interface").value.trim();
        const cipher = document.getElementById("cipher").value;

        if (!serverIp || !port || !interfaceName) {
            alert("Пожалуйста, заполните все обязательные поля.");
            return;
        }

        fetch("/vpn/start", {
            method: "POST",
            headers: { "Content-Type": "application/x-www-form-urlencoded" },
            body: new URLSearchParams({ serverIp, protocol, port, interface: interfaceName, cipher })
        })
            .then(response => response.json())
            .then(data => {
                showSuccessMessage(data.message);
            })
            .catch(error => {
                console.error("Ошибка:", error);
                showSuccessMessage("Произошла ошибка при запуске VPN. Проверьте настройки.");
            });
    });

    function showSuccessMessage(message) {
        const successMessage = document.createElement("div");
        successMessage.classList.add("success-message");
        successMessage.innerHTML = `
            <p>${message}</p>
            <button onclick="this.parentElement.remove()">Закрыть</button>
        `;
        document.body.appendChild(successMessage);

        const continueBtn = document.createElement("button");
        continueBtn.textContent = "Продолжить";
        continueBtn.id = "continueBtn";
        continueBtn.classList.remove("hidden");
        continueBtn.classList.add("continue-button");
        document.body.appendChild(continueBtn);

        continueBtn.addEventListener("click", function () {
            window.location.href = "/next-step";  // Можно заменить на нужную ссылку для продолжения
        });
    }

    const style = document.createElement("style");
    style.textContent = `
        .continue-button {
            position: fixed;
            bottom: 20px;
            left: 50%;
            transform: translateX(-50%);
            padding: 20px 40px;
            font-size: 20px;
            background-color: #4CAF50;
            color: white;
            border: none;
            border-radius: 8px;
            cursor: pointer;
            box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.1);
        }
        .continue-button:hover {
            background-color: #45a049;
        }
    `;
    document.head.appendChild(style);
});
