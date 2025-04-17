package database

import (
	"context"
	"fmt"
	"strings"
	"swagtask/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ---- READ ----
func GetTaskWithTagsById(dbpool *pgxpool.Pool, id int) (*TaskWithTags, error) {
	var taskWithTags TaskWithTags

	// we need left join because some tasks dont have ids
	rows, err := dbpool.Query(context.Background(), `
		SELECT t.id, t.name, t.idea, t.completed, tg.id, tg.name
		FROM tasks t
		LEFT JOIN tag_task_relations rel ON t.id = rel.task_id
		LEFT JOIN tags tg ON rel.tag_id = tg.id
		WHERE t.id = $1
	`, id)

	allTags, errTags := GetAllTags(dbpool)
	if errTags != nil {
		return nil, errTags
	}

	for rows.Next() {
		var task Task
		var tag_id *int      // nullable or can be inexistant
		var tag_name *string // nullable or can be inexistant

		rows.Scan(&task.Id, &task.Name, &task.Idea, &task.Completed, &tag_id, &tag_name)
		fmt.Println(task)
		if tag_id != nil && tag_name != nil {
			taskWithTags = NewTaskWithTags(task, append(taskWithTags.Tags, Tag{Id: *tag_id, Name: *tag_name}), allTags)
		} else {
			taskWithTags = NewTaskWithTags(
				task,
				[]Tag{},
				allTags,
			)
		}
	}

	if err != nil {
		return nil, err
	}

	return &taskWithTags, nil
}

func GetAllTasksWithTags(dbpool *pgxpool.Pool) ([]TaskWithTags, error) {
	rows, err := dbpool.Query(context.Background(), `
		SELECT t.id, t.name, t.idea, t.completed, tg.id, tg.name
		FROM tasks t
		LEFT JOIN tag_task_relations rel ON t.id = rel.task_id
		LEFT JOIN tags tg ON rel.tag_id = tg.id
		ORDER BY t.id DESC
	`)
	if err != nil {
		return nil, err
	}

	taskMap := make(map[int]*TaskWithTags)
	var orderedTasks []TaskWithTags

	for rows.Next() {
		var taskId int
		var taskName, idea string
		var completed bool
		var tagId *int      // nullable
		var tagName *string // nullable

		err := rows.Scan(&taskId, &taskName, &idea, &completed, &tagId, &tagName)
		if err != nil {
			return nil, err
		}

		allTags, errTags := GetAllTags(dbpool)
		if errTags != nil {
			return nil, err
		}

		_, exists := taskMap[taskId]
		if !exists {
			task := Task{Id: taskId, Name: taskName, Idea: idea, Completed: completed}
			taskWithTags := NewTaskWithTags(task, []Tag{}, allTags)
			taskMap[taskId] = &taskWithTags
			orderedTasks = append(orderedTasks, taskWithTags)
		}

		if tagId != nil && tagName != nil {
			tag := Tag{Id: *tagId, Name: *tagName}
			taskMap[taskId].Tags = append(taskMap[taskId].Tags, tag)
		}
	}

	var finalTasks []TaskWithTags
	for _, t := range orderedTasks {
		finalTasks = append(finalTasks, *taskMap[t.Task.Id])
	}

	return finalTasks, nil
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
	utils.PrintList(args)
	fmt.Println(str)
	if errTask != nil {
		return nil, errTask
	}

	return &task, nil

}
