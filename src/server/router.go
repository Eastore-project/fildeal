package server

import (
	"fmt"
	"net/http"

	"github.com/eastore-project/fildeal/src/routes"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	// Middleware
	router.Use(loggingMiddleware)

	// Routes
	router.PathPrefix("/download").Handler(routes.DownloadRouter())

	// 404 Handler
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	return router
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.RemoteAddr + " - " + r.Method + " - " + r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("404 Not Found:", r.RequestURI)
	http.Error(w, "Endpoint does not exist.", http.StatusNotFound)
}
