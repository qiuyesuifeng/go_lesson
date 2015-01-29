package main

import (
	"net/http"
)

func MainController(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func main() {
	http.HandleFunc("/", MainController)
	http.ListenAndServe(":8888", nil)
}
