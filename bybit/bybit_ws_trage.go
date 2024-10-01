package bybit

import (
	"algotrading_v1/database"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
)

func StartTradeStream() {
	wsURL := "wss://stream.bybit.com/v5/public/spot"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket:", err)
	}
	defer conn.Close()

	subscribeMsg := map[string]interface{}{
		"op":   "subscribe",
		"args": []string{"publicTrade.BTCUSDT"},
	}
	conn.WriteJSON(subscribeMsg)
	log.Println("Subscribed to BTC/USDT trade...")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		var tradeResponse TradeResponse
		err = json.Unmarshal(message, &tradeResponse)
		if err != nil {
			log.Println("Error parsing message:", err)
			continue
		}

		for _, trade := range tradeResponse.Data {
			priceThreshold := 63000.00
			price, err := strconv.ParseFloat(trade.Price, 64)
			if err != nil {
				log.Println("Error parsing price:", err)
				continue
			}

			if price <= priceThreshold {
				orderId, err := database.PlaceOrder(price, trade.Quantity, "buy")
				if err != nil {
					log.Println("Error placing order:", err)
					continue
				}
				fmt.Printf("Order placed: %s at %f USDT\n", orderId, price)
			}
		}
	}
}
