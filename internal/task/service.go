package task

import (
	"context"
	"errors"
	"fmt"
	"swagtask/internal/auth"
	db "swagtask/internal/db/generated"
	"swagtask/internal/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// ---- READ ----

func GetTaskWithTagsById(queries *db.Queries, userId, vaultId, id pgtype.UUID, ctx context.Context) (*TaskWithTags, error) {
	taskWithRelations, err := queries.GetTaskWithTagRelations(ctx, db.GetTaskWithTagRelationsParams{
		ID:      id,
		UserID:  userId,
		VaultID: vaultId,
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
	var task TaskUI
	tagsOfTask := []relatedTag{}

	allTags, errTags := queries.GetAllTagsDesc(ctx, db.GetAllTagsDescParams{
		VaultID: vaultId,
		UserID:  userId,
	})
	if errTags != nil {
		return nil, fmt.Errorf("%w: %v", utils.ErrBadRequest, errTags)
	}

	for _, taskWithRelation := range taskWithRelations {
		task = TaskUI{
			ID:   taskWithRelation.ID.String(),
			Name: taskWithRelation.Name,
			Idea: taskWithRelation.Idea,
			Author: auth.Author{
				PathToPfp: taskWithRelation.AuthorPathToPfp,
				Username:  taskWithRelation.AuthorUsername,
			},
			Completed: taskWithRelation.Completed,
		}

		if taskWithRelation.TagID.Valid && taskWithRelation.TagName.Valid {
			tagsOfTask = append(tagsOfTask, relatedTag{ID: taskWithRelation.TagID.String(), Name: taskWithRelation.TagName.String})
		}
	}

	availableTags := getTaskAvailableTags(allTags, tagsOfTask)
	taskWithTags := newTaskWithTags(task, tagsOfTask, availableTags)
	return &taskWithTags, nil
}

func GetFilteredTasksWithTags(queries *db.Queries, filters TasksPageFilters,
	userId, vaultId pgtype.UUID, ctx context.Context) ([]TaskWithTags, error) {
	taskswithTagRelations, err := queries.GetFilteredTasks(ctx, db.GetFilteredTasksParams{
		TaskName: utils.StringToPgText(filters.SearchQuery),
		TagName:  utils.StringToPgText(filters.ActiveTag),
		UserID:   userId,
		VaultID:  vaultId,
	})
	if err != nil {
		fmt.Println("error here at first")
		return nil, fmt.Errorf("%w: %v", utils.ErrBadRequest, err)
	}

	allTags, errAllTags := queries.GetAllTagsDesc(ctx, db.GetAllTagsDescParams{
		VaultID: vaultId,
		UserID:  userId,
	})
	if errAllTags != nil {
		fmt.Println("error here at second")
		return nil, fmt.Errorf("%w: %v", utils.ErrBadRequest, errAllTags)
	}

	tasksWithTags := []TaskWithTags{}
	taskIdToTags := make(map[pgtype.UUID][]relatedTag)
	idToTask := make(map[pgtype.UUID]TaskUI)
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

		idToTask[task.ID] = TaskUI{
			ID:        task.ID.String(),
			Name:      task.Name,
			Idea:      task.Idea,
			Completed: task.Completed,
			Author: auth.Author{
				PathToPfp: task.AuthorPathToPfp,
				Username:  task.AuthorUsername,
			}}

		idSeen[task.ID] = true
	}

	for _, id := range orderedIds {
		task := idToTask[id]
		tagsOfTask := taskIdToTags[id]
		avaialbleTags := getTaskAvailableTags(allTags, tagsOfTask)

		taskWithTags := newTaskWithTags(task, tagsOfTask, avaialbleTags)
		tasksWithTags = append(tasksWithTags, taskWithTags)
	}

	return tasksWithTags, nil
}

func GetTaskNavigationButtons(ctx context.Context, queries *db.Queries, createdAt pgtype.Timestamp,
	userId, vaultId, id pgtype.UUID) (taskButton, taskButton) {
	prev, errPrev := queries.GetPreviousTaskDetails(ctx, db.GetPreviousTaskDetailsParams{
		CreatedAt: createdAt,
		VaultID:   vaultId,
		UserID:    userId,
	})
	next, errNext := queries.GetNextTaskDetails(ctx, db.GetNextTaskDetailsParams{
		UserID:    userId,
		CreatedAt: createdAt,
		VaultID:   vaultId,
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
func CreateTask(queries *db.Queries, name, idea string,
	userId, vaultId pgtype.UUID, ctx context.Context) (*TaskWithTags, error) {
	fmt.Println(userId.String())
	fmt.Println(vaultId.String())
	fmt.Println(name)
	fmt.Println(idea)

	id, errCreate := queries.CreateTask(ctx, db.CreateTaskParams{
		Name:    name,
		Idea:    idea,
		UserID:  userId,
		VaultID: vaultId,
	})
	if errCreate != nil {
		return nil, fmt.Errorf("%w: %v", utils.ErrUnprocessable, errCreate)
	}
	fmt.Println(id)

	taskWithTags, errGetTask := GetTaskWithTagsById(queries, userId, vaultId, id, ctx)
	if errGetTask != nil {
		return nil, errGetTask
	}

	return taskWithTags, nil
}

// ---- UPDATE ----
func UpdateTaskCompletion(queries *db.Queries, userId, vaultId pgtype.UUID, taskId pgtype.UUID, ctx context.Context) (*TaskWithTags, error) {
	errCompletion := queries.ToggleTaskCompletion(ctx, db.ToggleTaskCompletionParams{
		ID:      taskId,
		UserID:  userId,
		VaultID: vaultId,
	})
	if errCompletion != nil {
		if errors.Is(errCompletion, pgx.ErrNoRows) {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", utils.ErrBadRequest, errCompletion)
	}

	taskWithTags, errGetTask := GetTaskWithTagsById(queries, userId, vaultId, taskId, ctx)
	if errGetTask != nil {
		return nil, errGetTask
	}

	return taskWithTags, nil
}

func UpdateTask(queries *db.Queries, vaultId, taskId pgtype.UUID, userId pgtype.UUID,
	name string, idea string, ctx context.Context) (*TaskWithTags, error) {
	namePg := utils.StringToPgText(name)
	ideaPg := utils.StringToPgText(idea)
	if !namePg.Valid && !ideaPg.Valid {
		return nil, utils.ErrNoUpdateFields
	}

	errCompletion := queries.UpdateTask(ctx, db.UpdateTaskParams{
		Name:    namePg,
		Idea:    ideaPg,
		ID:      taskId,
		UserID:  userId,
		VaultID: vaultId,
	})
	if errCompletion != nil {
		return nil, fmt.Errorf("%w: %v", utils.ErrUnprocessable, errCompletion)
	}

	taskWithTags, errGetTask := GetTaskWithTagsById(queries, userId, vaultId, taskId, ctx)
	if errGetTask != nil {
		return nil, errGetTask
	}

	return taskWithTags, nil
}

func addTagToTask(queries *db.Queries,
	tagId, userId, taskId, vaultId pgtype.UUID, ctx context.Context) (*TaskWithTags, error) {
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

	taskWithTags, errTasks := GetTaskWithTagsById(queries, userId, vaultId, taskId, ctx)
	if errTasks != nil {
		return nil, errTasks
	}

	return taskWithTags, nil
}

// ---- DELETE ----
func DeleteTask(queries *db.Queries, taskId, vaultId, userId pgtype.UUID, ctx context.Context) error {
	err := queries.DeleteTask(ctx, db.DeleteTaskParams{
		ID:      taskId,
		UserID:  userId,
		VaultID: vaultId,
	})
	if err != nil {
		return fmt.Errorf("%w: %v", utils.ErrBadRequest, err)
	}

	return nil
}

func deleteTagRelationFromTask(queries *db.Queries, tagId, userId, vaultId, taskId pgtype.UUID, ctx context.Context) (*TaskWithTags, error) {
	errRelations := queries.DeleteTagTaskRelation(ctx, db.DeleteTagTaskRelationParams{
		TaskID: taskId,
		TagID:  tagId})

	if errRelations != nil {
		if errors.Is(errRelations, pgx.ErrNoRows) {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", utils.ErrBadRequest, errRelations)
	}

	taskWithTags, err := GetTaskWithTagsById(queries, userId, vaultId, taskId, ctx)
	if err != nil {
		return nil, err
	}

	return taskWithTags, nil
}

// this is a one off for the task page
func GetTaskPage(queries *db.Queries, userId, vaultId, id pgtype.UUID, ctx context.Context) (*TaskWithTags, pgtype.Timestamp, error) {
	taskWithRelations, err := queries.GetTaskWithTagRelations(ctx, db.GetTaskWithTagRelationsParams{
		ID:      id,
		UserID:  userId,
		VaultID: vaultId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgtype.Timestamp{}, utils.ErrNotFound
		}
		return nil, pgtype.Timestamp{}, fmt.Errorf("%w: %v", utils.ErrBadRequest, err)
	}
	if len(taskWithRelations) == 0 {
		return nil, pgtype.Timestamp{}, utils.ErrNotFound
	}
	var task TaskUI
	tagsOfTask := []relatedTag{}

	allTags, errTags := queries.GetAllTagsDesc(ctx, db.GetAllTagsDescParams{
		VaultID: vaultId,
		UserID:  userId,
	})
	if errTags != nil {
		return nil, pgtype.Timestamp{}, fmt.Errorf("%w: %v", utils.ErrBadRequest, errTags)
	}

	var createdAt pgtype.Timestamp
	for _, taskWithRelation := range taskWithRelations {
		task = TaskUI{
			ID:   taskWithRelation.ID.String(),
			Name: taskWithRelation.Name,
			Idea: taskWithRelation.Idea,
			Author: auth.Author{
				PathToPfp: taskWithRelation.AuthorPathToPfp,
				Username:  taskWithRelation.AuthorUsername,
			},
			Completed: taskWithRelation.Completed,
		}
		createdAt = taskWithRelation.CreatedAt

		if taskWithRelation.TagID.Valid && taskWithRelation.TagName.Valid {
			tagsOfTask = append(tagsOfTask, relatedTag{ID: taskWithRelation.TagID.String(), Name: taskWithRelation.TagName.String})
		}
	}

	availableTags := getTaskAvailableTags(allTags, tagsOfTask)
	taskWithTags := newTaskWithTags(task, tagsOfTask, availableTags)
	return &taskWithTags, createdAt, nil
}
