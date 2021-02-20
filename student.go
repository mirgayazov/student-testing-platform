package main

import (
	"fmt"
	"net/http"
	"html/template"
	"database/sql"
)

func studentPanel(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/studentPanel.html","templates/header.html","templates/footer.html")	

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	var info Info
	info.UserName = getUserName(r)
	info.UserStatus = getUserStatus(r)
	info.UserPosition = getUserPosition(r)

	t.ExecuteTemplate(w, "studentPanel", info)	
}


var tasks =[]Task{}
var ids =[]uint16{}

func testing(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/testing.html","templates/header.html","templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	
	connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	
	res, err := db.Query("SELECT id, question, answer FROM public.questions")
	if err != nil {
		panic(err)
	}
	defer res.Close()

	tasks = []Task{}
	ids = []uint16{}
	for res.Next() {
		var task Task
		err = res.Scan(&task.ID, &task.Question, &task.Answer)
		if err != nil {
			panic(err)
		}
		tasks = append(tasks, task)
		ids = append(ids, task.ID)
		//создал массив id-шников осталось рандомом достать 3 штуки и по id достать вопрос + ответ
	}
	fmt.Println(len(tasks))

	t.ExecuteTemplate(w, "testing", tasks)
}

// func saveQuestion(w  http.ResponseWriter, r *http.Request) {
// 	question := r.FormValue("question")
// 	answer := r.FormValue("answer")

// 	if question == "" || answer == ""{
// 		fmt.Fprintf(w, "Не все данные заполнены")
// 	} else {
// 		connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
// 		db, err := sql.Open("postgres", connStr)
// 		if err != nil {
// 			panic(err)
// 		}
// 		defer db.Close()
		
// 		insert, err := db.Query(fmt.Sprintf("INSERT INTO questions (question, answer) VALUES('%s','%s')", question, answer))
// 		if err != nil {
// 			panic(err)
// 		}
// 		defer insert.Close()
// 		http.Redirect(w, r, "/questionsСontrol", http.StatusSeeOther)
// 	}
// }