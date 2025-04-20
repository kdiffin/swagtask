package database

import (
	"context"
	"fmt"
	"strings"

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

// redo this whole thing
// TODO: fixed the n+1 by querying all tags, ill do it with joins again later
func GetAllTasksWithTags(dbpool *pgxpool.Pool) ([]TaskWithTags, error) {
	// get the tasks
	// get the relation table
	// fetch and append any tag in the rel table whichs task_id to the respective task

	rows, errTasks := dbpool.Query(context.Background(), "SELECT name,idea,id,completed FROM tasks ORDER BY id DESC")
	if errTasks != nil {
		return nil, errTasks
	}
	rowsRel, errRel := dbpool.Query(context.Background(), "SELECT task_id, tag_id FROM tag_task_relations")
	if errRel != nil {
		return nil, errRel
	}
	tags, errTags := GetAllTags(dbpool)
	if errTags != nil {
		return nil, errTags
	}

	tasks := []Task{}
	tasksTags := make(map[int][]Tag)
	for rows.Next() {
		var task Task
		rows.Scan(&task.Name, &task.Idea, &task.Id, &task.Completed)
		tasks = append(tasks, task)
		tasksTags[task.Id] = []Tag{}

	}

	// for fast lookup
	idToTag := make(map[int]Tag)
	for _, tag := range tags {
		idToTag[tag.Id] = tag
	}
	for rowsRel.Next() {
		var taskId int
		var tagId int
		rowsRel.Scan(&taskId, &tagId)
		tasksTags[taskId] = append(tasksTags[taskId], idToTag[tagId])
	}

	tasksWithTags := []TaskWithTags{}

	// maybe make this faster
	for _, task := range tasks {
		availableTags := []Tag{}
		for _, tag := range tags {
			notInTask := true
			for _, taskTag := range tasksTags[task.Id] {
				if taskTag == tag {
					notInTask = false
				}
			}

			if notInTask {
				availableTags = append(availableTags, tag)
			}
		}
		taskWithTag := NewTaskWithTags(task, tasksTags[task.Id], availableTags)
		tasksWithTags = append(tasksWithTags, taskWithTag)
	}

	return tasksWithTags, nil

}

func GetAllFilteredTasksWithTags(dbpool *pgxpool.Pool, tagFilter string) ([]TaskWithTags, error) {
	var tagFilterId string
	errGetIdOfFilter := dbpool.QueryRow(context.Background(), `
		SELECT id  
		FROM tags 
		WHERE tags.name = $1 
		LIMIT 1`, tagFilter).Scan(&tagFilterId)

	// TODO: maybe implement a badgateway error with errors
	if errGetIdOfFilter != nil {
		return nil, errGetIdOfFilter
	}

	// get the tasks related to this tag
	rowsIdOfTasksWithTag, errGettingTasks := dbpool.Query(context.Background(), `SELECT task_id FROM tag_task_relations rel WHERE rel.tag_id = $1`, tagFilterId)
	if errGettingTasks != nil {
		return nil, errGettingTasks
	}

	idOfTasksToQuery := []int{}
	for rowsIdOfTasksWithTag.Next() {
		var taskId int
		errIdOfTaskScan := rowsIdOfTasksWithTag.Scan(&taskId)
		if errIdOfTaskScan != nil {
			return nil, errIdOfTaskScan
		}

		idOfTasksToQuery = append(idOfTasksToQuery, taskId)
	}

	finalTasks := []TaskWithTags{}
	for _, id := range idOfTasksToQuery {
		// n+1 query eliyende oglanlar
		// TODO: FIX THIS STUPID AHH N+1 ISSUE!
		taskWithTags, errGetTask := GetTaskWithTagsById(dbpool, id)

		if errGetTask != nil {
			return nil, errGetTask
		}

		finalTasks = append(finalTasks, *taskWithTags)

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
	if errTask != nil {
		return nil, errTask
	}

	return &task, nil

}
