package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func deleteFileRoute(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	enableCors(&w)
	const userID = 1 // TODO: use real userID
	type Payload struct {
		ID []int `json:"id"`
	}
	var payload Payload
	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		panic(err)
	}

	err = deleteFiles(payload.ID, userID)
	if err != nil {
		fmt.Println("deleteFile", err)
	}

	w.WriteHeader(http.StatusOK)
}
