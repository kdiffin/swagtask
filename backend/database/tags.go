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
	rows, errTags := dbpool.Query(context.Background(), `SELECT name, id FROM tags`)
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

func GetTagById(dbpool *pgxpool.Pool) ([]Tag, error) {
	allTags := []Tag{}
	rows, errTags := dbpool.Query(context.Background(), `SELECT name, id FROM tags`)
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
