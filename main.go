package main

import (
	"fmt"
	"net/http"
	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
)

func handleFunc() {
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r := mux.NewRouter()
	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
    r.PathPrefix("/static/").Handler(s)
	s = http.StripPrefix("/src/", http.FileServer(http.Dir("./src/")))
	r.PathPrefix("/src/").Handler(s)
	//
	r.HandleFunc("/", index)
	//
	r.HandleFunc("/registration", registration)
	r.HandleFunc("/authorization", authorization)
	//
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	//
	r.HandleFunc("/about", about)
	r.HandleFunc("/tizd", tizd)
	//
	r.HandleFunc("/adminPanel", adminPanel)
	r.HandleFunc("/accessСontrol", accessСontrol)
	r.HandleFunc("/saveUserPosition", saveUserPosition)
	//
	r.HandleFunc("/teacherPanel", teacherPanel)
	r.HandleFunc("/questionsСontrol", questionsСontrol)
	r.HandleFunc("/saveQuestion", saveQuestion)
	//
	r.HandleFunc("/studentPanel", studentPanel)
	r.HandleFunc("/testing", testing)
	r.HandleFunc("/checkAndSaveTest", checkAndSaveTest)
	//
	r.HandleFunc("/createCourse", createCourse)
	r.HandleFunc("/requestToСreateСourse", requestToСreateСourse)
	//
	r.HandleFunc("/adminRequests", adminRequests)
    r.HandleFunc("/approveRequest/{id:[0-9]+}", approveRequest)
	r.HandleFunc("/rejectRequest/{id:[0-9]+}", rejectRequest)
	r.HandleFunc("/course/{id:[0-9]+}", course)
	r.HandleFunc("/student/course/{id:[0-9]+}", studentCourse)
	r.HandleFunc("/student/test/{id:[0-9]+}", studentTest)
	r.HandleFunc("/teacher/test/{id:[0-9]+}/{slug}", teacherTestInfoStd)
	r.HandleFunc("/teacher/test/{id:[0-9]+}", teacherTestInfo)
	r.HandleFunc("/saveStudentTest", saveStudentTest)
	r.HandleFunc("/teacherCourses", teacherCourses)
	//
	r.HandleFunc("/courseOverview", courseOverview)
	r.HandleFunc("/findCourse", findCourse)
	r.HandleFunc("/checkCodeword", checkCodeword)
	r.HandleFunc("/deleteQuestion/{id:[0-9]+}", deleteQuestion)
	r.HandleFunc("/addCourseQuestion/{id:[0-9]+}", addCourseQuestion)
	r.HandleFunc("/saveCourseQuestion", saveCourseQuestion)
	//
	r.HandleFunc("/studentCourses", studentCourses)
	r.HandleFunc("/tizd", tizd)
	r.HandleFunc("/classIris", classIris)
	
	//
	r.HandleFunc("/course/{id:[0-9]+}/createTest", createTest)
	r.HandleFunc("/saveNewTest", saveNewTest)
	//
	r.HandleFunc("/saveUser", saveUser)
	r.HandleFunc("/deleteUser", deleteUser)
	http.ListenAndServe(":8080", r)
}


func main() {
	fmt.Println("Переход -> http://localhost:8080/")
	handleFunc()
}