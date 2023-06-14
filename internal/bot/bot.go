package bot

import (
	"fmt"
	"modules/internal/db"
	tdmod "modules/internal/todoist"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/volyanyk/todoist"
	tdist "github.com/volyanyk/todoist"

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
					msg = msg + fmt.Sprintf("№ %d  %s \n", k+1, v.Content) // TODO Настроить форматирование вывода задач
				}

				bot.Send(
					tgbotapi.NewMessage(chatId, msg),
				)
			case "create_task": // Отладить данный блок, задача не создается
				bot.Send(
					tgbotapi.NewMessage(chatId, "Введите наименование задачи и ее содержимое, разделив их запятой:"),
				)

				// Регистрация обработчика для ввода данных от пользователя
				inputCh := make(chan tgbotapi.Message)
				go func() {
					for msg := range inputCh {
						// Обработка введенных данных
						taskData := strings.Split(msg.Text, ",")
						taskContent := strings.TrimSpace(taskData[0])
						taskDescription := strings.TrimSpace(taskData[1])

						// Создание задачи в todoist
						var req tdist.AddTaskRequest
						req.Content = taskContent
						req.Description = taskDescription
						tdClient.AddTask(req)

						bot.Send(
							tgbotapi.NewMessage(chatId, "Задача успешно создана в todoist!"),
						)

					}
				}()

				// Регистрация обработчика callback'ов для ввода данных от пользователя
				go func() {
					for callback := range callbackCh {
						if callback.Message.Chat.ID == chatId {
							inputCh <- *callback.Message
						}
					}
				}()
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

		// Если тип возвращаемого сообщения == text, начинаем проверку на команды
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
						tgbotapi.NewInlineKeyboardButtonData("Новая задача", "create_task"),
						tgbotapi.NewInlineKeyboardButtonData("Пользователи бота", "number_of_users"),
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
