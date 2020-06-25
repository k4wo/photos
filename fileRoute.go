package main

import (
	"encoding/json"
	"net/http"

	appDB "photos/db"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

func deleteFileRoute(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	enableCors(&w)
	type Payload struct {
		ID []int `json:"id"`
	}
	var payload Payload
	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		log.Error().Err(err).Caller().Int("user", userID).Msg("Can't parse files' id to delete")

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	notDeleted := appDB.DeleteFiles(payload.ID, userID, db)
	if len(notDeleted) > 0 {
		response, err := json.Marshal(Payload{notDeleted})
		if err != nil {
			log.Error().Err(err).Caller().Int("user", userID).Msg("Can't parse not deleted files")
			response, _ = json.Marshal(Payload{})
		}

		jsonResponse(w, http.StatusUnprocessableEntity, string(response))
		return
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
		log.Error().Err(err).Caller().Int("user", userID).Msg("Can't parse files")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(files)
}
