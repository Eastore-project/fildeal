package routes

import (
	configurations "fildeal/src/config"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)


func DownloadRouter() *mux.Router {
	router := mux.NewRouter().PathPrefix("/download").Subrouter()    // Create a sub-router for '/download' prefix
	router.HandleFunc("/car", DownloadCarHead).Methods("HEAD") // Adjusted to '/car'
	router.HandleFunc("/car", DownloadCar).Methods("GET")      // Adjusted to '/car'
	return router
}


func DownloadCarHead(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file_name")
	filePath := filepath.Join(configurations.LoadConfigurations().AggregateCarPath, fileName)

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		fmt.Println("failed to get file info: %w", err)
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
	filePath := filepath.Join(configurations.LoadConfigurations().AggregateCarPath, fileName)
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
	rangeHeader := r.Header.Get("Range") // Example: "bytes=0-999"

	if rangeHeader == "" {
		// If no range is specified, send the whole file
		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", fileSize))
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		io.Copy(w, file)
		return
	}

	// Parse range header to get byte start and end
	var start, end int64
	_, err = fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)
	if err != nil || start < 0 || end < start || start >= fileSize {
		// If the range is invalid or cannot be parsed, return 416 Range Not Satisfiable
		http.Error(w, "Range Not Satisfiable", http.StatusRequestedRangeNotSatisfiable)
		return
	}

	// If the end is beyond the file size, adjust it to the file size
	if end >= fileSize || end == 0 {
		end = fileSize - 1
	}

	// Set headers for partial content
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", end-start+1))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusPartialContent)

	// Write the requested range to response
	file.Seek(start, 0)
	io.CopyN(w, file, end-start+1)
}
