package telegram

import (
	"algotrading_v1/bybit"
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"io/ioutil"
	"net/http"
	"strings"
)

// Установка Webhook
func setWebhook(token, webhookURL string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook", token)
	data := fmt.Sprintf("url=%s", webhookURL)

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

// Инициализация бота и установка Webhook
func InitBot(token, webhookURL string) (*bot.Bot, error) {
	// Устанавливаем Webhook
	err := setWebhook(token, webhookURL)
	if err != nil {
		return nil, fmt.Errorf("error setting webhook: %v", err)
	}

	// Инициализация бота
	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}
	b, err := bot.New(token, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to init bot: %v", err)
	}

	// Регистрация команд
	b.RegisterHandler(bot.HandlerTypeMessageText, "/stream_spot", bot.MatchTypeExact, streamSpotHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/stream_trade", bot.MatchTypeExact, streamTradeHandler)

	return b, nil
}

// Обработчики команд
func streamSpotHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	go bybit.StartSpotStream()
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Started streaming spot prices.",
	})
}

func streamTradeHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	go bybit.StartTradeStream()
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Started trading BTC/USDT.",
	})
}

// Хендлер по умолчанию
func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Available commands: /stream_spot, /stream_trade",
	})
}
