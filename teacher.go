package main

import (
	"fmt"
	"net/http"
	"html/template"
	"database/sql"
)

func teacherPanel(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/teacherPanel.html","templates/header.html","templates/footer.html")	

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	var info Info
	info.UserName = getUserName(r)
	info.UserStatus = getUserStatus(r)
	info.UserPosition = getUserPosition(r)

	t.ExecuteTemplate(w, "teacherPanel", info)	
}

func questionsСontrol(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/questionsСontrol.html","templates/header.html","templates/footer.html")	

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	var info Info
	info.UserName = getUserName(r)
	info.UserStatus = getUserStatus(r)
	info.UserPosition = getUserPosition(r)

	t.ExecuteTemplate(w, "questionsСontrol", info)
}

func saveQuestion(w  http.ResponseWriter, r *http.Request) {
	question := r.FormValue("question")
	answer := r.FormValue("answer")

	if question == "" || answer == ""{
		fmt.Fprintf(w, "Не все данные заполнены")
	} else {
		connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
		defer db.Close()
		
		insert, err := db.Query(fmt.Sprintf("INSERT INTO questions (question, answer) VALUES('%s','%s')", question, answer))
		if err != nil {
			panic(err)
		}
		defer insert.Close()
		http.Redirect(w, r, "/questionsСontrol", http.StatusSeeOther)
	}
}

func createCourse(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/createCourse.html","templates/header.html","templates/footer.html")	

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	var info Info
	info.UserName = getUserName(r)
	info.UserStatus = getUserStatus(r)
	info.UserPosition = getUserPosition(r)

	t.ExecuteTemplate(w, "createCourse", info)
}

func requestToСreateСourse(w http.ResponseWriter, r *http.Request) {
	connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	
	insert, err := db.Query(fmt.Sprintf("INSERT INTO course_requests (teacher_name, course_name) VALUES('%s','%s')", getUserName(r), r.FormValue("courseName")))
	if err != nil {
		panic(err)
	} else {
		message := "Вы успешно отправили заявку!"
		t, err := template.ParseFiles("templates/message.html","templates/footer.html")	
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		t.ExecuteTemplate(w, "message", message)
	}
	defer insert.Close()	
	
}