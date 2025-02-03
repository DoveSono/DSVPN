package main

import (
	"DSVPN/gost" // Модуль для генерации ключа и S-блока
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
	case "darwin": // для macOS
		err = exec.Command("open", url).Start()
	default:
		log.Println("Неизвестная операционная система, не удается открыть браузер")
		return
	}
	if err != nil {
		log.Println("Не удалось открыть браузер:", err)
	}
}

func main() {
	var wg sync.WaitGroup
	r := gin.Default()

	// Создаем сессию с использованием cookie
	store := cookie.NewStore([]byte("secret"))
	// Время жизни сессии 30 минут
	store.Options(sessions.Options{
		MaxAge:   1800, // 30 минут в секундах
		HttpOnly: true,
	})
	r.Use(sessions.Sessions("mysession", store))

	// Загрузка статических файлов
	r.Static("/web", "./web")
	r.LoadHTMLGlob("web/*")

	// Директ для главной страницы
	r.GET("/", func(c *gin.Context) {
		// Проверка авторизации
		session := sessions.Default(c)
		if session.Get("authenticated") != nil {
			// Если авторизован -  редирект на страницу генерации
			c.Redirect(http.StatusFound, "/generate")
			return
		}

		// Если не авторизован - редирект на страницу авторизации
		c.HTML(http.StatusOK, "authorization.html", nil)
	})

	// Директ для авторизации
	r.POST("/login", func(c *gin.Context) {
		var form struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		// Данные из формы
		if err := c.ShouldBindJSON(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Проверка логина и пароля
		if form.Username == "admin" && form.Password == "admindsvpn" {
			// Сохранение сессии
			session := sessions.Default(c)
			session.Set("authenticated", true)
			session.Save()

			c.JSON(http.StatusOK, gin.H{"message": "Авторизация успешна!"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		}
	})

	// Директ для генерации ключа и S-блока
	r.GET("/generate", func(c *gin.Context) {
		// Проверка авторизации
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

	// Генерация ключа
	r.POST("/generate/key", func(c *gin.Context) {
		key, err := gost.GenerateGostKey()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"key": key})
	})

	// Генерация S-блока
	r.POST("/generate/sblock", func(c *gin.Context) {
		sBlock := gost.GenerateGostSBlock()
		c.JSON(http.StatusOK, gin.H{"sBlock": sBlock})
	})

	// Добавление горутины с сервером в wait-группу
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := r.Run(":8080")
		if err != nil {
			log.Fatal("Ошибка запуска сервера: ", err)
		}
	}()

	// Задержка, чтобы не было траблов
	time.Sleep(2 * time.Second)

	// Открытие сайта
	openBrowser("http://localhost:8080")

	// Ожидание завершения горутины
	wg.Wait()
}
