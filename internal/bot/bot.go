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
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKENTGBOT"))
	if err != nil {
		panic(err)
	}

	u := tgbotapi.NewUpdate(0) // Время обновления
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}

	// Создание канала для обработки callback'ов
	callbackCh := make(chan tgbotapi.CallbackQuery)

	// Регистрация обработчика callback'ов
	go func() {
		for callback := range callbackCh {
			// Получение chatId и data из callback'а
			chatId := callback.Message.Chat.ID
			data := callback.Data

			switch data {
			case "tasks":
				tasks, err := tdmod.GetTasks(tdClient)
				if err != nil {
					bot.Send(
						tgbotapi.NewMessage(chatId, fmt.Sprintf("не удалось получить задачи из todoist, Ошибка: %s", err)),
					)
				}

				var msg string

				for k, v := range *tasks {
					msg = msg + fmt.Sprintf("№ %d \n %s \n ", k, v.Content) // TODO Настроить форматирование вывода задач
				}

				bot.Send(
					tgbotapi.NewMessage(chatId, msg),
				)
			case "number_of_users":
				num, err := db.GetNumberOfUsers()
				if err != nil {
					bot.Send(
						tgbotapi.NewMessage(chatId, "Ошибка базы данных. Обратитесь к администратору"),
					)
				}

				bot.Send(
					tgbotapi.NewMessage(chatId, fmt.Sprintf("%d людей используют данного бота. Дата запроса: %s", num, time.Now().GoString())),
				)
			}
		}
	}()

	for update := range updates {
		// Обработка callback'ов
		if update.CallbackQuery != nil {
			callback := *update.CallbackQuery

			// Отправка callback-ответа
			callbackResponse := tgbotapi.NewCallback(callback.ID, callback.Data)
			_, err := bot.AnswerCallbackQuery(callbackResponse)
			if err != nil {
				panic(err)
			}

			callbackCh <- callback
		}
		if update.Message == nil {
			continue
		}

		// Если тип возвращаемого сообщения == text, начинаем проверку на компанды
		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
			var chatId int64 = update.Message.Chat.ID

			switch update.Message.Text {
			case "/start":
				bot.Send(
					tgbotapi.NewMessage(chatId, "Бот активирован и готов принимать Ваши поручения по менеджменту задач! :)"),
				)

				// Создание массива кнопок и добавление его в объект InlineKeyboardMarkup
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					[]tgbotapi.InlineKeyboardButton{
						tgbotapi.NewInlineKeyboardButtonData("Список задач", "tasks"),
						tgbotapi.NewInlineKeyboardButtonData("Получить пользователей бота", "number_of_users"),
					},
				)

				// Создание сообщения с кнопками
				msg := tgbotapi.NewMessage(chatId, "Выберите кнопку:")
				msg.ReplyMarkup = keyboard

				bot.Send(msg)
			}
		}

		// Обработка callback'ов
		if update.CallbackQuery != nil {
			callbackCh <- *update.CallbackQuery
		}
	}
}
