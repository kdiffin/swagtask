package service

import (
	"context"
	"fmt"
	db "swagtask/db/generated"
	"swagtask/models"
)

// ---- READ ----

func GetTaskWithTagsById(queries *db.Queries, id int32) (*models.TaskWithTags, error) {
	taskWithRelations, err := queries.GetTaskWithTagRelations(context.Background(),id)
	if err != nil {
		return nil,err
	}
	var task models.Task
	tagsOfTask := []db.Tag{}

	allTags, errTags := queries.GetAllTagsDesc(context.Background())
	if errTags != nil {
		return nil, errTags
	}
	
	for _,taskWithRelation := range taskWithRelations {
		task = models.Task{
			ID: taskWithRelation.ID,
			Name: taskWithRelation.Name,
			Idea: taskWithRelation.Idea,
			Completed: taskWithRelation.Completed.Bool,
		}

		if taskWithRelation.TagID.Valid && taskWithRelation.TagName.Valid {
			tagsOfTask = append(tagsOfTask, db.Tag{ID: taskWithRelation.TagID.Int32, Name: taskWithRelation.TagName.String  })
		} 
	}

	availableTags := GetTaskAvailableTags(allTags, tagsOfTask)
	taskWithTags := models.NewTaskWithTags(task, tagsOfTask, availableTags)
	return &taskWithTags, nil
}

func GetAllTasksWithTags(queries *db.Queries) ([]models.TaskWithTags, error) {
	// 1. gets the task with their related tags
	// 2. gets all tags
	// 3. finds the non related tags for the dropdown and then sends them
	// TODO: make these run in paralell
	taskswithTagRelations, errTasks := queries.GetTasksWithTagRelations(context.Background())
	if errTasks != nil {
		return nil, fmt.Errorf("service.GetAllTasksWithTags: %w", errTasks)
	}

	allTags, errAllTags := queries.GetAllTagsDesc(context.Background())
	if errAllTags != nil {
		return nil, fmt.Errorf("service.GetAllTasksWithTags: %w", errAllTags)
	}

	tasksWithTags := []models.TaskWithTags{}
	taskIdToTags := make(map[int32][]db.Tag)
	// for easy lookup
	idToTask := make(map[int32]db.Task)
	// so that we can make the items ordered (this is kindof a set)
	orderedIds := []int32{}
	// so that we can mimic the set logic
	idSeen := make(map[int32]bool)
	for _, task := range taskswithTagRelations {
		
		if task.TagID.Valid && task.TagName.Valid {
			taskIdToTags[task.ID] = append(taskIdToTags[task.ID], db.Tag{ID: task.TagID.Int32, Name: task.TagName.String})

		}

		// mimic a set
		if !idSeen[task.ID] {
			orderedIds = append(orderedIds, task.ID)
		}

		idToTask[task.ID] = db.Task{
			ID: task.ID,
			Name: task.Name,
			Idea: task.Idea,
			Completed: task.Completed,
		}

		idSeen[task.ID] = true
	}

	for _, id := range orderedIds {
		task := idToTask[id]
		tagsOfTask := taskIdToTags[id]
		avaialbleTags := GetTaskAvailableTags(allTags, tagsOfTask)

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

// // ---- CREATE ----
func CreateTask(queries *db.Queries, name string, idea string) (*models.TaskWithTags, error) {
	task, errCreate := queries.CreateTask(context.Background(), db.CreateTaskParams{Name: name, Idea: idea})
	if errCreate != nil {
		return nil, fmt.Errorf("error creating task: %w", errCreate)
	}

	allTags, errAllTags := queries.GetAllTagsDesc(context.Background())
	if errAllTags != nil {
		return nil, fmt.Errorf("error getting all tags: %w", errAllTags)
	}

	taskWithTag := models.NewTaskWithTags(
		models.NewUITask(task),
		[]db.Tag{},
		allTags,
	)

	return &taskWithTag, nil
}

//  ---- UPDATE ---
func UpdateTaskCompletion(queries *db.Queries, taskId int32) (*models.TaskWithTags, error) {
	errCompletion := queries.ToggleTaskCompletion(context.Background(), taskId)
	if errCompletion != nil {
		return nil, errCompletion
	}

	taskWithTags, errGetTask  := GetTaskWithTagsById(queries, taskId)
	if errGetTask != nil {
		return nil, fmt.Errorf("err getting task after updating: %w", errGetTask )
	}
	
	return taskWithTags, nil
}

func AddTagToTask(queries *db.Queries, tagId int32, taskId int32)  (*models.TaskWithTags, error) {
	err := queries.CreateTagTaskRelation(context.Background(), db.CreateTagTaskRelationParams{
		TaskID: Int32ToInt4Psql(taskId),
		TagID: Int32ToInt4Psql(tagId),
	}) 
	if err != nil {
		return nil , err 
	}
		
	// get updated task
	taskWithTags, errTasks := GetTaskWithTagsById(queries, taskId)
	if errTasks != nil {
		return nil, errTasks
	}

	return taskWithTags, nil
}


// DELETE
func DeleteTask(queries *db.Queries, taskId int32) error {
	errRelations := queries.DeleteAllTagRelationsForTask(context.Background(), Int32ToInt4Psql(taskId))
	if errRelations != nil {
		return errRelations
	}
	
	err := queries.DeleteTask(context.Background(), taskId)
	if err != nil {
		return err
	}

	return nil
}

func DeleteTagRelationFromTask(queries *db.Queries, tagId int32, taskId int32) (*models.TaskWithTags, error) {
	errRelations := queries.DeleteSingleTagRelation(context.Background(), db.DeleteSingleTagRelationParams{
		TaskID: Int32ToInt4Psql(taskId),
		TagID: Int32ToInt4Psql(tagId)})
			
	if errRelations != nil {
		return  nil,errRelations
	}
	
	
	taskWithTags, err := GetTaskWithTagsById(queries, taskId)
	if err != nil {
		return nil,err
	}

	return taskWithTags, nil
}