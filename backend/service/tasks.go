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

func GetTaskWithTagsById(queries *db.Queries, userId int32,
     id int32, ctx context.Context) (*models.TaskWithTags, error) {
    taskWithRelations, err := queries.GetTaskWithTagRelations(ctx, db.GetTaskWithTagRelationsParams{
        ID: id,
        UserID: userId,
    })
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

    allTags, errTags := queries.GetAllTagsDesc(ctx, userId)
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


func GetFilteredTasksWithTags(queries *db.Queries, filters *models.TasksPageFilters,
    userId int32, ctx context.Context) ([]models.TaskWithTags, error) {
    taskswithTagRelations, err := queries.GetFilteredTasks(ctx, db.GetFilteredTasksParams{
        TaskName: stringtoPgText(filters.SearchQuery),
        TagName: stringtoPgText(filters.ActiveTag),
    })  
    if err != nil {
        return nil, fmt.Errorf("%w: %v", ErrBadRequest, err)
    }


    allTags, errAllTags := queries.GetAllTagsDesc(ctx, userId)
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

func GetTaskNavigationButtons(ctx context.Context, queries *db.Queries,userId int32, id int32) (models.TaskButton, models.TaskButton) {
	prev, errPrev := queries.GetPreviousTaskDetails(ctx, db.GetPreviousTaskDetailsParams{
        ID: id,
        UserID: userId,
    })
	next, errNext := queries.GetNextTaskDetails(ctx, db.GetNextTaskDetailsParams{
        ID: id,
        UserID: userId,
    })

	var prevButton, nextButton models.TaskButton
	if errPrev == nil {
		prevButton = models.TaskButton{ID: prev.ID, Name: prev.Name, Exists: true}
	}
	if errNext == nil {
		nextButton = models.TaskButton{ID: next.ID, Name: next.Name, Exists: true}
	}
	return prevButton, nextButton
}


// ---- CREATE ----
func CreateTask(queries *db.Queries, name string, 
    userId int32, idea string, ctx context.Context) (*models.TaskWithTags, error) {
    task, errCreate := queries.CreateTask(ctx, db.CreateTaskParams{Name: name, Idea: idea, UserID: userId})
    if errCreate != nil {
        if errors.Is(errCreate, pgx.ErrNoRows) {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("%w: %v", ErrUnprocessable, errCreate)
    }

    allTags, errAllTags := queries.GetAllTagsDesc(ctx, userId)
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
func UpdateTaskCompletion(queries *db.Queries, userId int32, taskId int32, ctx context.Context) (*models.TaskWithTags, error) {
    errCompletion := queries.ToggleTaskCompletion(ctx, db.ToggleTaskCompletionParams{
        ID: taskId,
        UserID: userId,
    })
    if errCompletion != nil {
        if errors.Is(errCompletion, pgx.ErrNoRows) {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("%w: %v", ErrBadRequest, errCompletion)
    }

    taskWithTags, errGetTask := GetTaskWithTagsById(queries, userId, taskId, ctx)
    if errGetTask != nil {
        return nil, errGetTask
    }

    return taskWithTags, nil
}

func UpdateTask(queries *db.Queries, taskId int32, userId int32,
     name string, idea string, ctx context.Context) (*models.TaskWithTags, error) {
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

    errCompletion := queries.UpdateTask(ctx, db.UpdateTaskParams{
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

    taskWithTags, errGetTask := GetTaskWithTagsById(queries, userId, taskId, ctx)
    if errGetTask != nil {
        return nil, errGetTask
    }

    return taskWithTags, nil
}

func AddTagToTask(queries *db.Queries,
     tagId, userId, taskId int32, ctx context.Context) (*models.TaskWithTags, error) {
    err := queries.CreateTagTaskRelation(ctx, db.CreateTagTaskRelationParams{
        TaskID: taskId,
        TagID:  tagId,
    })
	
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, ErrNotFound
        }   
        return nil, fmt.Errorf("%w: %v", ErrBadRequest, err)
    }

    taskWithTags, errTasks := GetTaskWithTagsById(queries,userId, taskId,ctx)
    if errTasks != nil {
        return nil, errTasks    
    }

    return taskWithTags, nil
}

// ---- DELETE ----
func DeleteTask(queries *db.Queries,  taskId int32, ctx context.Context) error {
    err := queries.DeleteTask(ctx, db.DeleteTaskParams{
        ID: ,
        taskId: 
    })
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return ErrNotFound
        }
        return fmt.Errorf("%w: %v", ErrBadRequest, err)
    }

    return nil
}

func DeleteTagRelationFromTask(queries *db.Queries, tagId int32, taskId int32, ctx context.Context) (*models.TaskWithTags, error) {
    errRelations := queries.DeleteTagTaskRelation(ctx, db.DeleteTagTaskRelationParams{
        TaskID: taskId,
        TagID:  tagId})
    
		if errRelations != nil {
        if errors.Is(errRelations, pgx.ErrNoRows) {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("%w: %v", ErrBadRequest, errRelations)
    }

    taskWithTags, err := GetTaskWithTagsById(queries, taskId,ctx)
    if err != nil {
        return nil, err
    }

    return taskWithTags, nil
}