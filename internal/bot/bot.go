package bot

import (
	"os"
	"reflect"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func StartBot() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
	}

	u := tgbotapi.NewUpdate(0) // Время обновления
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Если тип возвращаемого сообщения == text, начинаем проверку на компанды
		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
			switch update.Message.Text {
			case "/start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Бот активирован и готов принимать Ваши поручения по менеджменту задач! :)")
				bot.Send(msg)
			}
		}
	}
}
