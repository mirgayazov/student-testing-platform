package main

import (
	"fmt"
	"net/http"
	_ "github.com/lib/pq"
)

func handleFunc() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	//
	http.HandleFunc("/", index)
	//
	http.HandleFunc("/registration", registration)
	http.HandleFunc("/authorization", authorization)
	//
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	//
	http.HandleFunc("/about", about)
	//
	http.HandleFunc("/adminPanel", adminPanel)
	http.HandleFunc("/accessСontrol", accessСontrol)
	http.HandleFunc("/saveUserPosition", saveUserPosition)
	//
	http.HandleFunc("/teacherPanel", teacherPanel)
	http.HandleFunc("/questionsСontrol", questionsСontrol)
	http.HandleFunc("/saveQuestion", saveQuestion)
	//
	http.HandleFunc("/studentPanel", studentPanel)
	http.HandleFunc("/testing", testing)
	http.HandleFunc("/checkAndSaveTest", checkAndSaveTest)
	//
	http.HandleFunc("/saveUser", saveUser)
	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("Переход -> http://localhost:8080/")
	handleFunc()
}