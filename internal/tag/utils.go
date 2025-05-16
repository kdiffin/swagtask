package tag

import (
	db "swagtask/internal/db/generated"
)

func getTagAvailableTasks(allTasksOptions []db.GetAllTaskOptionsRow, relatedTaskOptions []relatedTask) []availableTask {
	// think of this as a set checking if the task is a task of the tag
	// int is id
	taskExists := make(map[string]bool)
	for _, task := range relatedTaskOptions {
		taskExists[task.ID] = true
	}

	availableTasks := []availableTask{}
	for _, taskOption := range allTasksOptions {
		realAvailableTask := availableTask{
			Name: taskOption.Name,
			ID:   taskOption.ID.String(),
		}
		if !taskExists[taskOption.ID.String()] {
			availableTasks = append(availableTasks, realAvailableTask)
		}
	}

	return availableTasks
}
