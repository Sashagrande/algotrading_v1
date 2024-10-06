package main

import (
	"algotrading_v1/database"
	"algotrading_v1/telegram"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-telegram/bot/models"
	"log"
	"os"
	"os/signal"
	"strconv"
)

func main() {
	// Подключение к базе данных PostgreSQL
	connStr := os.Getenv("DATABASE_URL")
	err := database.InitDB(connStr)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Чтение токена и URL Webhook из переменных окружения
	token := os.Getenv("TELEGRAM_TOKEN")
	webhookURL := os.Getenv("WEBHOOK_URL")
	chatIDStr := os.Getenv("CHAT_ID")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if token == "" || webhookURL == "" || chatID == 0 {
		log.Fatal("TELEGRAM_TOKEN or WEBHOOK_URL or CHAT_ID not set")
	}

	// Инициализация бота
	b, err := telegram.InitBot(token, webhookURL, chatID)
	if err != nil {
		fmt.Printf("Error initializing bot: %v\n", err)
		return
	}

	// Создание Gin-сервера для Webhook
	r := gin.Default()

	// Обрабатываем webhook
	r.POST("/webhook", func(c *gin.Context) {
		var update models.Update
		if err := c.ShouldBindJSON(&update); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		b.ProcessUpdate(context.Background(), &update)
		c.Status(200)
	})

	// Простой хендлер для корневого пути
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Webhook is active!",
		})
	})

	// Запуск сервера
	go func() {
		if err := r.Run(":6127"); err != nil {
			fmt.Printf("Failed to run server: %v\n", err)
		}
	}()

	// Ожидание сигнала завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	fmt.Println("Server and bot are running. Press Ctrl+C to exit...")
	<-sigChan

	fmt.Println("Shutting down...")
}
