package database

func GetTaskAvailableTags(allTags []Tag, relatedTags []Tag) []Tag {
	// think of this as a set checking if the tag is a tag of the task
	// int is id
	tagExists := make(map[int]bool)
	for _, tag := range relatedTags {
		tagExists[tag.Id] = true
	}

	availableTags := []Tag{}
	for _, tag := range allTags {
		if !tagExists[tag.Id] {
			availableTags = append(availableTags, tag)
		}
	}

	return availableTags
}
