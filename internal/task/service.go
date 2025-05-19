package task

import (
	"context"
	"errors"
	"fmt"
	db "swagtask/internal/db/generated"
	"swagtask/internal/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// ---- READ ----

func GetTaskWithTagsById(queries *db.Queries, userId, id pgtype.UUID, ctx context.Context) (*taskWithTags, error) {
	taskWithRelations, err := queries.GetTaskWithTagRelations(ctx, db.GetTaskWithTagRelationsParams{
		ID:     id,
		UserID: userId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", utils.ErrBadRequest, err)
	}
	if len(taskWithRelations) == 0 {
		return nil, utils.ErrNotFound
	}
	var task taskUI
	tagsOfTask := []relatedTag{}

	allTags, errTags := queries.getAllTagsDesc(ctx, userId)
	if errTags != nil {
		if errors.Is(errTags, pgx.ErrNoRows) {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", utils.ErrBadRequest, errTags)
	}

	for _, taskWithRelation := range taskWithRelations {
		task = taskUI{
			ID:        taskWithRelation.ID.String(),
			Name:      taskWithRelation.Name,
			Idea:      taskWithRelation.Idea,
			Completed: taskWithRelation.Completed,
		}

		if taskWithRelation.TagID.Valid && taskWithRelation.TagName.Valid {
			tagsOfTask = append(tagsOfTask, relatedTag{ID: taskWithRelation.TagID.Int32, Name: taskWithRelation.TagName.String})
		}
	}

	availableTags := getTaskAvailableTags(allTags, tagsOfTask)
	taskWithTags := newTaskWithTags(task, tagsOfTask, availableTags)
	return &taskWithTags, nil
}

func GetFilteredTasksWithTags(queries *db.Queries, filters *tasksPageFilters,
	userId pgtype.UUID, ctx context.Context) ([]taskWithTags, error) {
	taskswithTagRelations, err := queries.GetFilteredTasks(ctx, db.GetFilteredTasksParams{
		TaskName: utils.StringToPgText(filters.SearchQuery),
		TagName:  utils.StringToPgText(filters.ActiveTag),
		UserID:   userId,
	})
	if err != nil {
		fmt.Println("error here at first")
		return nil, fmt.Errorf("%w: %v", utils.ErrBadRequest, err)
	}

	allTags, errAllTags := queries.GetAllTagsDesc(ctx, userId)
	if errAllTags != nil {
		fmt.Println("error here at second")
		return nil, fmt.Errorf("%w: %v", utils.ErrBadRequest, errAllTags)
	}

	tasksWithTags := []taskWithTags{}
	taskIdToTags := make(map[pgtype.UUID][]relatedTag)
	idToTask := make(map[pgtype.UUID]db.Task)
	orderedIds := []pgtype.UUID{}
	idSeen := make(map[pgtype.UUID]bool)
	for _, task := range taskswithTagRelations {
		if task.TagID.Valid && task.TagName.Valid && task.TagUserID.Valid {
			taskIdToTags[task.ID] = append(taskIdToTags[task.ID],
				relatedTag{ID: task.TagID.String(),
					Name: task.TagName.String})
		}

		if !idSeen[task.ID] {
			orderedIds = append(orderedIds, task.ID)
		}

		idToTask[task.ID] = db.Task{
			ID:        task.ID,
			Name:      task.Name,
			Idea:      task.Idea,
			Completed: task.Completed,
			UserID:    task.UserID,
			CreatedAt: task.CreatedAt,
			UpdatedAt: task.UpdatedAt,
		}

		idSeen[task.ID] = true
	}

	for _, id := range orderedIds {
		task := idToTask[id]
		tagsOfTask := taskIdToTags[id]
		avaialbleTags := getTaskAvailableTags(allTags, tagsOfTask)

		taskWithTags := newTaskWithTags(newUITask(task), tagsOfTask, avaialbleTags)
		tasksWithTags = append(tasksWithTags, taskWithTags)
	}

	return tasksWithTags, nil
}

func GetTaskNavigationButtons(ctx context.Context, queries *db.Queries, userId pgtype.UUID, id pgtype.UUID) (taskButton, taskButton) {
	prev, errPrev := queries.GetPreviousTaskDetails(ctx, db.GetPreviousTaskDetailsParams{
		ID:     id,
		UserID: userId,
	})
	next, errNext := queries.GetNextTaskDetails(ctx, db.GetNextTaskDetailsParams{
		ID:     id,
		UserID: userId,
	})

	var prevButton, nextButton taskButton
	if errPrev == nil {
		prevButton = taskButton{ID: prev.ID, Name: prev.Name, Exists: true}
	}
	if errNext == nil {
		nextButton = taskButton{ID: next.ID, Name: next.Name, Exists: true}
	}
	return prevButton, nextButton
}

// ---- CREATE ----
func CreateTask(queries *db.Queries, name string,
	userId pgtype.UUID, idea string, ctx context.Context) (*taskWithTags, error) {
	task, errCreate := queries.CreateTask(ctx, db.CreateTaskParams{Name: name, Idea: idea, UserID: userId})
	if errCreate != nil {
		if errors.Is(errCreate, pgx.ErrNoRows) {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", utils.ErrUnprocessable, errCreate)
	}

	allTags, errAllTags := queries.GetAllTagsDesc(ctx, userId)
	if errAllTags != nil {
		if errors.Is(errAllTags, pgx.ErrNoRows) {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", utils.ErrBadRequest, errAllTags)
	}

	taskWithTag := newTaskWithTags(
		newUITask(task),
		[]relatedTag{},
		allTags,
	)

	return &taskWithTag, nil
}

// ---- UPDATE ----
func UpdateTaskCompletion(queries *db.Queries, userId pgtype.UUID, taskId pgtype.UUID, ctx context.Context) (*taskWithTags, error) {
	errCompletion := queries.ToggleTaskCompletion(ctx, db.ToggleTaskCompletionParams{
		ID:     taskId,
		UserID: userId,
	})
	if errCompletion != nil {
		if errors.Is(errCompletion, pgx.ErrNoRows) {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", utils.ErrBadRequest, errCompletion)
	}

	taskWithTags, errGetTask := GetTaskWithTagsById(queries, userId, taskId, ctx)
	if errGetTask != nil {
		return nil, errGetTask
	}

	return taskWithTags, nil
}

func UpdateTask(queries *db.Queries, taskId pgtype.UUID, userId pgtype.UUID,
	name string, idea string, ctx context.Context) (*taskWithTags, error) {
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
		ID:     taskId,
		Name:   nameNullable,
		UserID: userId,
		Idea:   ideaNullable,
	})
	if errCompletion != nil {
		if errors.Is(errCompletion, pgx.ErrNoRows) {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", utils.ErrUnprocessable, errCompletion)
	}

	taskWithTags, errGetTask := GetTaskWithTagsById(queries, userId, taskId, ctx)
	if errGetTask != nil {
		return nil, errGetTask
	}

	return taskWithTags, nil
}

func AddTagToTask(queries *db.Queries,
	tagId, userId, taskId pgtype.UUID, ctx context.Context) (*taskWithTags, error) {
	err := queries.CreateTagTaskRelation(ctx, db.CreateTagTaskRelationParams{
		TaskID: taskId,
		TagID:  tagId,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", utils.ErrBadRequest, err)
	}

	taskWithTags, errTasks := GetTaskWithTagsById(queries, userId, taskId, ctx)
	if errTasks != nil {
		return nil, errTasks
	}

	return taskWithTags, nil
}

// ---- DELETE ----
func DeleteTask(queries *db.Queries, taskId, userId pgtype.UUID, ctx context.Context) error {
	err := queries.DeleteTask(ctx, db.DeleteTaskParams{
		ID:     taskId,
		UserID: userId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return utils.ErrNotFound
		}
		return fmt.Errorf("%w: %v", utils.ErrBadRequest, err)
	}

	return nil
}

func DeleteTagRelationFromTask(queries *db.Queries, tagId, userId, taskId pgtype.UUID, ctx context.Context) (*taskWithTags, error) {
	errRelations := queries.DeleteTagTaskRelation(ctx, db.DeleteTagTaskRelationParams{
		TaskID: taskId,
		TagID:  tagId})

	if errRelations != nil {
		if errors.Is(errRelations, pgx.ErrNoRows) {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", utils.ErrBadRequest, errRelations)
	}

	taskWithTags, err := GetTaskWithTagsById(queries, userId, taskId, ctx)
	if err != nil {
		return nil, err
	}

	return taskWithTags, nil
}
