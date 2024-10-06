package bybit

import (
	"algotrading_v1/notifier"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

var StopTradeCh = make(chan bool)

// Запуск торгового потока с покупками и продажами каждые 5 секунд
func StartTradeStream(n notifier.Notifier) {
	var mu sync.Mutex
	for {
		select {
		case <-StopTradeCh: // Добавляем проверку остановки
			log.Println("Received stop signal, stopping trade stream...")
			return
		default:
			wsURL := "wss://stream.bybit.com/v5/public/spot"
			conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
			if err != nil {
				log.Printf("Error connecting to WebSocket: %v. Retrying...", err)
				time.Sleep(5 * time.Second)
				continue
			}
			log.Println("Connected to WebSocket")

			defer conn.Close()

			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					log.Printf("Error reading message: %v. Reconnecting...", err)
					break
				}

				var response TradeResponse
				err = json.Unmarshal(message, &response)
				if err != nil {
					log.Printf("Error parsing message: %v", err)
					continue
				}

				if response.Topic == "publicTrade.BTCUSDT" {
					for _, tradeData := range response.Data {
						lastPrice := tradeData.Price
						mu.Lock()
						performTrade(n, tradeData.TradeID, lastPrice)
						mu.Unlock()
						time.Sleep(5 * time.Second)

						// Проверяем сигнал остановки
						select {
						case <-StopTradeCh:
							log.Println("Received stop signal during trade, stopping...")
							return
						default:
						}
					}
				}
			}
		}
	}
}

// Выполнение сделки
func performTrade(n notifier.Notifier, tradeID, lastPrice string) {
	quantity := "0.1" // Пример объёма

	// Покупка
	buyTime := time.Now().Format("2006-01-02 15:04:05")
	n.SendTradeResult(tradeID, "buy", lastPrice, quantity, buyTime)

	time.Sleep(5 * time.Second) // Ждём 5 секунд перед продажей

	// Продажа
	sellTime := time.Now().Format("2006-01-02 15:04:05")
	n.SendTradeResult(tradeID, "sell", lastPrice, quantity, sellTime)
}
