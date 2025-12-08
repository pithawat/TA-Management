package main

import (
	"TA-management/internal/logs"
	router "TA-management/internal/routers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	log := logs.InitializeLogger()
	defer logs.SyncLogger(log)

	routes := router.InitRouter()

	server := &http.Server{
		Addr:    ":8084",
		Handler: routes,
	}

	server.ListenAndServe()
}
