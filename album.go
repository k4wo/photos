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
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	Size  int      `json:"size"`
	Owner int      `json:"owner"`
	Cover null.Int `json:"cover"`
}

func (store *dbStruct) getAlbum(name string, userID int) (Album, error) {
	query := `
		SELECT
			id,
			owner,
			name,
			size,
			cover
		FROM albums
		WHERE owner = $1 AND name = $2`

	row := store.connection.QueryRow(query, userID, name)
	album := Album{}
	err := row.Scan(&album.ID, &album.Owner, &album.Name, &album.Size, &album.Cover)

	return album, err
}

func (store *dbStruct) getAlbums() ([]Album, error) {
	query := `
		SELECT
			id,
			owner,
			name,
			size,
			cover
		FROM albums`
	rows, _ := store.connection.Query(query)
	defer rows.Close()

	var albums []Album
	for rows.Next() {
		album := Album{}
		err := rows.Scan(
			&album.ID,
			&album.Owner,
			&album.Name,
			&album.Size,
			&album.Cover,
		)

		if err == nil {
			albums = append(albums, album)
		}

		if err != nil {
			return albums, err
		}
	}

	return albums, nil
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

	query := `INSERT INTO albums(owner, name) VALUES($1, $2)`
	_, err = store.connection.Exec(
		query,
		userID,
		name,
	)

	if err != nil {
		return album, err
	}

	return album, err
}

func fetchAlbums(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	enableCors(&w)
	album, err := db.getAlbums()
	if err != nil {
		fmt.Println("fetchAlbums", err)
	}

	json.NewEncoder(w).Encode(album)
}

func setNewAlbum(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	enableCors(&w)
	type Payload struct {
		Name string `json:"name"`
	}
	var payload Payload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		fmt.Println("setNewAlbum", err)
	}

	album, err := db.createAlbum(payload.Name)
	if err != nil {
		fmt.Println("setNewAlbum", err)
	}

	json.NewEncoder(w).Encode(album)
}
