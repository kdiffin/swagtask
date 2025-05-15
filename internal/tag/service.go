package tag

import (
	"context"
	"fmt"
	"swagtask/backend/models"
	db "swagtask/internal/db/generated"
	"swagtask/internal/utils"

	"github.com/jackc/pgx/v5/pgtype"
)

func getTagWithTasksById(queries *db.Queries, tagId, userId pgtype.UUID,
	ctx context.Context) (*models.TagWithTasks, error) {
	tagsWithTaskRelations, err := queries.GetTagWithTaskRelations(ctx, db.GetTagWithTaskRelationsParams{
		ID:     tagId,
		UserID: userId,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", utils.ErrBadRequest, err)
	}
	if len(tagsWithTaskRelations) == 0 {
		return nil, utils.ErrNotFound
	}

	var tag db.Tag
	relatedTasks := []models.RelatedTask{}
	for _, tagWithTaskRelation := range tagsWithTaskRelations {
		tag = db.Tag{
			ID:        tagWithTaskRelation.ID,
			Name:      tagWithTaskRelation.Name,
			CreatedAt: tagWithTaskRelation.CreatedAt,
			UpdatedAt: tagWithTaskRelation.UpdatedAt,
		}
		if tagWithTaskRelation.TaskID.Valid && tagWithTaskRelation.TaskName.Valid {
			relatedTasks = append(relatedTasks, models.RelatedTask{
				Name: tagWithTaskRelation.TaskName.String,
				ID:   tagWithTaskRelation.TaskID,
			})
		}
	}

	allTaskOptions, errGettingAllTasks := queries.GetAllTaskOptions(ctx, userId)
	if errGettingAllTasks != nil {
		return nil, fmt.Errorf("%w: %v", utils.ErrBadRequest, errGettingAllTasks)
	}
	availableTags := getTagAvailableTasks(allTaskOptions, relatedTasks)
	tagWithTasks := models.NewTagWithTasks(tag, relatedTasks, availableTags)

	return &tagWithTasks, nil
}

// func getTagsWithTasks(queries *db.Queries, userId pgtype.UUID, ctx context.Context) ([]models.TagWithTasks, error) {
// 	tagsWithTasksRelations, errTags := queries.GetTagsWithTaskRelations(ctx, userId)
// 	if errTags != nil {
// 		return nil, fmt.Errorf("%w: %v", ErrBadRequest, errTags)
// 	}

// 	allTaskOptions, errAllTaskOptions := queries.GetAllTaskOptions(ctx, userId)
// 	if errAllTaskOptions != nil {
// 		return nil, fmt.Errorf("%w: %v", ErrBadRequest, errAllTaskOptions)
// 	}

// 	tagsWithTasks := []models.TagWithTasks{}
// 	tagIdToTaskOptions := make(map[int32][]models.RelatedTask)
// 	idToTag := make(map[int32]db.Tag)
// 	orderedIds := []int32{}
// 	idSeen := make(map[int32]bool)
// 	for _, tag := range tagsWithTasksRelations {
// 		if tag.TaskID.Valid && tag.TaskName.Valid {
// 			tagIdToTaskOptions[tag.ID] = append(tagIdToTaskOptions[tag.ID], models.RelatedTask{ID: tag.TaskID.Int32, Name: tag.TaskName.String})
// 		}
// 		if !idSeen[tag.ID] {
// 			orderedIds = append(orderedIds, tag.ID)
// 		}
// 		idToTag[tag.ID] = db.Tag{
// 			ID:        tag.ID,
// 			Name:      tag.Name,
// 			CreatedAt: tag.CreatedAt,
// 			UpdatedAt: tag.UpdatedAt,
// 			UserID:    tag.UserID,
// 		}
// 		idSeen[tag.ID] = true
// 	}

// 	for _, id := range orderedIds {
// 		tag := idToTag[id]
// 		tagsOfTask := tagIdToTaskOptions[id]
// 		avaialbleTags := getTagAvailableTasks(allTaskOptions, tagsOfTask)

// 		tagWithTasks := models.NewTagWithTasks(tag, tagsOfTask, avaialbleTags)
// 		tagsWithTasks = append(tagsWithTasks, tagWithTasks)
// 	}

// 	return tagsWithTasks, nil
// }

// func updateTag(queries *db.Queries, tagId int32, userId pgtype.UUID,
// 	tagName string, ctx context.Context) (*models.TagWithTasks, error) {
// 	err := queries.UpdateTag(ctx, db.UpdateTagParams{
// 		Name:   tagName,
// 		ID:     tagId,
// 		UserID: userId,
// 	})

// 	if err != nil {
// 		return nil, fmt.Errorf("%w: %v", ErrBadRequest, err)
// 	}

// 	tagWithTasks, errTags := GetTagWithTasksById(queries, tagId, userId, ctx)
// 	if errTags != nil {
// 		return nil, errTags
// 	}

// 	return tagWithTasks, nil
// }

// func deleteTag(queries *db.Queries, tagId int32,
// 	userId pgtype.UUID, ctx context.Context) error {
// 	errDelete := queries.DeleteTag(ctx, db.DeleteTagParams{
// 		ID:     tagId,
// 		UserID: userId,
// 	})
// 	if errDelete != nil {

// 		return fmt.Errorf("%w: %v", ErrBadRequest, errDelete)
// 	}

// 	return nil
// }

// func deleteTaskRelationFromTag(queries *db.Queries, tagId int32,
// 	userId pgtype.UUID, taskId int32, ctx context.Context) (*models.TagWithTasks, error) {
// 	err := queries.DeleteTagTaskRelation(ctx, db.DeleteTagTaskRelationParams{
// 		TaskID: taskId,
// 		TagID:  tagId,
// 	})

// 	if err != nil {
// 		return nil, fmt.Errorf("%w: %v", ErrBadRequest, err)
// 	}

// 	tagWithTasks, errTag := GetTagWithTasksById(queries, tagId, userId, ctx)
// 	if errTag != nil {
// 		return nil, err
// 	}

// 	return tagWithTasks, nil
// }

// func addTaskToTag(queries *db.Queries, tagId int32,
// 	userId pgtype.UUID, taskId int32, ctx context.Context) (*models.TagWithTasks, error) {
// 	err := queries.CreateTagTaskRelation(ctx, db.CreateTagTaskRelationParams{
// 		TaskID: taskId,
// 		TagID:  tagId,
// 	})

// 	if err != nil {
// 		return nil, fmt.Errorf("%w: %v", ErrBadRequest, err)
// 	}

// 	tagWithTasks, errTag := GetTagWithTasksById(queries, tagId, userId, ctx)
// 	if errTag != nil {
// 		return nil, errTag
// 	}

// 	return tagWithTasks, nil
// }
