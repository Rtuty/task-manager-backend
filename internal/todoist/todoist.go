package todoist

import (
	"errors"
	"fmt"
	"os"

	"github.com/volyanyk/todoist"
)

func initTodoistToken() (string, error) {
	token, exists := os.LookupEnv("TODOIST")
	if !exists {
		return "", errors.New("todoist API token not found in .env")
	}

	return token, nil
}

func initFolderName() (string, error) {
	folderName, exists := os.LookupEnv("FOLDERNAME")
	if !exists {
		return "", errors.New("folder name not found in .env")
	}

	return folderName, nil
}

// Получаем токен, создаем и возвращаем новый клиент
func NewClient() (*todoist.Client, error) {
	token, err := initTodoistToken()
	if err != nil {
		return nil, err
	}

	client := todoist.New(token)

	return client, nil
}

func GetTasks(client *todoist.Client) (*[]todoist.Task, error) {
	// Получаем список проектов
	projects, err := client.GetProjects()
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка проектов. details: %v", err)
	}

	//  Находим id проекта (папки) по наименованию, которое указано в .env
	folderName, err := initFolderName()
	if err != nil {
		return nil, err
	}

	var folderId string

	for _, v := range *projects {
		if v.Name == folderName {
			folderId = v.ID
		}
	}

	// Получаем список задач из указанного проекта
	folderTasks, err := client.GetActiveTasks(todoist.GetActiveTasksRequest{ProjectId: folderId})
	if err != nil {
		return nil, fmt.Errorf("не удалось получить задачи из указанного проекта. details: %v", err)
	}

	return folderTasks, nil
}
