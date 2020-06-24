package db

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"testing"

	constants "photos/constants"
	dev "photos/dev"

	_ "github.com/lib/pq"
	"github.com/subosito/gotenv"
)

var db *sql.DB

func TestGetAlbumContent(t *testing.T) {
	userID := 10

	albumID := "56"
	expected := 13
	albumContent, err := GetAlbumContent(userID, albumID, db)
	if len(albumContent) != expected || err != nil {
		t.Errorf(
			"GetAlbumContent - owned album returns %d expect %d - error: %s",
			len(albumContent),
			expected,
			err,
		)
	}

	albumID = "1"
	expected = 10
	albumContent, err = GetAlbumContent(userID, albumID, db)
	if len(albumContent) != expected || err != nil {
		t.Errorf(
			"GetAlbumContent - shared album returns %d expect %d - error: %s",
			len(albumContent),
			expected,
			err,
		)
	}

	albumID = "1"
	albumContent, err = GetAlbumContent(userID, albumID, db)
	if err != nil && err.Error() != constants.STRINGS["noAccessToAlbum"] {
		t.Errorf("GetAlbumContent - don't have access to the album - error: %s", err)
	}
}

func TestGetAlbums(t *testing.T) {
	userID := 1
	expected := 2
	albums, _ := GetAlbums(userID, db)

	if len(albums) != expected {
		t.Errorf("GetAlbums(%d) = %d; want %d", userID, len(albums), expected)
	}
}

func TestCreateAlbum(t *testing.T) {
	albumName := "album name - test"
	album, err := CreateAlbum(20, albumName, db)

	if err != nil || album.Name != albumName {
		t.Errorf("CreateAlbum(20, %s) = %s; want `%s`", albumName, album.Name, albumName)
	}
}

func TestAddFilesToAlbum(t *testing.T) {
	userID := 9
	filesID := []int{4, 994, 572}

	albumID := "88"
	status := AddFilesToAlbum(albumID, userID, filesID, db)
	albumContent, _ := GetAlbumContent(userID, albumID, db)
	if status != http.StatusOK || len(albumContent) != len(filesID) {
		t.Errorf(
			"AddFilesToAlbum - %d, expected %d - user adds to his album",
			len(albumContent),
			len(filesID),
		)
	}

	albumID = "14"
	status = AddFilesToAlbum(albumID, userID, filesID, db)
	albumContent, _ = GetAlbumContent(userID, albumID, db)
	if status != http.StatusOK || len(albumContent) != 11 {
		t.Errorf(
			"AddFilesToAlbum - %d, expected %d - user adds to shared album",
			len(albumContent),
			len(filesID),
		)
	}

	albumID = "26"
	status = AddFilesToAlbum(albumID, userID, filesID, db)
	albumContent, _ = GetAlbumContent(userID, albumID, db)
	if status != http.StatusOK || len(albumContent) != 12 {
		t.Errorf(
			"AddFilesToAlbum - %d, expected %d - user adds file which already is in the album",
			len(albumContent),
			12,
		)
	}

	albumID = "1"
	status = AddFilesToAlbum(albumID, userID, filesID, db)
	if status != http.StatusForbidden {
		t.Errorf(
			"AddFilesToAlbum - status: %d, expected %d - user doesn't have access to the album",
			status,
			http.StatusForbidden,
		)
	}

	albumID = "7"
	filesID = []int{1, 2, 3}
	status = AddFilesToAlbum(albumID, userID, filesID, db)
	albumContent, _ = GetAlbumContent(userID, albumID, db)
	if status != http.StatusForbidden || len(albumContent) != 13 {
		t.Errorf(
			"AddFilesToAlbum - %d, expected %d - user doesn't have access to the files",
			len(albumContent),
			13,
		)
	}
}

func TestSetAlbumCover(t *testing.T) {
	userID := 15
	fileID := 306
	albumID := "28"
	_, file := SetAlbumCover(albumID, userID, fileID, db)
	if file.ID.ValueOrZero() != int64(fileID) {
		t.Errorf(
			"SetAlbumCover - %d, expected %d - owner adds a file", file.ID.ValueOrZero(), fileID,
		)
	}

	userID = 9
	fileID = 421
	albumID = "69"
	_, file = SetAlbumCover(albumID, userID, fileID, db)
	if file.ID.ValueOrZero() != int64(fileID) {
		t.Errorf(
			"SetAlbumCover - %d, expected %d - user and shared album", file.ID.ValueOrZero(), fileID,
		)
	}

	fileID = 35
	_, file = SetAlbumCover(albumID, userID, fileID, db)
	if file.ID.ValueOrZero() == int64(fileID) {
		t.Errorf(
			"SetAlbumCover - %d, expected %d - no file in the album", file.ID.ValueOrZero(), fileID,
		)
	}
}

func TestRemoveFromAlbum(t *testing.T) {
	userID := 6

	albumID := "60"
	fileID := 111
	status := RemoveFromAlbum(albumID, userID, fileID, db)
	isInAlbum := isFileInAlbum(fileID, albumID, db)
	if status != http.StatusOK || isInAlbum == true {
		t.Errorf("RemoveFromAlbum - remove shared file %d from the album %s", fileID, albumID)
	}

	albumID = "99"
	fileID = 391
	status = RemoveFromAlbum(albumID, userID, fileID, db)
	isInAlbum = isFileInAlbum(fileID, albumID, db)
	if status != http.StatusOK || isInAlbum == true {
		t.Errorf("RemoveFromAlbum - remove owned file %d from the album %s", fileID, albumID)
	}

	albumID = "44"
	fileID = 998
	status = RemoveFromAlbum(albumID, userID, fileID, db)
	isInAlbum = isFileInAlbum(fileID, albumID, db)
	if status == http.StatusOK || isInAlbum == false {
		t.Errorf("RemoveFromAlbum - user don't have access the album %s", albumID)
	}
}

func TestDeleteAlbum(t *testing.T) {
	userID := 8

	albumID := 26
	isAlbumAvailable := false
	albumIDString := fmt.Sprintf("%d", albumID)
	err := DeleteAlbum(albumIDString, userID, db)
	albums, albumsErr := GetAlbums(userID, db)
	files, _ := GetAlbumContent(userID, albumIDString, db)
	for _, album := range albums {
		if album.ID == albumID {
			isAlbumAvailable = true
		}
	}
	if err != nil || albumsErr != nil || len(files) > 0 || isAlbumAvailable == true {
		t.Errorf("DeleteAlbumTest - delete album %d by owner %d", albumID, userID)
	}

	albumID = 52
	isAlbumAvailable = false
	err = DeleteAlbum(albumIDString, userID, db)
	albums, albumsErr = GetAlbums(userID, db)
	for _, album := range albums {
		if album.ID == albumID {
			isAlbumAvailable = true
		}
	}
	if err == nil || isAlbumAvailable == true {
		t.Errorf("DeleteAlbumTest - album %d by user who is not the owner %d", albumID, userID)
	}
}

func TestMain(m *testing.M) {
	gotenv.Load("../.env_test")
	dbConfig := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	database, err := sql.Open("postgres", dbConfig)
	if err != nil {
		panic(err)
	}
	db = database
	dev.ResetDatabase(database)

	os.Exit(m.Run())
}
