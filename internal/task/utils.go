package task

func getTaskAvailableTags(allTags []availableTag, relatedTags []relatedTag) []availableTag {
	// think of this as a set checking if the tag is a tag of the task
	// int is id
	tagExists := make(map[string]bool)
	for _, tag := range relatedTags {
		tagExists[tag.ID] = true
	}

	availableTags := []availableTag{}
	for _, tag := range allTags {
		if !tagExists[tag.ID] {
			availableTags = append(availableTags, tag)
		}
	}

	return availableTags
}
