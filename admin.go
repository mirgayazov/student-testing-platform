package main

import (
	"fmt"
	"net/http"
	"html/template"
	"database/sql"
	s "strings"
)

func adminPanel(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/adminPanel.html","templates/header.html","templates/footer.html")	

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	var info Info
	info.UserName = getUserName(r)
	info.UserStatus = getUserStatus(r)
	info.UserPosition = getUserPosition(r)

	t.ExecuteTemplate(w, "adminPanel", info)	
}

func accessСontrol(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/accessСontrol.html","templates/header.html","templates/footer.html")	
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	var info Info
	info.UserName = getUserName(r)
	info.UserStatus = getUserStatus(r)
	info.UserPosition = getUserPosition(r)

	t.ExecuteTemplate(w, "accessСontrol", info)
}

func saveUserPosition(w  http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("user_name")
	position := r.FormValue("position")

	if userName == "" || position == ""{
		fmt.Fprintf(w, "Не все данные заполнены")
	} else {
		connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
		defer db.Close()
		
		_, err = db.Query(fmt.Sprintf("UPDATE users SET position = '%s' WHERE user_name = '%s'", position, userName))
		if err != nil {
			panic(err)
		}
		http.Redirect(w, r, "/accessСontrol", http.StatusSeeOther)
	}
}

func adminRequests(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/adminRequests.html","templates/header.html","templates/footer.html")	
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	
	res, err := db.Query("SELECT * FROM course_requests")
	if err != nil {
		panic(err)
	}
	requests := []Request{}
	for res.Next() {
		var req Request
		err = res.Scan(&req.ID, &req.TeacherName, &req.CourseName)
		if err != nil {
			panic(err)
		}
		requests = append(requests, req)
	}

	var info Info
	info.UserName = getUserName(r)
	info.UserStatus = getUserStatus(r)
	info.UserPosition = getUserPosition(r)

	t.ExecuteTemplate(w, "adminRequests", struct{Info, Request interface{}}{info, requests});
}

func approveRequest(w http.ResponseWriter, r *http.Request) {
	id := s.Replace(fmt.Sprint(r.URL), "/approveRequest/", "", -1)
	//по этому id я вытягиваю запрос препод-курс 
	connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	
	res, err := db.Query(fmt.Sprintf("SELECT teacher_name, course_name FROM course_requests where id='%s'",id))
	if err != nil {
		panic(err)
	}
	var req Request
	for res.Next() {
		req.ID = 0
		err = res.Scan(&req.TeacherName, &req.CourseName)
		if err != nil {
			panic(err)
		}
		fmt.Println(req)
	}
	defer res.Close()
	//в таблицу курсов добавляю пару курс-препод
	courseNameWithTeacher :=  req.CourseName+"("+req.TeacherName+")"
	insert, err := db.Query(fmt.Sprintf("INSERT INTO courses (course_name, teacher_name) VALUES('%s','%s')",courseNameWithTeacher, req.TeacherName))
	if err != nil {
		panic(err)
	} else {
		message := fmt.Sprintf("Вы одобрили заявку преподавателю %s на создание курса '%s'.",req.TeacherName, req.CourseName)
		t, err := template.ParseFiles("templates/message.html","templates/footer.html")	
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		db.Query(fmt.Sprintf("delete from course_requests where id='%s'", id))
		t.ExecuteTemplate(w, "message", message)
	}
	defer insert.Close()
}

func rejectRequest(w http.ResponseWriter, r *http.Request) {
	id := s.Replace(fmt.Sprint(r.URL), "/rejectRequest/", "", -1)
	//по этому id я вытягиваю запрос препод-курс 
	connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	res, err := db.Query(fmt.Sprintf("SELECT teacher_name, course_name FROM course_requests where id='%s'",id))
	if err != nil {
		panic(err)
	}
	var req Request
	for res.Next() {
		req.ID = 0
		err = res.Scan(&req.TeacherName, &req.CourseName)
		if err != nil {
			panic(err)
		}
		fmt.Println(req)
	}
	defer res.Close()
	db.Query(fmt.Sprintf("delete from course_requests where id='%s'", id))
	message := fmt.Sprintf("Вы отклонили заявку преподавателя %s на создание курса '%s'.",req.TeacherName, req.CourseName)
		t, err := template.ParseFiles("templates/message.html","templates/footer.html")	
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		t.ExecuteTemplate(w, "message", message)
}
	