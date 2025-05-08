package models

import (
	db "swagtask/db/generated"
)

// ---- DB TO UI MAPPING ----
type Task struct {
	ID        int32
	Name      string
	Idea      string
	Completed bool
}

func NewUITask(task db.Task) Task{
	if !task.Completed.Valid {
		panic("this wasnt supposed to happen since the db has false as its default")
	}

	return Task{
		ID: task.ID,
		Name: task.Name,
		Idea:  task.Idea,
		Completed: task.Completed.Bool,
	}
}


// ---- FOR UI ----
// tasks
type TaskWithTags struct {
	Task
	Tags          []db.Tag
	AvailableTags []db.Tag
}

func NewTaskWithTags(task Task, tags []db.Tag, availableTags []db.Tag) TaskWithTags {
	return TaskWithTags{
		Task:          task,
		Tags:          tags,
		AvailableTags: availableTags,
	}
}

// tags
type AvailableTask struct {
	Name string
	ID   int32
}
type RelatedTask struct {
	Name string
	ID   int32
}
type TagWithTasks struct {
	db.Tag
	RelatedTasks   []RelatedTask
	AvailableTasks []AvailableTask
}

func NewTagWithTasks(tag db.Tag, relatedTasks []RelatedTask, availableTasks []AvailableTask) TagWithTasks {
	return TagWithTasks{
		Tag:            tag,
		RelatedTasks:   relatedTasks,
		AvailableTasks: availableTasks,
	}
}

type TasksPageFilters struct {
	SearchQuery  string
	ActiveTag string
}

func NewTasksPageFilters(tagName string, taskName string) TasksPageFilters {
	return TasksPageFilters{
		ActiveTag:  tagName,
		SearchQuery: taskName,
	}
}


// user
type User struct{
	ID int32
	PathToPfp string
	Username string
}