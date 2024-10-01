package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // PostgreSQL драйвер
)

var db *sql.DB

// Инициализация подключения к базе данных
func InitDB(connStr string) error {
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Проверим соединение с базой
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")
	return nil
}

// Функция для размещения ордера в таблице
func PlaceOrder(price float64, quantity string, orderType string) (int, error) {
	// Время выполнения сделки
	now := time.Now()

	// SQL-запрос для вставки данных
	query := `
		INSERT INTO trades (trade_id, buy_price, quantity, buy_time, status)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`

	var lastInsertID int
	err := db.QueryRow(query, generateTradeID(), price, quantity, now, orderType).Scan(&lastInsertID)
	if err != nil {
		log.Printf("Error inserting trade into database: %v", err)
		return 0, err
	}

	log.Printf("Order placed successfully with ID: %d", lastInsertID)
	return lastInsertID, nil
}

// Вспомогательная функция для генерации уникального ID сделки
func generateTradeID() string {
	return fmt.Sprintf("trade-%d", time.Now().UnixNano())
}
