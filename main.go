package main

import (
	"log"
	"net/http"
	"io/ioutil"
	"mime"
	"path/filepath"
	"fmt"
	"os"
	"crypto/rand"
	"strings"
)

const maxUploadSize = 2 * 1024 * 1024 // 2 mb
const uploadPath = "."

func main() {

	fs := http.FileServer(http.Dir(uploadPath))
	http.Handle("/files/", http.StripPrefix("/files", fs))

	http.HandleFunc("/upload", uploadFileHandler())
	http.HandleFunc("/delete/",  deleteFile)

	log.Print("Server started on localhost:8080, use /upload for uploading files and /files/{fileName} for downloading")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func deleteFile(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.RequestURI, "/")
	filename := path[len(path) - 1]
	err := os.Remove(filepath.Join(uploadPath, filename))
	if err != nil {
		renderError(w, err.Error(), http.StatusBadRequest)
	}
	w.Write([]byte(fmt.Sprintf(filename, "has been deleted")))
}

func uploadFileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// validate file size
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			renderError(w, "File Too Big", http.StatusBadRequest)
			return
		}

		// parse and validate file and post parameters

		file, _, err := r.FormFile("uploadFile")
		if err != nil {
			renderError(w, "Invalid File", http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			renderError(w, "Invalid File", http.StatusBadRequest)
			return
		}

		filetype := http.DetectContentType(fileBytes)
		fileName := randToken(6)
		fileEndings, err := mime.ExtensionsByType(filetype)
		if err != nil {
			renderError(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
			return
		}
		newPath := filepath.Join(uploadPath, fileName+fileEndings[0])
		fmt.Printf("File: %s\n", newPath)

		// write file
		newFile, err := os.Create(newPath)
		if err != nil {
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		defer newFile.Close() // idempotent, okay to call twice
		if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(fmt.Sprint(fileName, fileEndings[0])))
	})
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}