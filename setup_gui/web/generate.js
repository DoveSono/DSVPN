document.addEventListener("DOMContentLoaded", function () {
    let keyGenerated = false;
    let sBlockGenerated = false;

    const continueBtn = document.createElement("button");
    continueBtn.textContent = "Продолжить";
    continueBtn.id = "continueBtn";
    continueBtn.classList.add("hidden", "continue-button");
    document.body.appendChild(continueBtn);

    // Проверка, можно ли показать кнопку "Продолжить"
    function checkContinue() {
        if (keyGenerated && sBlockGenerated) {
            continueBtn.classList.remove("hidden");
        } else {
            continueBtn.classList.add("hidden");
        }
    }

    // Генерация ключа
    document.getElementById("generateKeyBtn").addEventListener("click", function () {
        fetch("/generate/key", { method: "POST" })
            .then(response => response.json())
            .then(data => {
                if (data.key) {
                    const keyStatus = document.getElementById("keyStatus");
                    keyStatus.textContent = "Ключ успешно сгенерирован!";
                    keyStatus.classList.remove("hidden");
                    document.getElementById("saveKeyBtn").classList.remove("hidden");
                    keyGenerated = true;
                    checkContinue();
                }
            })
            .catch(error => console.error("Ошибка при генерации ключа:", error));
    });

    // Генерация S-блока
    document.getElementById("generateSBlockBtn").addEventListener("click", function () {
        fetch("/generate/sblock", { method: "POST" })
            .then(response => response.json())
            .then(data => {
                if (data.sBlock) {
                    const sBlockStatus = document.getElementById("sBlockStatus");
                    sBlockStatus.textContent = "S-блок успешно сгенерирован!";
                    sBlockStatus.classList.remove("hidden");
                    document.getElementById("saveSBlockBtn").classList.remove("hidden");
                    sBlockGenerated = true;
                    checkContinue();
                }
            })
            .catch(error => console.error("Ошибка при генерации S-блока:", error));
    });

    // Загрузка файла ключа
    document.getElementById("uploadKeyBtn").addEventListener("click", function() {
        document.getElementById("uploadKeyInput").click();
    });

    document.getElementById("uploadKeyInput").addEventListener("change", function (e) {
        const file = e.target.files[0];
        const reader = new FileReader();
        reader.onload = function (event) {
            const keyStatus = document.getElementById("keyStatus");
            keyStatus.textContent = "Ключ успешно загружен!";
            keyStatus.classList.remove("hidden");
            document.getElementById("saveKeyBtn").classList.remove("hidden");
            keyGenerated = true;
            checkContinue();
        };
        reader.readAsText(file);
    });

    // Загрузка файла S-блока
    document.getElementById("uploadSBlockBtn").addEventListener("click", function() {
        document.getElementById("uploadSBlockInput").click();
    });

    document.getElementById("uploadSBlockInput").addEventListener("change", function (e) {
        const file = e.target.files[0];
        const reader = new FileReader();
        reader.onload = function (event) {
            const sBlockStatus = document.getElementById("sBlockStatus");
            sBlockStatus.textContent = "S-блок успешно загружен!";
            sBlockStatus.classList.remove("hidden");
            document.getElementById("saveSBlockBtn").classList.remove("hidden");
            sBlockGenerated = true;
            checkContinue();
        };
        reader.readAsText(file);
    });

    // Сохранение ключа в файл
    document.getElementById("saveKeyBtn").addEventListener("click", function () {
        const key = "Ваш сгенерированный ключ"; // Поскольку ключ больше не отображается, это место будет только для сохранения
        const blob = new Blob([key], { type: "text/plain;charset=utf-8" });
        saveAs(blob, "generated_key.txt");
    });

    // Сохранение S-блока в файл
    document.getElementById("saveSBlockBtn").addEventListener("click", function () {
        const sBlock = "Ваш сгенерированный S-блок"; // Поскольку S-блок больше не отображается
        const blob = new Blob([sBlock], { type: "text/plain;charset=utf-8" });
        saveAs(blob, "generated_sblock.txt");
    });

    // Обработчик для кнопки "Продолжить"
    continueBtn.addEventListener("click", function () {
        window.location.href = "/setup";
    });
});
