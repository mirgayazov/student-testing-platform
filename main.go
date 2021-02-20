package main

import (
	"fmt"
	"net/http"
	_ "github.com/lib/pq"
)

func handleFunc() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/create", create)
	http.HandleFunc("/save_article", saveArticle)
	http.HandleFunc("/registration", registration)
	http.HandleFunc("/authorization", authorization)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/about", about)
	http.HandleFunc("/saveUser", saveUser)
	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("Переход -> http://localhost:8080/")
	handleFunc()
}