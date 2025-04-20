package database

func GetTaskAvailableTags(allTags []Tag, relatedTags []Tag) []Tag {
	// think of this as a set checking if the tag is a tag of the task
	// int is id
	tagExists := make(map[int]bool)
	for _, tag := range relatedTags {
		tagExists[tag.Id] = true
	}

	availableTags := []Tag{}
	for _, tag := range allTags {
		if !tagExists[tag.Id] {
			availableTags = append(availableTags, tag)
		}
	}

	return availableTags
}

func GetTagAvailableTasks(allTasksOptions []AvailableTask, relatedTaskOptions []RelatedTask) []AvailableTask {
	// think of this as a set checking if the task is a task of the tag
	// int is id
	taskExists := make(map[int]bool)
	for _, task := range relatedTaskOptions {
		taskExists[task.Id] = true
	}

	availableTasks := []AvailableTask{}
	for _, availableTask := range allTasksOptions {
		if !taskExists[availableTask.Id] {
			availableTasks = append(availableTasks, availableTask)
		}
	}

	return availableTasks
}
