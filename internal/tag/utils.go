package tag

import (
	"swagtask/backend/models"
	db "swagtask/internal/db/generated"
)

func getTagAvailableTasks(allTasksOptions []db.GetAllTaskOptionsRow, relatedTaskOptions []models.RelatedTask) []models.AvailableTask {
	// think of this as a set checking if the task is a task of the tag
	// int is id
	taskExists := make(map[int]bool)
	for _, task := range relatedTaskOptions {
		taskExists[int(task.ID)] = true
	}

	availableTasks := []models.AvailableTask{}
	for _, availableTask := range allTasksOptions {
		realAvailableTask := models.AvailableTask{
			Name: availableTask.Name,
			ID:   availableTask.ID,
		}
		if !taskExists[availableTask.ID.String()] {
			availableTasks = append(availableTasks, realAvailableTask)
		}
	}

	return availableTasks
}
