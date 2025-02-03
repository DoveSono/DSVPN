document.getElementById("generateKeyBtn").addEventListener("click", function() {
    fetch("/generate/key", {
        method: "POST"
    })
        .then(response => response.json())
        .then(data => {
            if (data.key) {
                document.getElementById("keyDisplay").textContent = data.key;
                alert("Ключ успешно сгенерирован!");
                document.getElementById("saveKeyBtn").classList.remove("hidden");  // Кнопка сохранить
            }
        })
        .catch(error => console.error("Error generating key:", error));
});

document.getElementById("generateSBlockBtn").addEventListener("click", function() {
    fetch("/generate/sblock", {
        method: "POST"
    })
        .then(response => response.json())
        .then(data => {
            if (data.sBlock) {
                document.getElementById("sBlockDisplay").textContent = data.sBlock;
                alert("S-блок успешно сгенерирован!");
                document.getElementById("saveSBlockBtn").classList.remove("hidden");  // Кнопка сохранить
            }
        })
        .catch(error => console.error("Error generating S-block:", error));
});

// Загрузка ключа из файла
document.getElementById("uploadKeyBtn").addEventListener("click", function() {
    document.getElementById("uploadKeyInput").click();
});

document.getElementById("uploadKeyInput").addEventListener("change", function(e) {
    const file = e.target.files[0];
    const reader = new FileReader();
    reader.onload = function(event) {
        document.getElementById("keyDisplay").textContent = event.target.result;
        document.getElementById("saveKeyBtn").classList.remove("hidden");  // Кнопка сохранить
    };
    reader.readAsText(file);
});

// Загрузка S-блока из файла
document.getElementById("uploadSBlockInput").addEventListener("change", function(e) {
    const file = e.target.files[0];
    const reader = new FileReader();
    reader.onload = function(event) {
        document.getElementById("sBlockDisplay").textContent = event.target.result;
        document.getElementById("saveSBlockBtn").classList.remove("hidden");  // Кнопка сохранить
    };
    reader.readAsText(file);
});

// Сохранение ключа в файл
document.getElementById("saveKeyBtn").addEventListener("click", function() {
    const key = document.getElementById("keyDisplay").textContent;
    const blob = new Blob([key], { type: "text/plain;charset=utf-8" });
    saveAs(blob, "generated_key.txt");
});

// Сохранение S-блока в файл
document.getElementById("saveSBlockBtn").addEventListener("click", function() {
    const sBlock = document.getElementById("sBlockDisplay").textContent;
    const blob = new Blob([sBlock], { type: "text/plain;charset=utf-8" });
    saveAs(blob, "generated_sblock.txt");
});
