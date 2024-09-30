package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

// Структура для подписки на канал сделок
type SubscribeMessage struct {
	Op   string   `json:"op"`
	Args []string `json:"args"`
}

// Структура для ответа с информацией о сделках
type TradeResponse struct {
	Topic string `json:"topic"`
	Data  []struct {
		Price     string `json:"p"` // Цена сделки
		Quantity  string `json:"v"` // Количество (объём сделки)
		Timestamp int64  `json:"T"` // Время сделки в миллисекундах
		TradeId   string `json:"i"` // Id сделки
	} `json:"data"`
}

func main() {
	// Подключение к WebSocket Bybit
	wsURL := "wss://stream.bybit.com/v5/public/spot"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket:", err)
	}
	defer conn.Close()

	// Подписываемся на последние сделки для пары BTC/USDT
	subscribeMsg := SubscribeMessage{
		Op:   "subscribe",
		Args: []string{"publicTrade.BTCUSDT"},
	}

	err = conn.WriteJSON(subscribeMsg)
	if err != nil {
		log.Fatal("Error subscribing to trades:", err)
	}

	fmt.Println("Subscribed to BTC/USDT trades...")

	// Канал для обработки сигналов завершения работы (Ctrl+C)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Чтение сообщений из WebSocket
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				return
			}

			//// Логируем все входящие сообщения для анализа
			//log.Println("Received message:", string(message))

			// Парсим сообщение
			var tradeResponse TradeResponse
			err = json.Unmarshal(message, &tradeResponse)
			if err != nil {
				log.Println("Unmarshal error:", err)
				continue
			}

			// Если пришли данные о сделках
			if tradeResponse.Topic == "publicTrade.BTCUSDT" {
				for _, trade := range tradeResponse.Data {
					// Время сделки в секундах
					tradeTime := time.UnixMilli(trade.Timestamp)
					fmt.Printf("TradeId: %s, Price: %s USDT, Quantity: %s BTC, Time: %s\n", trade.TradeId, trade.Price, trade.Quantity, tradeTime)
				}
			}
		}
	}()

	// Ожидание завершения
	for {
		select {
		case <-interrupt:
			fmt.Println("Interrupt received, shutting down...")
			// Отключаемся от WebSocket
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error closing WebSocket:", err)
				return
			}
			return
		}
	}
}
