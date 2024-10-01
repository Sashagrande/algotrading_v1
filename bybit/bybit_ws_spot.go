package bybit

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

func StartSpotStream() {
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
	log.Println("Subscribed to BTC/USDT spot prices...")

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
			tradeTime := time.UnixMilli(trade.Timestamp)
			fmt.Printf("TradeId: %s, Price: %s USDT, Quantity: %s BTC, Time: %s\n",
				trade.TradeId, trade.Price, trade.Quantity, tradeTime)
		}
	}
}
