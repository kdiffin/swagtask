package service

import (
	db "swagtask/db/generated"
	"swagtask/models"

	"github.com/jackc/pgx/v5/pgtype"
)

func getTaskAvailableTags(allTags []db.Tag, relatedTags []db.Tag) []db.Tag {
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
			ID: availableTask.ID,
		}
		if !taskExists[int(availableTask.ID)] {
			availableTasks = append(availableTasks, realAvailableTask)
		}
	}

	return availableTasks
}


func int32ToInt4Psql(i int32) pgtype.Int4 {
	var pgi pgtype.Int4
	pgi.Int32 = i	
	pgi.Valid = true

	return pgi
}
