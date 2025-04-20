package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ---- READ ----
func GetTaskWithTagsById(dbpool *pgxpool.Pool, id int) (*TaskWithTags, error) {
	// we need left join because some tasks dont have ids
	rows, err := dbpool.Query(context.Background(), `
			SELECT t.id, t.name, t.idea, t.completed, tg.id, tg.name
			FROM tasks t
			LEFT JOIN tag_task_relations rel ON t.id = rel.task_id
			LEFT JOIN tags tg ON rel.tag_id = tg.id
			WHERE t.id = $1
	`, id)
	if err != nil {
		return nil, err
	}

	var task Task
	tagsOfTask := []Tag{}

	allTags, errTags := GetAllTags(dbpool)
	if errTags != nil {
		return nil, errTags
	}
	for rows.Next() {
		var tag_id *int      // nullable or can be inexistant
		var tag_name *string // nullable or can be inexistant

		errScan := rows.Scan(&task.Id, &task.Name, &task.Idea, &task.Completed, &tag_id, &tag_name)
		if errScan != nil {
			return nil, errScan
		}

		if tag_id != nil && tag_name != nil {
			tagsOfTask = append(tagsOfTask, Tag{Id: *tag_id, Name: *tag_name})
		}
	}

	availableTags := GetTaskAvailableTags(allTags, tagsOfTask)
	taskWithTags := NewTaskWithTags(task, tagsOfTask, availableTags)
	return &taskWithTags, nil
}

func GetAllTasksWithTags(dbpool *pgxpool.Pool) ([]TaskWithTags, error) {
	// 1. gets the task with their related tags
	// 2. gets all tags
	// 3. finds the non related tags for the dropdown and then sends them
	// TODO: make these run in paralell
	rowsTasks, errTasks := dbpool.Query(context.Background(), `
		SELECT t.name, t.idea, t.id, t.completed, tg.id, tg.name
		FROM tasks t
		LEFT JOIN tag_task_relations rel 
			ON t.id = rel.task_id 
		LEFT JOIN tags tg 
			ON tg.id = rel.tag_id
		ORDER BY t.id DESC`)
	if errTasks != nil {
		return nil, errTasks
	}

	allTags, errAllTags := GetAllTags(dbpool)
	if errAllTags != nil {
		return nil, errAllTags
	}

	tasksWithTags := []TaskWithTags{}
	taskIdToTags := make(map[int][]Tag)
	// for easy lookup
	idToTask := make(map[int]Task)
	// so that we can make the items ordered (this is kindof a set)
	orderedIds := []int{}
	// so that we can mimic the set logic
	idSeen := make(map[int]bool)
	for rowsTasks.Next() {
		var task Task
		var tagId *int      // nullable
		var tagName *string // nullable

		errScanTask := rowsTasks.Scan(&task.Name, &task.Idea, &task.Id, &task.Completed, &tagId, &tagName)
		if errScanTask != nil {
			return nil, errScanTask
		}

		if tagId != nil && tagName != nil {
			// add tag if it exists
			taskIdToTags[task.Id] = append(taskIdToTags[task.Id], Tag{Id: *tagId, Name: *tagName})
		}

		// mimic a set
		if !idSeen[task.Id] {
			orderedIds = append(orderedIds, task.Id)
		}

		idToTask[task.Id] = task
		idSeen[task.Id] = true
	}

	for _, id := range orderedIds {
		task := idToTask[id]
		tagsOfTask := taskIdToTags[id]
		avaialbleTags := GetTaskAvailableTags(allTags, tagsOfTask)

		taskWithTags := NewTaskWithTags(task, tagsOfTask, avaialbleTags)
		tasksWithTags = append(tasksWithTags, taskWithTags)
	}

	return tasksWithTags, nil
}

