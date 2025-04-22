package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetTagsOfTask(dbpool *pgxpool.Pool, id int) ([]Tag, error) {
	tagsOfTask := []Tag{}

	rows, errTags := dbpool.Query(context.Background(), `SELECT tg.name, tg.id FROM tags tg JOIN tag_task_relations rel ON rel.tag_id = tg.id WHERE rel.task_id = $1`, id)
	if errTags != nil {
		return nil, errTags
	}

	for rows.Next() {
		var tag Tag
		rows.Scan(&tag.Name, &tag.Id)
		tagsOfTask = append(tagsOfTask, tag)
	}

	return tagsOfTask, nil
}

func GetAllTags(dbpool *pgxpool.Pool) ([]Tag, error) {
	allTags := []Tag{}
	rows, errTags := dbpool.Query(context.Background(), `SELECT name, id FROM tags ORDER BY id DESC`)
	if errTags != nil {
		return nil, errTags
	}
	for rows.Next() {
		var tag Tag
		rows.Scan(&tag.Name, &tag.Id)
		allTags = append(allTags, tag)
	}

	return allTags, nil
}

func GetRelatedTasksOfTag(dbpool *pgxpool.Pool, id int) ([]RelatedTask, error) {
	rows, err := dbpool.Query(context.Background(), `
		SELECT t.name, t.id FROM tasks t
		JOIN tag_task_relations rel ON t.id = rel.task_id
		JOIN tags tg ON tg.id = rel.tag_id
		WHERE tg.id = $1`, id)
	if err != nil {
		return nil, err
	}

	relatedTasks := []RelatedTask{}
	for rows.Next() {
		var relatedTask RelatedTask
		errScan := rows.Scan(&relatedTask.Name, &relatedTask.Id)
		if errScan != nil {
			fmt.Println("err on rows scan")
			return nil, errScan
		}
		relatedTasks = append(relatedTasks, relatedTask)
	}

	return relatedTasks, nil
}

func GetTagById(dbpool *pgxpool.Pool, id int) (*Tag, error) {
	var tag Tag
	errTags := dbpool.QueryRow(context.Background(), `
		SELECT name, id 
		FROM tags 
		WHERE id = $1`, id).Scan(&tag.Name, &tag.Id)

	if errTags != nil {
		return nil, errTags
	}

	return &tag, nil
}

func GetTagWithTasks(dbpool *pgxpool.Pool, id int) (*TagWithTasks, error) {
	rows, errTasks := dbpool.Query(context.Background(), "SELECT name,id FROM tasks")
	if errTasks != nil {
		fmt.Println("error here tasks query")
		return nil, errTasks
	}

	allTasks := []AvailableTask{}
	for rows.Next() {
		var option AvailableTask
		rows.Scan(&option.Name, &option.Id)

		allTasks = append(allTasks, option)
	}

	relatedTasks, errRelatedTasks := GetRelatedTasksOfTag(dbpool, id)
	if errRelatedTasks != nil {
		fmt.Println("error here relations")
		return nil, errRelatedTasks
	}

	tag, errTag := GetTagById(dbpool, id)
	if errTag != nil {
		return nil, errTag
	}
	availableTasks := GetTagAvailableTasks(allTasks, relatedTasks)
	tagWithTasks := NewTagWithTasks(*tag, relatedTasks, availableTasks)
	return &tagWithTasks, nil
}
func GetAllTagsWithTasks(dbpool *pgxpool.Pool) ([]TagWithTasks, error) {
	// 1. get tags with their relations to tasks
	// 2. get all tasks as options
	// 3. get the tasks that aint related to task

	rows, err := dbpool.Query(context.Background(), `
		SELECT tg.id, tg.name, t.id, t.name 
		FROM tags tg
		LEFT JOIN tag_task_relations rel
			ON tg.id = rel.tag_id
		LEFT JOIN tasks t
			ON t.id = rel.task_id
		ORDER BY tg.id DESC`)
	if err != nil {
		return nil, err
	}
	rowsAllTasks, errRowsAllTasks := dbpool.Query(context.Background(), `SELECT id, name FROM tasks`)
	if errRowsAllTasks != nil {
		return nil, errRowsAllTasks
	}

	allTasks := []AvailableTask{}
	for rowsAllTasks.Next() {
		var task AvailableTask
		rowsAllTasks.Scan(&task.Id, &task.Name)
		allTasks = append(allTasks, task)
	}

	idToTag := make(map[int]Tag)
	tagIdToTasks := make(map[int][]RelatedTask)
	// simulating a set
	orderedIds := []int{}
	seenId := make(map[int]bool)
	for rows.Next() {
		var tag Tag
		var taskId *int
		var taskName *string
		rows.Scan(&tag.Id, &tag.Name, &taskId, &taskName)

		if taskId != nil && taskName != nil {
			tagIdToTasks[tag.Id] = append(tagIdToTasks[tag.Id], RelatedTask{Name: *taskName, Id: *taskId})
		}

		if !seenId[tag.Id] {
			orderedIds = append(orderedIds, tag.Id)
		}

		seenId[tag.Id] = true
		idToTag[tag.Id] = tag
	}

	tagsWithTasks := []TagWithTasks{}
	for _, id := range orderedIds {
		availableTasks := GetTagAvailableTasks(allTasks, tagIdToTasks[id])
		tagsWithTasks = append(tagsWithTasks, NewTagWithTasks(idToTag[id], tagIdToTasks[id], availableTasks))
	}

	// for _,
	// tagsWithTasksOptions = append(tagsWithTasksOptions, tagWithTasks)
	// availableTasks := GetTagAvailableTasks(allTasks, )

	return tagsWithTasks, nil
}
