package routes

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

// AggregateCarPath will be set with the cli flag value.
var AggregateCarPath string

func DownloadRouter() *mux.Router {
	router := mux.NewRouter().PathPrefix("/download").Subrouter()
	router.HandleFunc("/car", DownloadCarHead).Methods("HEAD")
	router.HandleFunc("/car", DownloadCar).Methods("GET")
	return router
}

func DownloadCarHead(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file_name")
	filePath := filepath.Join(AggregateCarPath, fileName)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		fmt.Println("failed to get file info:", err)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Length", fmt.Sprint(fileInfo.Size()))
	w.Header().Set("Last-Modified", fileInfo.ModTime().UTC().Format(http.TimeFormat))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
}

func DownloadCar(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file_name")
	filePath := filepath.Join(AggregateCarPath, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found.", http.StatusNotFound)
		return
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Error getting file info.", http.StatusInternalServerError)
		return
	}
	fileSize := fileInfo.Size()
	rangeHeader := r.Header.Get("Range")
	if rangeHeader == "" {
		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", fileSize))
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		io.Copy(w, file)
		return
	}
	var start, end int64
	_, err = fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)
	if err != nil || start < 0 || end < start || start >= fileSize {
		http.Error(w, "Range Not Satisfiable", http.StatusRequestedRangeNotSatisfiable)
		return
	}
	if end >= fileSize || end == 0 {
		end = fileSize - 1
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", end-start+1))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusPartialContent)
	file.Seek(start, 0)
	io.CopyN(w, file, end-start+1)
}
