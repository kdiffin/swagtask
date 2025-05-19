package tag

import (
	"context"
	"fmt"
	db "swagtask/internal/db/generated"
	"swagtask/internal/utils"

	"github.com/jackc/pgx/v5/pgtype"
)

func getTagWithTasksById(queries *db.Queries, tagId, userId, vaultId pgtype.UUID,
	ctx context.Context) (*tagWithTasks, error) {
	tagsWithTaskRelations, err := queries.GetTagWithTaskRelations(ctx, db.GetTagWithTaskRelationsParams{
		ID:      tagId,
		UserID:  userId,
		VaultID: vaultId,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", utils.ErrBadRequest, err)
	}
	if len(tagsWithTaskRelations) == 0 {
		return nil, utils.ErrNotFound
	}

	var tagUIStruct tagUI
	relatedTasks := []relatedTask{}
	for _, tagWithTaskRelation := range tagsWithTaskRelations {
		tagUIStruct.Name = tagWithTaskRelation.Name
		tagUIStruct.ID = tagWithTaskRelation.ID.String()
		tagUIStruct.Author = tagAuthor{
			PathToPfp: tagWithTaskRelation.AuthorPathToPfp,
			Username:  tagWithTaskRelation.AuthorUsername,
		}
		tagUIStruct.VaultID = tagWithTaskRelation.VaultID.String()

		if tagWithTaskRelation.TaskID.Valid && tagWithTaskRelation.TaskName.Valid {
			relatedTasks = append(relatedTasks, relatedTask{
				Name: tagWithTaskRelation.TaskName.String,
				ID:   tagWithTaskRelation.TaskID.String(),
			})
		}
	}

	allTaskOptions, errGettingAllTasks := queries.GetAllTaskOptions(ctx, userId)
	if errGettingAllTasks != nil {
		return nil, fmt.Errorf("%w: %v", utils.ErrBadRequest, errGettingAllTasks)
	}
	availableTags := getTagAvailableTasks(allTaskOptions, relatedTasks)
	tagWithTasks := newTagWithTasks(tagUIStruct, relatedTasks, availableTags)

	return &tagWithTasks, nil
}

func GetTagsWithTasks(queries *db.Queries, userId, vaultId pgtype.UUID, ctx context.Context) ([]tagWithTasks, error) {

	tagsWithTasksRelations, errTags := queries.GetTagsWithTaskRelations(ctx, db.GetTagsWithTaskRelationsParams{
		UserID:  userId,
		VaultID: vaultId,
	})
	if errTags != nil {
		return nil, fmt.Errorf("%w: %v", utils.ErrBadRequest, errTags)
	}

	allTaskOptions, errAllTaskOptions := queries.GetAllTaskOptions(ctx, userId)
	if errAllTaskOptions != nil {
		return nil, fmt.Errorf("%w: %v", utils.ErrBadRequest, errAllTaskOptions)
	}

	tagsWithTasks := []tagWithTasks{}
	tagIdToTaskOptions := make(map[pgtype.UUID][]relatedTask)
	idToTag := make(map[pgtype.UUID]db.GetTagsWithTaskRelationsRow)
	orderedIds := []pgtype.UUID{}
	idSeen := make(map[pgtype.UUID]bool)
	for _, tag := range tagsWithTasksRelations {

		if tag.TaskID.Valid && tag.TaskName.Valid {
			tagIdToTaskOptions[tag.ID] = append(tagIdToTaskOptions[tag.ID], relatedTask{ID: tag.TaskID.String(), Name: tag.TaskName.String})
		}
		if !idSeen[tag.ID] {
			orderedIds = append(orderedIds, tag.ID)
		}
		idToTag[tag.ID] = tag
		idSeen[tag.ID] = true
	}

	for _, id := range orderedIds {
		tag := idToTag[id]
		tagsOfTask := tagIdToTaskOptions[id]
		avaialbleTags := getTagAvailableTasks(allTaskOptions, tagsOfTask)

		tagWithTasks := newTagWithTasks(tagUI{
			VaultID: tag.VaultID.String(),
			Name:    tag.Name,
			Author: tagAuthor{
				PathToPfp: tag.AuthorPathToPfp,
				Username:  tag.AuthorUsername,
			},
			ID: tag.ID.String(),
		}, tagsOfTask, avaialbleTags)
		tagsWithTasks = append(tagsWithTasks, tagWithTasks)
	}

	return tagsWithTasks, nil
}

func updateTag(queries *db.Queries, tagId, userId, vaultId pgtype.UUID,
	tagName string, ctx context.Context) (*tagWithTasks, error) {
	err := queries.UpdateTag(ctx, db.UpdateTagParams{
		VaultID: vaultId,
		Name:    tagName,
		ID:      tagId,
		UserID:  userId,
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %v", utils.ErrBadRequest, err)
	}

	tagWithTasks, errTags := getTagWithTasksById(queries, tagId, userId, vaultId, ctx)
	if errTags != nil {

		return nil, errTags
	}

	return tagWithTasks, nil
}

func deleteTag(queries *db.Queries, tagId, userId, vaultId pgtype.UUID, ctx context.Context) error {
	errDelete := queries.DeleteTag(ctx, db.DeleteTagParams{
		ID:      tagId,
		VaultID: vaultId,
		UserID:  userId,
	})
	if errDelete != nil {
		return fmt.Errorf("%w: %v", utils.ErrBadRequest, errDelete)
	}

	return nil
}

func deleteTaskRelationFromTag(queries *db.Queries, tagId, taskId, userId, vaultId pgtype.UUID, ctx context.Context) (*tagWithTasks, error) {
	err := queries.DeleteTagTaskRelation(ctx, db.DeleteTagTaskRelationParams{
		TaskID: taskId,
		TagID:  tagId,
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %v", utils.ErrBadRequest, err)
	}

	tagWithTasks, errTag := getTagWithTasksById(queries, tagId, userId, vaultId, ctx)
	if errTag != nil {
		return nil, errTag
	}

	return tagWithTasks, nil
}

func addTaskToTag(queries *db.Queries, tagId, userId, taskId, vaultId pgtype.UUID, ctx context.Context) (*tagWithTasks, error) {
	err := queries.CreateTagTaskRelation(ctx, db.CreateTagTaskRelationParams{
		TaskID: taskId,
		TagID:  tagId,
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %v", utils.ErrBadRequest, err)
	}

	tagWithTasks, errTag := getTagWithTasksById(queries, tagId, userId, vaultId, ctx)
	if errTag != nil {
		return nil, errTag
	}

	return tagWithTasks, nil
}
