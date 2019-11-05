package main

import (
	. "fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)


// UploadFile uploads a file to the server
func UploadFile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	file, handle, err := r.FormFile(FileField)
	if file == nil {
		Fprint(w, STRINGS["noFile"])
		return
	}

	if err != nil {
		Fprintf(w, "%v", err)
		return
	}
	defer file.Close()

	mimeType := handle.Header.Get("Content-Type")
	switch mimeType {
	case "image/jpeg":
		saveFile(w, file, handle)
	case "image/png":
		saveFile(w, file, handle)
	default:
		jsonResponse(w, http.StatusBadRequest, STRINGS["fileFormatInvalid"])
	}
}
