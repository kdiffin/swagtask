package service

import (
	db "swagtask/db/generated"

	"github.com/jackc/pgx/v5/pgtype"
)

func GetTaskAvailableTags(allTags []db.Tag, relatedTags []db.Tag) []db.Tag {
	// think of this as a set checking if the tag is a tag of the task
	// int is id
	tagExists := make(map[int32]bool)
	for _, tag := range relatedTags {
		tagExists[tag.ID] = true
	}

	availableTags := []db.Tag{}
	for _, tag := range allTags {
		if !tagExists[tag.ID] {
			availableTags = append(availableTags, tag)
		}
	}

	return availableTags
}

func Int32ToInt4Psql(i int32) pgtype.Int4 {
	var pgi pgtype.Int4
	pgi.Int32 = i
	pgi.Valid = true

	return pgi
}
// func GetTagAvailableTasks(allTasksOptions []AvailableTask, relatedTaskOptions []RelatedTask) []AvailableTask {
// 	// think of this as a set checking if the task is a task of the tag
// 	// int is id
// 	taskExists := make(map[int]bool)
// 	for _, task := range relatedTaskOptions {
// 		taskExists[task.ID] = true
// 	}

// 	availableTasks := []AvailableTask{}
// 	for _, availableTask := range allTasksOptions {
// 		if !taskExists[availableTask.ID] {
// 			availableTasks = append(availableTasks, availableTask)
// 		}
// 	}

// 	return availableTasks
// }
