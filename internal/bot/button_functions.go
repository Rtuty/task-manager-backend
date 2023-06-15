package bot

import (
	"fmt"
	"modules/internal/db"
	tdmod "modules/internal/todoist"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	tdist "github.com/volyanyk/todoist"
)

func tasksList(tdClient *tdist.Client, bot *tgbotapi.BotAPI, chatId int64) {
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
}

func projectsList(tdClient *tdist.Client, bot *tgbotapi.BotAPI, chatId int64) {
	projects, err := tdClient.GetProjects()
	if err != nil {
		bot.Send(
			tgbotapi.NewMessage(chatId, "Не удалось получить список проектов"),
		)
	}
	var mass []tgbotapi.InlineKeyboardButton

	for _, v := range *projects {
		mass = append(mass, tgbotapi.NewInlineKeyboardButtonData(v.Name, v.ID))
	}

	msg := tgbotapi.NewMessage(chatId, "Проекты")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(mass)

	bot.Send(msg)
}

func numberOfUsers(bot *tgbotapi.BotAPI, chatId int64) {
	dbCon, err := db.NewDbClient()
	if err != nil {
		bot.Send(
			tgbotapi.NewMessage(chatId, fmt.Sprintf("Не удалость подключиться к базе данных. Ошибка: %s", err)),
		)
	}

	num, err := db.GetNumberOfUsers(dbCon)
	if err != nil {
		bot.Send(
			tgbotapi.NewMessage(chatId, "Ошибка базы данных. Обратитесь к администратору"),
		)
	}

	bot.Send(
		tgbotapi.NewMessage(chatId, fmt.Sprintf("%d людей используют данного бота. Дата запроса: %s", num, time.Now().GoString())),
	)
}
