package bybit

// Общая структура для ответа по сделкам
type TradeResponse struct {
	Topic string `json:"topic"`
	Data  []struct {
		Price     string `json:"p"`
		Quantity  string `json:"v"`
		Timestamp int64  `json:"T"`
		TradeId   string `json:"i"`
	} `json:"data"`
}
