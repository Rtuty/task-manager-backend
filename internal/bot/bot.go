package bot

import (
	"fmt"
	"modules/internal/db"
	tdmod "modules/internal/todoist"
	"os"
	"reflect"
	"time"

	"github.com/volyanyk/todoist"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func StartBot(tdClient *todoist.Client) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOTTOKEN"))
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
				bot.Send(
					tgbotapi.NewMessage(update.Message.Chat.ID, "Бот активирован и готов принимать Ваши поручения по менеджменту задач! :)"),
				)
			case "/number_of_users":
				num, err := db.GetNumberOfUsers()
				if err != nil {
					bot.Send(
						tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка базы данных. Обратитесь к администратору"),
					)
				}

				bot.Send(
					tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%d людей используют данного бота. Дата запроса: %s", num, time.Now().GoString())),
				)
			case "/tasks":
				tasks, err := tdmod.GetTasks(tdClient)
				if err != nil {
					bot.Send(
						tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("не удалось получить задачи из todoist, Ошибка: %s", err)),
					)
				}

				var msg string

				for k, v := range *tasks {
					msg = msg + fmt.Sprintf("№ %d \n %s \n ", k, v.Content) // TODO Настроить форматирование вывода задач
				}

				bot.Send(
					tgbotapi.NewMessage(update.Message.Chat.ID, msg),
				)
			}
		}
	}
}
