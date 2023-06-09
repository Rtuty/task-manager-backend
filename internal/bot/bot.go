package bot

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func initBotToken() (string, error) {
	token, exists := os.LookupEnv("BOTOKEN")
	if !exists {
		return "", errors.New("bot token is not found")
	}

	return token, nil
}

func StartBotInstance(ctx context.Context) {
	token, err := initBotToken()
	if err != nil {
		panic(err)
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
		bot.WithCallbackQueryDataHandler("button", bot.MatchTypePrefix, callbackHandler),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		panic(err)
	}

	user, _ := b.GetMe(context.Background())

	fmt.Printf("User: %#v\n", user)

	b.Start(ctx)
}

func callbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Chat.ID,
		Text:   "Выбран пункт: " + update.CallbackQuery.Data,
	})

}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Список задач", CallbackData: "tasks"},
			},
		},
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Привет, дружок, вот мой текущий функционал:",
		ReplyMarkup: kb,
	})
}
