package service

// import (
// 	"context"
// 	"fmt"

// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// func GetRelatedTasksOfTag(dbpool *pgxpool.Pool, id int) ([]RelatedTask, error) {
// 	rows, err := dbpool.Query(context.Background(), `
// 		SELECT t.name, t.ID FROM tasks t
// 		JOIN tag_task_relations rel ON t.ID = rel.task_id
// 		JOIN tags tg ON tg.ID = rel.tag_id
// 		WHERE tg.ID = $1`, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	relatedTasks := []RelatedTask{}
// 	for rows.Next() {
// 		var relatedTask RelatedTask
// 		errScan := rows.Scan(&relatedTask.Name, &relatedTask.ID)
// 		if errScan != nil {
// 			fmt.Println("err on rows scan")
// 			return nil, errScan
// 		}
// 		relatedTasks = append(relatedTasks, relatedTask)
// 	}

// 	return relatedTasks, nil
// }

// func GetTagWithTasks(dbpool *pgxpool.Pool, id int) (*TagWithTasks, error) {
// 	rows, errTasks := dbpool.Query(context.Background(), "SELECT name,id FROM tasks")
// 	if errTasks != nil {
// 		fmt.Println("error here tasks query")
// 		return nil, errTasks
// 	}

// 	allTasks := []AvailableTask{}
// 	for rows.Next() {
// 		var option AvailableTask
// 		rows.Scan(&option.Name, &option.ID)

// 		allTasks = append(allTasks, option)
// 	}

// 	relatedTasks, errRelatedTasks := GetRelatedTasksOfTag(dbpool, id)
// 	if errRelatedTasks != nil {
// 		fmt.Println("error here relations")
// 		return nil, errRelatedTasks
// 	}

// 	tag, errTag := GetTagById(dbpool, id)
// 	if errTag != nil {
// 		return nil, errTag
// 	}
// 	availableTasks := GetTagAvailableTasks(allTasks, relatedTasks)
// 	tagWithTasks := NewTagWithTasks(*tag, relatedTasks, availableTasks)
// 	return &tagWithTasks, nil
// }
// func GetAllTagsWithTasks(dbpool *pgxpool.Pool) ([]TagWithTasks, error) {
// 	// 1. get tags with their relations to tasks
// 	// 2. get all tasks as options
// 	// 3. get the tasks that aint related to task

// 	rows, err := dbpool.Query(context.Background(), `
// 		SELECT tg.ID, tg.name, t.ID, t.name
// 		FROM tags tg
// 		LEFT JOIN tag_task_relations rel
// 			ON tg.ID = rel.tag_id
// 		LEFT JOIN tasks t
// 			ON t.ID = rel.task_id
// 		ORDER BY tg.ID DESC`)
// 	if err != nil {
// 		return nil, err
// 	}
// 	rowsAllTasks, errRowsAllTasks := dbpool.Query(context.Background(), `SELECT id, name FROM tasks`)
// 	if errRowsAllTasks != nil {
// 		return nil, errRowsAllTasks
// 	}

// 	allTasks := []AvailableTask{}
// 	for rowsAllTasks.Next() {
// 		var task AvailableTask
// 		rowsAllTasks.Scan(&task.ID, &task.Name)
// 		allTasks = append(allTasks, task)
// 	}

// 	idToTag := make(map[int]Tag)
// 	tagIdToTasks := make(map[int][]RelatedTask)
// 	// simulating a set
// 	orderedIds := []int{}
// 	seenId := make(map[int]bool)
// 	for rows.Next() {
// 		var tag Tag
// 		var taskId *int
// 		var taskName *string
// 		rows.Scan(&tag.ID, &tag.Name, &taskId, &taskName)

// 		if taskId != nil && taskName != nil {
// 			tagIdToTasks[tag.ID] = append(tagIdToTasks[tag.ID], RelatedTask{Name: *taskName, Id: *taskId})
// 		}

// 		if !seenId[tag.ID] {
// 			orderedIds = append(orderedIds, tag.ID)
// 		}

// 		seenId[tag.ID] = true
// 		idToTag[tag.ID] = tag
// 	}

// 	tagsWithTasks := []TagWithTasks{}
// 	for _, id := range orderedIds {
// 		availableTasks := GetTagAvailableTasks(allTasks, tagIdToTasks[id])
// 		tagsWithTasks = append(tagsWithTasks, NewTagWithTasks(idToTag[id], tagIdToTasks[id], availableTasks))
// 	}

// 	// for _,
// 	// tagsWithTasksOptions = append(tagsWithTasksOptions, tagWithTasks)
// 	// availableTasks := GetTagAvailableTasks(allTasks, )

// 	return tagsWithTasks, nil
// }
