package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"goblog/app/http/middlewares"
	"goblog/bootstrap"
	"goblog/pkg/database"
	"net/http"
)

var router *mux.Router
var db *sql.DB

func main() {
	database.Initialize()
	db = database.DB

	bootstrap.SetupDB()
	router = bootstrap.SetupRoute()

	// 通过路由名称获取路由
	homeUrl, _ := router.Get("about").URL()
	fmt.Println(homeUrl)
	articlesUrl, _ := router.Get("articles.show").URL("id", "123")
	fmt.Println(articlesUrl)

	http.ListenAndServe(":3000", middlewares.RemoveTrailingSlash(router))

	fmt.Printf("CCTV")
}
