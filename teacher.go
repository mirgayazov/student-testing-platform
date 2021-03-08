package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	s "strings"

	// "encoding/json"
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

	connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query(fmt.Sprintf("SELECT course_name FROM courses where id='%s'",id))
	if err != nil {
		panic(err)
	}
	var courseName string
	for res.Next() {
		res.Scan(&courseName)
	}
	defer res.Close()

	var info Info
	info.UserName = getUserName(r)
	info.UserStatus = getUserStatus(r)
	info.UserPosition = getUserPosition(r)

	res, err = db.Query(fmt.Sprintf("SELECT array_to_json(subscribers) FROM courses where id='%s'",id))
	if err != nil {
		panic(err)
	}
	var arrStr string
	for res.Next() {
		res.Scan(&arrStr)
		arrStr = s.Replace(arrStr, "[", "", 1)
		arrStr = s.Replace(arrStr, "]", "", 1)
		// fmt.Println(arrStr)
	}
	defer res.Close()
	

	students :=[]Student{}
	identificators := strings.Split(arrStr, ",")
	for _, identificator := range identificators {
		// fmt.Println(identificator)
		res, err := db.Query(fmt.Sprintf("SELECT last_name, first_name FROM users where id='%s'",identificator))
		if err != nil {
			panic(err)
		}
		for res.Next() {
			var student Student
			res.Scan(&student.LastName, &student.FirstName)
			students = append(students, student)
		}
		defer res.Close()
	}
	count := len(students)
	//00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
	res, err = db.Query(fmt.Sprintf("SELECT id, question, answer FROM questions where course_id='%s'",id))
	if err != nil {
		panic(err)
	}
	defer res.Close()

	indx :=1
	tasks :=[]Task{}
	for res.Next() {
		var task Task
		err = res.Scan(&task.ID, &task.Question, &task.Answer)
		if err != nil {
			panic(err)
		}
		task.Index=uint16(indx)
		indx++
		tasks = append(tasks, task)
	}
	taskCount := len(tasks)

	t, err := template.ParseFiles("templates/header.html","templates/coursePage.html","templates/footer.html")	
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "coursePage", struct{Info, Students, Count, Tasks, TaskCount, CourseName interface{}}{info, students, count, tasks, taskCount, courseName});
}