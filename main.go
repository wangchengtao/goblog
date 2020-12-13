package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello, 这里是 goblog</h1>")
}
func main() {
	fmt.Printf("中央电视台")
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":3000", nil)

	fmt.Printf("CCTV")
}