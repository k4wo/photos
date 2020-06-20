package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	appDB "photos/db"

	"github.com/julienschmidt/httprouter"
)

func deleteFileRoute(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	enableCors(&w)
	type Payload struct {
		ID []int `json:"id"`
	}
	var payload Payload
	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		panic(err)
	}

	err = appDB.DeleteFiles(payload.ID, userID, db)
	if err != nil {
		fmt.Println("deleteFile", err)
	}

	w.WriteHeader(http.StatusOK)
}

func uploadFilesRoute(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	enableCors(&w)

	r.ParseMultipartForm(32 << 20) // 32MB is the default used by FormFile
	files := r.MultipartForm.File["files"]
	status := appDB.ProcessFiles(files, userID, UploadDir, db)

	w.WriteHeader(status)
}

func fetchFilesRoute(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	enableCors(&w)
	files, err := appDB.GetFiles(userID, db)
	if err != nil {
		fmt.Println("fetchFilesRoute", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	json.NewEncoder(w).Encode(files)
}
