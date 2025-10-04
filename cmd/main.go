package main

import (
	router "TA-management/internal/routers"
	"net/http"
)

func main() {
	routes := router.InitRouter()

	server := &http.Server{
		Addr:    ":8084",
		Handler: routes,
	}

	server.ListenAndServe()
}
