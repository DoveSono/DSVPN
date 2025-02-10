package main

import (
	"DSVPN/gost" // Импорт библиотеки для генерации ключей и S-блока
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

// Открытие браузера
func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("cmd", "/C", "start", url).Start()
	case "darwin": // macOS
		err = exec.Command("open", url).Start()
	default:
		log.Println("Неизвестная ОС, не удается открыть браузер")
		return
	}
	if err != nil {
		log.Println("Ошибка при открытии браузера:", err)
	}
}

func main() {
	var wg sync.WaitGroup
	r := gin.Default()

	// Создание сессий
	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{
		MaxAge:   1800, // 30 минут
		HttpOnly: true,
	})
	r.Use(sessions.Sessions("mysession", store))

	// Подключение статики и HTML
	r.Static("/web", "./web")
	r.LoadHTMLGlob("web/*")

	// Главная страница
	r.GET("/", func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("authenticated") != nil {
			c.Redirect(http.StatusFound, "/generate")
			return
		}
		c.HTML(http.StatusOK, "authorization.html", nil)
	})

	// Авторизация
	r.POST("/login", func(c *gin.Context) {
		var form struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if form.Username == "admin" && form.Password == "admindsvpn" {
			session := sessions.Default(c)
			session.Set("authenticated", true)
			session.Save()
			c.JSON(http.StatusOK, gin.H{"message": "Авторизация успешна!"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		}
	})

	// Страница генерации ключа
	r.GET("/generate", func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("authenticated") == nil {
			c.Redirect(http.StatusFound, "/")
			return
		}
		c.HTML(http.StatusOK, "generate.html", gin.H{
			"key":    "",
			"sBlock": "",
		})
	})

	// Генерация ключа и S-блока
	r.POST("/generate/key", func(c *gin.Context) {
		key, err := gost.GenerateGostKey()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"key": key})
	})

	r.POST("/generate/sblock", func(c *gin.Context) {
		sBlock := gost.GenerateGostSBlock()
		c.JSON(http.StatusOK, gin.H{"sBlock": sBlock})
	})

	// Страница настройки VPN после генерации
	r.GET("/setup", func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("authenticated") == nil {
			c.Redirect(http.StatusFound, "/")
			return
		}
		c.HTML(http.StatusOK, "servsetup.html", nil)
	})

	// Запуск сервера
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := r.Run(":8080")
		if err != nil {
			log.Fatal("Ошибка запуска сервера:", err)
		}
	}()

	// Открытие браузера через 2 секунды
	time.Sleep(2 * time.Second)
	openBrowser("http://localhost:8080")

	wg.Wait()
}
