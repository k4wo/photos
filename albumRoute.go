package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
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

	status := setAlbumCover(albumID, userID, payload.File)
	jsonResponse(w, status, "")
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

	status := addFilesToAlbum(albumID, userID, payload.Files)
	jsonResponse(w, status, "")
}

func fetchAlbumContentRoute(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	enableCors(&w)
	albumID := p.ByName("id")
	files, err := getAlbumContent(userID, albumID)

	if err != nil {
		fmt.Println("fetchAlbumContentRoute", err)
		jsonResponse(w, http.StatusInternalServerError, "")
		return
	}

	json.NewEncoder(w).Encode(files)
}

func fetchAlbumsRoute(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	enableCors(&w)
	album, err := getAlbums(userID)

	if err != nil {
		fmt.Println("fetchAlbumsRoute", err)
		jsonResponse(w, http.StatusInternalServerError, "")
		return
	}

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

	status := removeFromAlbum(albumID, userID, payload.File)
	jsonResponse(w, status, "")
}
