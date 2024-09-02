package server

import (
	configurations "fildeal/src/config"
	"fmt"
	"net/http"
	"time"
)


func StartServer(config configurations.Configurations, handler http.Handler) {
    server := &http.Server{
        Handler:      handler,
        Addr:         fmt.Sprintf("0.0.0.0:%d", config.Port), // Listen on all network interfaces
        WriteTimeout: 5 * time.Minute,
        ReadTimeout:  5 * time.Minute,
    }

    fmt.Printf("Server listening on port %d\n", config.Port)
    err := server.ListenAndServe()
    if err != nil {
        fmt.Println("Error starting server:", err)
    }
}