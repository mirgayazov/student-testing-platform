package main

import (
	"fmt"
	"net/http"
	"html/template"
	s "strings"
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
	
	insert, err := db.Query(fmt.Sprintf("INSERT INTO course_requests (teacher_name, course_name, codeword) VALUES('%s','%s', '%s')", getUserName(r), r.FormValue("courseName"), r.FormValue("codeword")))
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

func teacherCourses(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/teacherCourses.html","templates/header.html","templates/footer.html")	
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var info Info
	info.UserName = getUserName(r)
	info.UserStatus = getUserStatus(r)
	info.UserPosition = getUserPosition(r)

	res, err := db.Query(fmt.Sprintf("SELECT id, course_name FROM courses where teacher_name='%s'",info.UserName))
	if err != nil {
		panic(err)
	}
	courses := []Course{}
	for res.Next() {
		var course Course
		err = res.Scan(&course.ID, &course.CourseName)
		if err != nil {
			panic(err)
		}
		courses = append(courses, course)
	}
	defer res.Close()

	t.ExecuteTemplate(w, "teacherCourses", struct{Info, Course interface{}}{info, courses});
}

func course(w http.ResponseWriter, r *http.Request) {
	id := s.Replace(fmt.Sprint(r.URL), "/course/", "", -1)
	message := id
	
	t, err := template.ParseFiles("templates/message.html","templates/footer.html")	
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "message", message)
}