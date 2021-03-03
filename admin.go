package main

import (
	"fmt"
	"net/http"
	"html/template"
	"database/sql"
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
	fmt.Println(r.URL)
}
	