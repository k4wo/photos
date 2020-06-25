package db

import (
	"testing"
)

func TestGetFiles(t *testing.T) {
	userID := 17
	files, err := GetFiles(userID, db)

	if err != nil || len(files) != 39 {
		t.Errorf("GetFiles = %d; want `%d`", len(files), 39)
	}
}
func TestDeleteFiles(t *testing.T) {
	userID := 16

	files := DeleteFiles([]int{523, 525, 778, 808}, userID, db)
	allFiles, err := GetFiles(userID, db)
	if len(files) != 0 || err != nil || len(allFiles) != 37 {
		t.Errorf("GetFiles = %d; want `%d`, user is the owner", len(files), 0)
	}

	files = DeleteFiles([]int{432, 43, 98}, userID, db)
	if len(files) != 3 {
		t.Errorf("GetFiles = %d; want `%d`, user isn't the owner", len(files), 3)
	}
}