func GetAllFilteredTasksWithTags(dbpool *pgxpool.Pool, param string) ([]TaskWithTags, error) {
	// 0. gets the tags id by name
	// 1. gets the tasks with the param tag
	// 2. gets all tags
	// 3. finds the non related tags for the dropdown and then sends them

	// 0.
	var tagId int
	// TODO: implement custom errors here
	errTagId := dbpool.QueryRow(context.Background(), `SELECT id FROM tags WHERE name = $1`, param).Scan(&tagId)
	if errTagId != nil {
		return nil, errTagId
	}

	rowsTasks, errTasks := dbpool.Query(context.Background(), `
		SELECT t.name, t.idea, t.id, t.completed, tg.id, tg.name
		FROM tasks t
		LEFT JOIN tag_task_relations rel 
			ON t.id = rel.task_id 
		LEFT JOIN tags tg 
			ON tg.id = rel.tag_id
		WHERE tg.id = $1
		ORDER BY t.id DESC`, tagId)
	if errTasks != nil {
		return nil, errTasks
	}

	allTags, errAllTags := GetAllTags(dbpool)
	if errAllTags != nil {
		return nil, errAllTags
	}

	tasksWithTags := []TaskWithTags{}
	taskIdToTags := make(map[int][]Tag)
	// for easy lookup
	idToTask := make(map[int]Task)
	// so that we can make the items ordered (this is kindof a set)
	orderedIds := []int{}
	// so that we can mimic the set logic
	idSeen := make(map[int]bool)
	for rowsTasks.Next() {
		var task Task
		var tagId *int      // nullable
		var tagName *string // nullable

		errScanTask := rowsTasks.Scan(&task.Name, &task.Idea, &task.Id, &task.Completed, &tagId, &tagName)
		if errScanTask != nil {
			return nil, errScanTask
		}

		if tagId != nil && tagName != nil {
			// add tag if it exists
			taskIdToTags[task.Id] = append(taskIdToTags[task.Id], Tag{Id: *tagId, Name: *tagName})
		}

		// mimic a set
		if !idSeen[task.Id] {
			orderedIds = append(orderedIds, task.Id)
		}

		idToTask[task.Id] = task
		idSeen[task.Id] = true
	}

	for _, id := range orderedIds {
		task := idToTask[id]
		tagsOfTask := taskIdToTags[id]
		avaialbleTags := GetTaskAvailableTags(allTags, tagsOfTask)

		taskWithTags := NewTaskWithTags(task, tagsOfTask, avaialbleTags)
		tasksWithTags = append(tasksWithTags, taskWithTags)
	}

	return tasksWithTags, nil
}

// ---- CREATE ----
func CreateTask(dbpool *pgxpool.Pool, name string, idea string) (*Task, error) {
	var task Task
	err := dbpool.QueryRow(context.Background(), "INSERT INTO tasks (name, idea) VALUES ($1, $2) RETURNING name, idea, id, completed",
		name, idea).Scan(&task.Name, &task.Idea, &task.Id, &task.Completed)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

// ---- UPDATE ---

func UpdateTask(dbpool *pgxpool.Pool, name string, idea string, id int) (*Task, error) {
	var task Task
	args := []interface{}{}
	n := 1

	updateTaskString := "UPDATE tasks SET"

	if name != "" {
		updateTaskString += fmt.Sprintf(" name = $%v,", n)
		args = append(args, name)
		n++
	}
	if idea != "" {
		updateTaskString += fmt.Sprintf(" idea = $%v,", n)
		args = append(args, idea)
		n++
	}
	// nothing to update if no args added
	if len(args) == 0 {
		return nil, ErrNoUpdateFields
	}

	args = append(args, id)

	// remove trailing comma
	updateTaskString = strings.TrimSuffix(updateTaskString, ",")
	str := updateTaskString + fmt.Sprintf(" WHERE id = $%v RETURNING name,idea,id,completed", n)
	errTask := dbpool.QueryRow(context.Background(), str, args...).Scan(&task.Name, &task.Idea, &task.Id, &task.Completed)
	if errTask != nil {
		return nil, errTask
	}

	return &task, nil

}
