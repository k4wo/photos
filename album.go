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
		id,
		owner,
		name,
		size,
		cover,
		updated_at,
		created_at
	FROM albums
`

func getAlbumContent(userID int, albumID string) ([]model.File, error) {
	rawQuery := `
		SELECT
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
			AND user_album. "user" = $2;
	`

	fmt.Println(userID, albumID)
	rows, _ := db.Query(rawQuery, albumID, userID)
	defer rows.Close()

	return filesScanner(rows)
}

func getAlbum(name string, userID int) (model.Album, error) {
	query := selectAlbum + " WHERE owner = $1 AND name = $2"
	row := db.QueryRow(query, userID, name)

	return albumScanner(row)
}

func getAlbums() ([]model.Album, error) {
	rows, _ := db.Query(selectAlbum)
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

func fetchAlbums(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	enableCors(&w)
	album, err := getAlbums()
	if err != nil {
		fmt.Println("fetchAlbums", err)
	}

	json.NewEncoder(w).Encode(album)
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

func fetchAlbumContent(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	enableCors(&w)
	albumID := p.ByName("id")
	files, err := getAlbumContent(2, albumID)

	if err != nil {
		fmt.Println("getAlbumContent", err)
	}

	json.NewEncoder(w).Encode(files)
}
