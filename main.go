package main

import (
	"fmt"
	"net/http"
	_ "github.com/lib/pq"
	"./handlers"
)



func handleFunc() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("/create", handlers.Create)
	http.HandleFunc("/save_article", handlers.SaveArticle)
	http.HandleFunc("/registration", handlers.Registration)
	http.HandleFunc("/authorization", handlers.Authorization)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/logout", handlers.Logout)
	http.HandleFunc("/about", handlers.About)
	http.HandleFunc("/saveUser", handlers.SaveUser)
	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("Переход -> http://localhost:8080/")
	handleFunc()
}