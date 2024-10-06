package notifier

type Notifier interface {
	SendTradeResult(tradeID, status, price, quantity, time string)
}
