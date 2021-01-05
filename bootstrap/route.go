package bootstrap

import (
	"github.com/gorilla/mux"
	"goblog/pkg/routes"
)

func SetupRoute() *mux.Router {
	router := mux.NewRouter()
	routes.RegisterWebRoutes(router)
	return router
}