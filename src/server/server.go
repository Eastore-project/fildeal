package server

import (
	"fmt"
	"net/http"
	"time"
)

func StartServer(port int, handler http.Handler) {
	server := &http.Server{
		Handler:      handler,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: 5 * time.Minute,
		ReadTimeout:  5 * time.Minute,
	}
	fmt.Printf("Server listening on port %d\n", port)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
