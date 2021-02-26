package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"goblog/bootstrap"
	"goblog/pkg/database"
	"goblog/pkg/logger"
	"net/http"
	"strconv"
	"strings"
)

var router *mux.Router
var db *sql.DB

// GetRouteVariable 获取 URI 路由参数
func getRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}

type Article struct {
	Title, Body string
	ID          int64
}

func (article Article) Delete() (rowAffected int64, err error) {
	res, err := db.Exec("DELETE FROM articles WHERE id=" + strconv.FormatInt(article.ID, 10))

	if err != nil {
		return 0, err
	}

	if n, _ := res.RowsAffected(); n > 0 {
		return n, nil
	}
	return 0, nil
}

func forceHtmlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. 设置标头
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// 2. 继续处理请求
		next.ServeHTTP(w, r)
	})
}

func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}

		next.ServeHTTP(w, r)
	})
}

func getArticleByID(id string) (Article, error) {
	article := Article{}
	query := "SELECT * FROM articles WHERE id=?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)

	return article, err

}

func articlesDELETEHandler(w http.ResponseWriter, r *http.Request) {
	id := getRouteVariable("id", r)
	article, err := getArticleByID(id)

	if err != nil {
		if err == sql.ErrNoRows {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		rowAffected, err := article.Delete()

		if err != nil {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "服务器内部错误")
		} else {
			if rowAffected > 0 {
				indexURL, _ := router.Get("articles.index").URL()
				http.Redirect(w, r, indexURL.String(), http.StatusFound)

			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404 文章未找到")
			}
		}
	}
}

func main() {
	database.Initialize()
	db = database.DB

	bootstrap.SetupDB()
	router = bootstrap.SetupRoute()

	router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDELETEHandler).Methods("POST").Name("articles.delete")

	router.Use(forceHtmlMiddleware)

	// 通过路由名称获取路由
	homeUrl, _ := router.Get("about").URL()
	fmt.Println(homeUrl)
	articlesUrl, _ := router.Get("articles.show").URL("id", "123")
	fmt.Println(articlesUrl)

	http.ListenAndServe(":3000", removeTrailingSlash(router))

	fmt.Printf("CCTV")
}
