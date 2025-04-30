package service

import (
	"context"
	"errors"
	"fmt"
	db "swagtask/db/generated"
	"swagtask/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// ---- READ ----

func GetTaskWithTagsById(queries *db.Queries, id int32) (*models.TaskWithTags, error) {
    taskWithRelations, err := queries.GetTaskWithTagRelations(context.Background(), id)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("%w: %v", ErrBadRequest, err)
    }
    if len(taskWithRelations) == 0 {
        return nil, ErrNotFound
    }
    var task models.Task
    tagsOfTask := []db.Tag{}

    allTags, errTags := queries.GetAllTagsDesc(context.Background())
    if errTags != nil {
        if errors.Is(errTags, pgx.ErrNoRows) {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("%w: %v", ErrBadRequest, errTags)
    }

    for _, taskWithRelation := range taskWithRelations {
        task = models.Task{
            ID:        taskWithRelation.ID,
            Name:      taskWithRelation.Name,
            Idea:      taskWithRelation.Idea,
            Completed: taskWithRelation.Completed.Bool,
        }

        if taskWithRelation.TagID.Valid && taskWithRelation.TagName.Valid {
            tagsOfTask = append(tagsOfTask, db.Tag{ID: taskWithRelation.TagID.Int32, Name: taskWithRelation.TagName.String})
        }
    }

    availableTags := getTaskAvailableTags(allTags, tagsOfTask)
    taskWithTags := models.NewTaskWithTags(task, tagsOfTask, availableTags)
    return &taskWithTags, nil
}


func GetFilteredTasksWithTags(queries *db.Queries, filters *models.TasksPageFilters) ([]models.TaskWithTags, error) {
    fmt.Println(filters.TagName+":", "filter tagname")
    fmt.Println(filters.TaskName+":", "filter taskname")
    taskswithTagRelations, err := queries.GetFilteredTasks(context.Background(), db.GetFilteredTasksParams{
        TaskName: stringtoPgText(filters.TaskName),
        TagName: stringtoPgText(filters.TagName),
    }) 
    if err != nil {
        return nil, fmt.Errorf("%w: %v", ErrBadRequest, err)
    }


    allTags, errAllTags := queries.GetAllTagsDesc(context.Background())
    if errAllTags != nil {
        return nil, fmt.Errorf("%w: %v", ErrBadRequest, errAllTags)
    }

    tasksWithTags := []models.TaskWithTags{}
    taskIdToTags := make(map[int32][]db.Tag)
    idToTask := make(map[int32]db.Task)
    orderedIds := []int32{}
    idSeen := make(map[int32]bool)
    for _, task := range taskswithTagRelations {

        if task.TagID.Valid && task.TagName.Valid {
            taskIdToTags[task.ID] = append(taskIdToTags[task.ID], db.Tag{ID: task.TagID.Int32, Name: task.TagName.String})
        }

        if !idSeen[task.ID] {
            orderedIds = append(orderedIds, task.ID)
        }

        idToTask[task.ID] = db.Task{
            ID:        task.ID,
            Name:      task.Name,
            Idea:      task.Idea,
            Completed: task.Completed,
        }

        idSeen[task.ID] = true
    }

    for _, id := range orderedIds {
        task := idToTask[id]
        tagsOfTask := taskIdToTags[id]
        avaialbleTags := getTaskAvailableTags(allTags, tagsOfTask)

        taskWithTags := models.NewTaskWithTags(models.NewUITask(task), tagsOfTask, avaialbleTags)
        tasksWithTags = append(tasksWithTags, taskWithTags)
    }


    return tasksWithTags, nil
}
// func GetAllFilteredTasksWithTags(dbpool *pgxpool.Pool, params TasksPageFilters) ([]TaskWithTags, error) {
// 	// 0. gets the tags id by name
// 	// 1. gets the tasks with the param tag
// 	// 2. gets all tags
// 	// 3. finds the non related tags for the dropdown and then sends them
// 	index := 1
// 	var constraintStr string
// 	args := []interface{}{}

// 	if params.tagName != "" && params.taskName != "" {
// 		constraintStr += fmt.Sprintf(" tg.name = $%v", index)
// 		index++
// 		args = append(args, params.tagName)
// 		// sql wildcard
// 		constraintStr += fmt.Sprintf(" AND t.name ILIKE $%v", index)
// 		args = append(args, "%"+params.taskName+"%")
// 	} else if params.taskName != "" {
// 		constraintStr += fmt.Sprintf(" t.name ILIKE $%v", index)
// 		args = append(args, "%"+params.taskName+"%")
// 	} else if params.tagName != "" {
// 		constraintStr += fmt.Sprintf(" tg.name = $%v", index)
// 		args = append(args, params.tagName)
// 	} // handle the else case before u call this code

// 	utils.PrintList(args)
// 	queryString := `SELECT t.name, t.Idea, t.ID, t.completed, tg.ID, tg.name
// 		FROM tasks t
// 		LEFT JOIN tag_task_relations rel 
// 			ON t.ID = rel.task_id 
// 		LEFT JOIN tags tg 
// 			ON tg.ID = rel.tag_id
// 		WHERE` + constraintStr + " ORDER BY t.ID DESC"
// 	fmt.Println(queryString)
// 	rowsTasks, errTasks := dbpool.Query(context.Background(), queryString, args...)
// 	if errTasks != nil {
// 		return nil, errTasks
// 	}

// 	allTags, errAllTags := GetAllTags(dbpool)
// 	if errAllTags != nil {
// 		return nil, errAllTags
// 	}

// 	tasksWithTags := []TaskWithTags{}
// 	taskIdToTags := make(map[int][]Tag)
// 	// for easy lookup
// 	idToTask := make(map[int]Task)
// 	// so that we can make the items ordered (this is kindof a set)
// 	orderedIds := []int{}
// 	// so that we can mimic the set logic
// 	idSeen := make(map[int]bool)
// 	for rowsTasks.Next() {
// 		var task Task
// 		var tagId *int      // nullable
// 		var tagName *string // nullable

// 		errScanTask := rowsTasks.Scan(&task.Name, &task.Idea, &task.ID, &task.Completed, &tagId, &tagName)
// 		if errScanTask != nil {
// 			return nil, errScanTask
// 		}

// 		if tagId != nil && tagName != nil {
// 			// add tag if it exists
// 			taskIdToTags[task.ID] = append(taskIdToTags[task.ID], Tag{Id: *tagId, Name: *tagName})
// 		}

// 		// mimic a set
// 		if !idSeen[task.ID] {
// 			orderedIds = append(orderedIds, task.ID)
// 		}

// 		idToTask[task.ID] = task
// 		idSeen[task.ID] = true
// 	}

// 	for _, id := range orderedIds {
// 		task := idToTask[id]
// 		tagsOfTask := taskIdToTags[id]
// 		avaialbleTags := GetTaskAvailableTags(allTags, tagsOfTask)

// 		taskWithTags := NewTaskWithTags(task, tagsOfTask, avaialbleTags)
// 		tasksWithTags = append(tasksWithTags, taskWithTags)
// 	}

// 	return tasksWithTags, nil
// }


// ---- CREATE ----
func CreateTask(queries *db.Queries, name string, idea string) (*models.TaskWithTags, error) {
    task, errCreate := queries.CreateTask(context.Background(), db.CreateTaskParams{Name: name, Idea: idea})
    if errCreate != nil {
        if errors.Is(errCreate, pgx.ErrNoRows) {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("%w: %v", ErrUnprocessable, errCreate)
    }

    allTags, errAllTags := queries.GetAllTagsDesc(context.Background())
    if errAllTags != nil {
        if errors.Is(errAllTags, pgx.ErrNoRows) {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("%w: %v", ErrBadRequest, errAllTags)
    }

    taskWithTag := models.NewTaskWithTags(
        models.NewUITask(task),
        []db.Tag{},
        allTags,
    )

    return &taskWithTag, nil
}

// ---- UPDATE ----
func UpdateTaskCompletion(queries *db.Queries, taskId int32) (*models.TaskWithTags, error) {
    errCompletion := queries.ToggleTaskCompletion(context.Background(), taskId)
    if errCompletion != nil {
        if errors.Is(errCompletion, pgx.ErrNoRows) {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("%w: %v", ErrBadRequest, errCompletion)
    }

    taskWithTags, errGetTask := GetTaskWithTagsById(queries, taskId)
    if errGetTask != nil {
        return nil, errGetTask
    }

    return taskWithTags, nil
}

func UpdateTask(queries *db.Queries, taskId int32, name string, idea string) (*models.TaskWithTags, error) {
    var nameNullable pgtype.Text
    var ideaNullable pgtype.Text

    if name == "" {
        nameNullable.Valid = false
    } else {
        nameNullable.Valid = true
        nameNullable.String = name
    }
    if idea == "" {
        ideaNullable.Valid = false
    } else {
        ideaNullable.Valid = true
        ideaNullable.String = idea
    }
    if !nameNullable.Valid && !ideaNullable.Valid {
        return nil, ErrNoUpdateFields
    }

    errCompletion := queries.UpdateTask(context.Background(), db.UpdateTaskParams{
        ID:   taskId,
        Name: nameNullable,
        Idea: ideaNullable,
    })
    if errCompletion != nil {
        if errors.Is(errCompletion, pgx.ErrNoRows) {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("%w: %v", ErrUnprocessable, errCompletion)
    }

    taskWithTags, errGetTask := GetTaskWithTagsById(queries, taskId)
    if errGetTask != nil {
        return nil, errGetTask
    }

    return taskWithTags, nil
}

func AddTagToTask(queries *db.Queries, tagId int32, taskId int32) (*models.TaskWithTags, error) {
    err := queries.CreateTagTaskRelation(context.Background(), db.CreateTagTaskRelationParams{
        TaskID: int32ToPgInt4(taskId),
        TagID:  int32ToPgInt4(tagId),
    })
	
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("%w: %v", ErrBadRequest, err)
    }

    taskWithTags, errTasks := GetTaskWithTagsById(queries, taskId)
    if errTasks != nil {
        return nil, errTasks
    }

    return taskWithTags, nil
}

// ---- DELETE ----
func DeleteTask(queries *db.Queries, taskId int32) error {
    errRelations := queries.DeleteAllTaskRelations(context.Background(), int32ToPgInt4(taskId))
    if errRelations != nil {
        if errors.Is(errRelations, pgx.ErrNoRows) {
            return ErrNotFound
        }
        return fmt.Errorf("%w: %v", ErrBadRequest, errRelations)
    }

    err := queries.DeleteTask(context.Background(), taskId)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return ErrNotFound
        }
        return fmt.Errorf("%w: %v", ErrBadRequest, err)
    }

    return nil
}

func DeleteTagRelationFromTask(queries *db.Queries, tagId int32, taskId int32) (*models.TaskWithTags, error) {
    errRelations := queries.DeleteTagTaskRelation(context.Background(), db.DeleteTagTaskRelationParams{
        TaskID: int32ToPgInt4(taskId),
        TagID:  int32ToPgInt4(tagId)})
		fmt.Println(tagId)
		fmt.Println(taskId)
    
		if errRelations != nil {
        if errors.Is(errRelations, pgx.ErrNoRows) {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("%w: %v", ErrBadRequest, errRelations)
    }

    taskWithTags, err := GetTaskWithTagsById(queries, taskId)
    if err != nil {
        return nil, err
    }

    return taskWithTags, nil
}