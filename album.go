package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	model "photos/model"

	"github.com/julienschmidt/httprouter"
)

var selectAlbum = `
	SELECT
		albums.id,
		albums.owner,
		albums.name,
		albums.size,
		albums.updated_at,
		albums.created_at,

		files.id,
		files.owner,
		files.name,
		files.hash,
		files.size,
		files.extension,
		files.mime,
		files.latitude,
		files.longitude,
		files.orientation,
		files.model,
		files.camera,
		files.iso,
		files.focal_length,
		files.exposure_time,
		files.f_number,
		files.height,
		files.width,
		files.date
	FROM albums
	LEFT JOIN files ON albums.cover = files.id
	WHERE albums.owner = $1
`

func hasAlbumAccess(userID int, albumID string) bool {
	var count int
	rawQuery := `SELECT count(id) FROM albums WHERE owner = $1 AND id = $2`

	row := db.QueryRow(rawQuery, userID, albumID)
	err := row.Scan(&count)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return count > 0
}

func isFileInAlbum(fileID int, albumID string) bool {
	var count int
	rawQuery := `SELECT count(id) FROM album_file WHERE file = $1 AND album = $2`

	row := db.QueryRow(rawQuery, fileID, albumID)
	err := row.Scan(&count)
	if err != nil {
		fmt.Println(err)
		return true
	}

	return count > 0
}

func getAlbumContent(userID int, albumID string) ([]model.File, error) {
	rawQuery := `
		SELECT
			files.id,
			files.owner,
			files.name,
			files.hash,
			files.size,
			files.extension,
			files.mime,
			files.latitude,
			files.longitude,
			files.orientation,
			files.model,
			files.camera,
			files.iso,
			files.focal_length,
			files.exposure_time,
			files.f_number,
			files.height,
			files.width,
			files.date
		FROM
			album_file
			LEFT JOIN files ON files.id = album_file.file
			LEFT JOIN user_album ON user_album.album = album_file.album
		WHERE
			album_file. "album" = $1
			AND album_file. "user" = $2;
	`

	rows, _ := db.Query(rawQuery, albumID, userID)
	defer rows.Close()

	return filesScanner(rows)
}

func getAlbum(name string, userID int) (model.Album, error) {
	query := selectAlbum + " AND name = $2"
	row := db.QueryRow(query, userID, name)

	return albumScanner(row)
}

func getAlbums(userID int) ([]model.Album, error) {
	rows, _ := db.Query(selectAlbum, userID)
	defer rows.Close()

	return albumsScanner(rows)
}

func createAlbum(name string) (model.Album, error) {
	if name == "" {
		msg := fmt.Sprintf(STRINGS["noAlbumName"])
		return model.Album{}, errors.New(msg)
	}

	const userID = 1 // TODO: use real userID

	album, err := getAlbum(name, userID)
	if err != sql.ErrNoRows {
		return album, err
	}

	query := `INSERT INTO albums(owner, name) VALUES($1, $2) RETURNING *`
	row := db.QueryRow(
		query,
		userID,
		name,
	)

	return albumScanner(row)
}

func deleteAlbum(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	enableCors(&w)
	const userID = 1 // TODO: use real userID
	albumID := p.ByName("id")

	query := `DELETE FROM albums WHERE id = $1 AND owner = $2`
	_, err := db.Exec(
		query,
		albumID,
		userID,
	)

	if err != nil {
		fmt.Println("deleteAlbum", err)
	}

	w.WriteHeader(http.StatusOK)
}

func addNewAlbum(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	enableCors(&w)
	type Payload struct {
		Name string `json:"name"`
	}
	var payload Payload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		fmt.Println("addNewAlbum", err)
	}

	album, err := createAlbum(payload.Name)
	if err != nil {
		fmt.Println("addNewAlbum", err)
	}

	json.NewEncoder(w).Encode(album)
}

func addFilesToAlbum(albumID string, userID int, files []int) int {
	hasAccess := hasAlbumAccess(userID, albumID)
	if !hasAccess {
		return http.StatusForbidden
	}

	for _, fileID := range files {
		if !hasFileAccess(userID, fileID) {
			fmt.Println("addFilesToAlbum", "no access", fileID)

			return http.StatusForbidden
		}
	}

	for _, fileID := range files {
		if !isFileInAlbum(fileID, albumID) {
			rawQuery := `INSERT INTO "album_file"("album", "file", "user") VALUES($1, $2, $3);`
			_, err := db.Exec(
				rawQuery,
				albumID,
				fileID,
				userID,
			)

			if err != nil {
				fmt.Println("addFilesToAlbum", err)
			}
		}
	}

	return http.StatusOK
}

func setAlbumCover(albumID string, userID, fileID int) (int, model.File) {
	hasAccess := hasAlbumAccess(userID, albumID)
	if !hasAccess {
		return http.StatusForbidden, model.File{}
	}

	if !hasFileAccess(userID, fileID) {
		fmt.Println("setAlbumCover", "no access", fileID)
		return http.StatusForbidden, model.File{}
	}

	if isFileInAlbum(fileID, albumID) {
		rawQuery := `UPDATE albums SET cover = $1 WHERE id = $2 RETURNING *`
		db.QueryRow(
			rawQuery,
			fileID,
			albumID,
		)

		file, err := getFileByID(fileID)
		if err != nil {
			fmt.Println("setAlbumCover", err)
			return http.StatusInternalServerError, model.File{}
		}

		return http.StatusOK, file
	}

	return http.StatusBadRequest, model.File{}
}

func removeFromAlbum(albumID string, userID, fileID int) int {
	hasAccess := hasAlbumAccess(userID, albumID)
	if !hasAccess {
		return http.StatusForbidden
	}

	rawQuery := `DELETE FROM "album_file" WHERE file = $1;`
	db.QueryRow(rawQuery, fileID)

	return http.StatusOK
}
