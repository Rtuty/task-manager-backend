package bot

import (
	"context"
	"errors"
	"fmt"
	"modules/internal/todoist"
	"os"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func initBotToken() (string, error) {
	token, exists := os.LookupEnv("TOKENTGBOT")
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

	client, err := todoist.NewClient()
	if err != nil {
		panic(err)
	}

	var result string
	switch update.CallbackQuery.Data {
	case "tasks":
		r, err := todoist.GetTasks(client)
		if err != nil {
			panic(err)
		}

		for _, v := range *r {
			result = result + " " + v.Content
		}
	}

	// client.GetActiveTasks()
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Chat.ID,
		Text:   result,
	})

}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Button 1", CallbackData: "tasks"},
			},
		},
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Click by button",
		ReplyMarkup: kb,
	})
}
