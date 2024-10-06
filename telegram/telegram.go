package telegram

import (
	"algotrading_v1/bybit"
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

type Bot struct {
	*bot.Bot
	ChatID int64
}

// var StopSpotCh = make(chan bool)

// Установка Webhook
func setWebhook(token, webhookURL string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook", token)
	data := fmt.Sprintf("url=%s/webhook", webhookURL)

	resp, err := http.Post(apiURL, "application/x-www-form-urlencoded", strings.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to set webhook: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	fmt.Printf("Webhook set response: %s\n", string(body))
	return nil
}

// Инициализация бота и установка Webhook с регистрацией команд
func InitBot(token, webhookURL string, chatID int64) (*Bot, error) {
	// Устанавливаем Webhook
	err := setWebhook(token, webhookURL)
	if err != nil {
		return nil, fmt.Errorf("error setting webhook: %v", err)
	}

	// Инициализация бота
	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}
	botInstance, err := bot.New(token, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to init bot: %v", err)
	}

	// Создаем экземпляр telegram.Bot
	tgBot := &Bot{
		Bot:    botInstance,
		ChatID: chatID,
	}

	// Регистрация команд
	tgBot.RegisterHandler(bot.HandlerTypeMessageText, "/stream_spot", bot.MatchTypeExact, tgBot.streamSpotHandler)
	tgBot.RegisterHandler(bot.HandlerTypeMessageText, "/stream_trade", bot.MatchTypeExact, tgBot.streamTradeHandler)
	tgBot.RegisterHandler(bot.HandlerTypeMessageText, "/stop", bot.MatchTypeExact, tgBot.stopHandler)

	return tgBot, nil
}

// Обработчики команд как методы структуры Bot
func (b *Bot) streamSpotHandler(ctx context.Context, botInstance *bot.Bot, update *models.Update) {
	bybit.StopSpotCh = make(chan bool)
	go bybit.StartSpotStream()
	botInstance.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Started streaming spot prices.",
	})
}

func (b *Bot) streamTradeHandler(ctx context.Context, botInstance *bot.Bot, update *models.Update) {
	bybit.StopTradeCh = make(chan bool)
	go bybit.StartTradeStream(b)
	botInstance.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Started trading BTC/USDT.",
	})
}

// Хендлер для команды /stop, чтобы остановить все операции
var stopTradeOnce sync.Once // Создаем объект Once для безопасного закрытия каналов

// Хендлер для команды /stop
func (b *Bot) stopHandler(ctx context.Context, botInstance *bot.Bot, update *models.Update) {
	stopTradeOnce.Do(func() {
		close(bybit.StopTradeCh) // Закрываем канал только один раз
		close(bybit.StopSpotCh)

		botInstance.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "All streams have been stopped.",
		})
		log.Println("Received stop command. Stopping all streams...")
	})
}

// Реализация метода SendTradeResult интерфейса Notifier
func (b *Bot) SendTradeResult(tradeID, status, price, quantity, time string) {
	message := fmt.Sprintf("TradeId: %s\nКол-во: %s\nСтатус: %s\nЦена: %s\nВремя сделки: %s",
		tradeID, quantity, status, price, time)

	b.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID:    b.ChatID,
		Text:      message,
		ParseMode: models.ParseModeMarkdown,
	})
}

// Хендлер по умолчанию
func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Available commands: /stream_spot, /stream_trade, /stop",
	})
}
