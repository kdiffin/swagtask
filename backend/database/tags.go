package database

import (
	"context"

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

func GetAllTagsWithTasks(dbpool *pgxpool.Pool) ([]TagWithTasks, error) {
	tags, err := GetAllTags(dbpool)
	rows, errTasks := dbpool.Query(context.Background(), "SELECT name,id FROM tasks")
	if err != nil || errTasks != nil {
		return nil, err
	}

	allTags := []TagRelationOption{}
	for rows.Next() {
		var taskName string
		var taskId int
		rows.Scan(&taskName, &taskId)
		allTags = append(allTags, TagRelationOption{Name: taskName, Id: taskId})
	}
	tagsWithTasksOptions := []TagWithTasks{}
	for _, tag := range tags {
		tagWithTasks := NewTagWithTasks(tag.Id, tag.Name, allTags)
		tagsWithTasksOptions = append(tagsWithTasksOptions, tagWithTasks)

	}

	return tagsWithTasksOptions, nil
}
