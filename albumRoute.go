package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	appDB "photos/db"
)

func setAlbumCoverRoute(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	enableCors(&w)
	albumID := p.ByName("id")

	type Payload struct {
		File int `json:"file"`
	}
	var payload Payload
	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil || payload.File == 0 {
		fmt.Println("setAlbumCoverRoute", err)
		jsonResponse(w, http.StatusBadRequest, "")
		return
	}

	status, fileStruct := appDB.SetAlbumCover(albumID, userID, payload.File, db)
	file, err := json.Marshal(&fileStruct)

	if err != nil {
		fmt.Println("setAlbumCoverRoute", err)
		jsonResponse(w, http.StatusInternalServerError, "")
	} else {
		jsonResponse(w, status, string(file))
	}
}

func addFilesToAlbumRoute(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	enableCors(&w)
	albumID := p.ByName("id")

	type Payload struct {
		Files []int `json:"files"`
	}
	var payload Payload
	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		fmt.Println("addFilesToAlbumRoute", err)
		jsonResponse(w, http.StatusBadRequest, "")
		return
	}

	status := appDB.AddFilesToAlbum(albumID, userID, payload.Files, db)
	jsonResponse(w, status, "")
}

func fetchAlbumContentRoute(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	enableCors(&w)
	albumID := p.ByName("id")
	files, err := appDB.GetAlbumContent(userID, albumID, db)

	if err != nil {
		fmt.Println("fetchAlbumContentRoute", err)
		jsonResponse(w, http.StatusInternalServerError, "")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func fetchAlbumsRoute(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	enableCors(&w)
	album, err := appDB.GetAlbums(userID, db)

	if err != nil {
		fmt.Println("fetchAlbumsRoute", err)
		jsonResponse(w, http.StatusInternalServerError, "")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(album)
}

func removeFromAlbumRoute(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	enableCors(&w)
	albumID := p.ByName("id")

	type Payload struct {
		File int `json:"file"`
	}
	var payload Payload
	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil || payload.File == 0 {
		fmt.Println("removeFromAlbumRoute", err)
		jsonResponse(w, http.StatusBadRequest, "")
		return
	}

	status := appDB.RemoveFromAlbum(albumID, userID, payload.File, db)
	jsonResponse(w, status, "")
}

func addNewAlbumRoute(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	enableCors(&w)
	type Payload struct {
		Name string `json:"name"`
	}
	var payload Payload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		fmt.Println("addNewAlbum", err)
	}

	album, err := appDB.CreateAlbum(userID, payload.Name, db)
	if err != nil {
		fmt.Println("addNewAlbum", err)
	}

	json.NewEncoder(w).Encode(album)
}

func deleteAlbumRoute(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	enableCors(&w)
	albumID := p.ByName("id")

	err := appDB.DeleteAlbum(albumID, userID, db)

	if err != nil {
		fmt.Println("deleteAlbum", err)
	}

	w.WriteHeader(http.StatusOK)
}
