package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/guregu/null.v3"
)

// Album is representation of db album table
type Album struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Size      int      `json:"size"`
	Owner     int      `json:"owner"`
	Cover     null.Int `json:"cover"`
	UpdatedAt string   `json:"updatedAt"`
	CreatedAt string   `json:"createdAt"`
}

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

func (store *dbStruct) getAlbumContent(userID int, albumID string) ([]ImageInfo, error) {
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
	rows, _ := store.connection.Query(rawQuery, albumID, userID)
	defer rows.Close()

	return filesScanner(rows)
}

func (store *dbStruct) getAlbum(name string, userID int) (Album, error) {
	query := selectAlbum + " WHERE owner = $1 AND name = $2"
	row := store.connection.QueryRow(query, userID, name)

	return albumScanner(row)
}

func (store *dbStruct) getAlbums() ([]Album, error) {
	rows, _ := store.connection.Query(selectAlbum)
	defer rows.Close()

	return albumsScanner(rows)
}

func (store *dbStruct) createAlbum(name string) (Album, error) {
	if name == "" {
		msg := fmt.Sprintf(STRINGS["noAlbumName"])
		return Album{}, errors.New(msg)
	}

	const userID = 1 // TODO: use real userID

	album, err := db.getAlbum(name, userID)
	if err != sql.ErrNoRows {
		return album, err
	}

	query := `INSERT INTO albums(owner, name) VALUES($1, $2) RETURNING *`
	row := store.connection.QueryRow(
		query,
		userID,
		name,
	)

	return albumScanner(row)
}

func fetchAlbums(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	enableCors(&w)
	album, err := db.getAlbums()
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

	album, err := db.createAlbum(payload.Name)
	if err != nil {
		fmt.Println("addNewAlbum", err)
	}

	json.NewEncoder(w).Encode(album)
}

func getAlbumContent(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	enableCors(&w)
	albumID := p.ByName("id")
	files, err := db.getAlbumContent(2, albumID)

	if err != nil {
		fmt.Println("getAlbumContent", err)
	}

	json.NewEncoder(w).Encode(files)
}
