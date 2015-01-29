package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func MainController(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func ViewController(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/html/go_tutorial_1_how_to_install_go.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	t.Execute(w, nil)
}

func main() {
	// 静态文件
	// http://127.0.0.1:8888/static/css/default.css
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// 主页
	// http://127.0.0.1:8888
	http.HandleFunc("/", MainController)

	// html页面
	// http://127.0.0.1:8888/view
	http.HandleFunc("/view", ViewController)
	http.ListenAndServe(":8888", nil)
}
