package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	s "strings"
	"time"
)

func teacherPanel(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/teacherPanel.html", "templates/header.html", "templates/footer.html")

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
	t, err := template.ParseFiles("templates/questionsСontrol.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	var info Info
	info.UserName = getUserName(r)
	info.UserStatus = getUserStatus(r)
	info.UserPosition = getUserPosition(r)

	t.ExecuteTemplate(w, "questionsСontrol", info)
}

func saveQuestion(w http.ResponseWriter, r *http.Request) {
	question := r.FormValue("question")
	answer := r.FormValue("answer")

	if question == "" || answer == "" {
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
	t, err := template.ParseFiles("templates/createCourse.html", "templates/header.html", "templates/footer.html")

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

	res, err := db.Query(fmt.Sprintf("SELECT id FROM users where user_name='%s'", getUserName(r)))
	if err != nil {
		panic(err)
	}
	var user_id string
	for res.Next() {
		err = res.Scan(&user_id)
		if err != nil {
			panic(err)
		}
	}
	defer res.Close()

	insert, err := db.Query(fmt.Sprintf("INSERT INTO course_requests (teacher_name, course_name, codeword, teacher_id) VALUES('%s','%s', '%s', '%s')", getUserName(r), r.FormValue("courseName"), r.FormValue("codeword"), user_id))
	if err != nil {
		panic(err)
	} else {
		message := "Вы успешно отправили заявку!"
		t, err := template.ParseFiles("templates/message.html", "templates/footer.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		t.ExecuteTemplate(w, "message", message)
	}
	defer insert.Close()
}

func teacherCourses(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/teacherCourses.html", "templates/header.html", "templates/footer.html")
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

	res, err := db.Query(fmt.Sprintf("SELECT id, course_name FROM courses where teacher_name='%s'", info.UserName))
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

	t.ExecuteTemplate(w, "teacherCourses", struct{ Info, Course interface{} }{info, courses})
}

func course(w http.ResponseWriter, r *http.Request) {
	id := s.Replace(fmt.Sprint(r.URL), "/course/", "", -1)

	connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query(fmt.Sprintf("SELECT course_name FROM courses where id='%s'", id))
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

	res, err = db.Query(fmt.Sprintf("SELECT student_id FROM course_subscribers where course_id='%s'", id))
	if err != nil {
		panic(err)
	}

	student_ids := []string{}
	for res.Next() {
		var student_id string
		res.Scan(&student_id)
		student_ids = append(student_ids, student_id)
	}
	defer res.Close()

	students := []Student{}
	for _, student_id := range student_ids {
		res, err := db.Query(fmt.Sprintf("SELECT last_name, first_name FROM users where id='%s'", student_id))
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
	res, err = db.Query(fmt.Sprintf("SELECT id, question, answer FROM questions where course_id='%s'", id))
	if err != nil {
		panic(err)
	}
	defer res.Close()

	indx := 1
	tasks := []Task{}
	for res.Next() {
		var task Task
		err = res.Scan(&task.ID, &task.Question, &task.Answer)
		if err != nil {
			panic(err)
		}
		task.Index = uint16(indx)
		indx++
		tasks = append(tasks, task)
	}
	taskCount := len(tasks)

	//нужно передать массив тестов<<
	resTest, err := db.Query(fmt.Sprintf("SELECT DISTINCT test_name FROM tests where course_id='%s'", id))
	if err != nil {
		panic(err)
	}
	defer resTest.Close()

	type Test struct {
		ID   uint16
		Name string
	}
	tests := []Test{}
	for resTest.Next() {
		var test Test
		err = resTest.Scan(&test.Name)
		if err != nil {
			panic(err)
		}

		resid, err := db.Query(fmt.Sprintf("SELECT DISTINCT id FROM tests where test_name='%s'", test.Name))
		if err != nil {
			panic(err)
		}
		defer resid.Close()
		for resid.Next() {
			var id uint16
			err = resid.Scan(&id)
			if err != nil {
				panic(err)
			}
			test.ID = id
		}
		tests = append(tests, test)
	}
	//нужно передать массив тестов>>
	t, err := template.ParseFiles("templates/header.html", "templates/coursePage.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "coursePage", struct{ Info, Students, Count, Tasks, TaskCount, CourseName, CourseID, Tests interface{} }{info, students, count, tasks, taskCount, courseName, id, tests})
}

func deleteQuestion(w http.ResponseWriter, r *http.Request) {
	id := s.Replace(fmt.Sprint(r.URL), "/deleteQuestion/", "", -1)
	connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.Query(fmt.Sprintf("delete from questions where id='%s'", id))
	http.Redirect(w, r, r.Header.Get("Referer"), 302)
}

func addCourseQuestion(w http.ResponseWriter, r *http.Request) {
	id := s.Replace(fmt.Sprint(r.URL), "/addCourseQuestion/", "", -1) //id курса
	t, err := template.ParseFiles("templates/header.html", "templates/addCourseQuestion.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	var info Info
	info.UserName = getUserName(r)
	info.UserStatus = getUserStatus(r)
	info.UserPosition = getUserPosition(r)
	t.ExecuteTemplate(w, "addCourseQuestion", struct{ Info, CourseID interface{} }{info, id})
}

func saveCourseQuestion(w http.ResponseWriter, r *http.Request) {
	question := r.FormValue("question")
	answer := r.FormValue("answer")
	topic := r.FormValue("topic")
	courseID := r.FormValue("courseid")

	connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	insert, err := db.Query(fmt.Sprintf("INSERT INTO questions (question, answer, course_id, topic) VALUES('%s','%s','%s','%s')", question, answer, courseID, topic))
	if err != nil {
		panic(err)
	}
	defer insert.Close()
	http.Redirect(w, r, r.Header.Get("Referer"), 302)
}

func createTest(w http.ResponseWriter, r *http.Request) {
	courseID := s.Replace(fmt.Sprint(r.URL), "/createTest", "", -1)
	courseID = s.Replace(courseID, "/course/", "", -1)
	//--------------------------------------------------------
	connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//--------------------------------------------------------
	//берем из базы уникальные топики
	res, err := db.Query(fmt.Sprintf("SELECT distinct topic FROM questions where course_id='%s'", courseID))
	if err != nil {
		panic(err)
	}
	defer res.Close()

	topics := []string{}
	for res.Next() {
		var topic string
		err = res.Scan(&topic)
		if err != nil {
			panic(err)
		}
		// fmt.Println(topic)
		topics = append(topics, topic)
	}
	topicCount := len(topics)
	topicCounters := []int16{}
	Topics := []Topic{}
	for i := 0; i < topicCount; i++ {
		res, err := db.Query(fmt.Sprintf("select count(topic) from questions where topic='%s'", topics[i]))
		if err != nil {
			panic(err)
		}
		var topicCounter int16
		for res.Next() {
			err = res.Scan(&topicCounter)
			if err != nil {
				panic(err)
			}
			var topic Topic
			topic.MaxValue = topicCounter
			topic.Name = topics[i]
			Topics = append(Topics, topic)
			topicCounters = append(topicCounters, topicCounter)
		}

		fmt.Print(topics[i])
		fmt.Println(topicCounter)
		fmt.Println("=====")
	}
	fmt.Println(topicCounters)

	//--------------------------------------------------------

	//--------------------------------------------------------
	//--------------------------------------------------------

	//--------------------------------------------------------
	//--------------------------------------------------------
	t, err := template.ParseFiles("templates/header.html", "templates/createTest.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	var info Info
	info.UserName = getUserName(r)
	info.UserStatus = getUserStatus(r)
	info.UserPosition = getUserPosition(r)
	t.ExecuteTemplate(w, "createTest", struct{ Info, CourseID, Topics interface{} }{info, courseID, Topics})
}

func saveNewTest(w http.ResponseWriter, r *http.Request) {
	courseID := r.FormValue("myid") //достал id курса
	testName := r.FormValue("testName")
	testTime := r.FormValue("testTime")

	//--------------------------------------------------------
	connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//--------------------------------------------------------
	//берем из базы уникальные топики
	res, err := db.Query(fmt.Sprintf("SELECT distinct topic FROM questions where course_id='%s'", courseID))
	if err != nil {
		panic(err)
	}
	defer res.Close()

	topics := []string{}
	for res.Next() {
		var topic string
		err = res.Scan(&topic)
		if err != nil {
			panic(err)
		}
		// fmt.Println(topic)
		topics = append(topics, topic)
	}
	topicCount := len(topics)
	topicCounters := []int16{}
	Topics := []Topic{}
	for i := 0; i < topicCount; i++ {
		res, err := db.Query(fmt.Sprintf("select count(topic) from questions where topic='%s'", topics[i]))
		if err != nil {
			panic(err)
		}
		var topicCounter int16
		for res.Next() {
			err = res.Scan(&topicCounter)
			if err != nil {
				panic(err)
			}
			var topic Topic
			topic.MaxValue = topicCounter
			topic.Name = topics[i]
			Topics = append(Topics, topic)
			topicCounters = append(topicCounters, topicCounter)
		}
	}

	today := time.Now()

	topicsQuestionsCount := []string{}
	for i := 0; i < len(Topics); i++ {
		var topicQuestionsCount string
		topicQuestionsCount = r.FormValue("name_" + Topics[i].Name)
		topicsQuestionsCount = append(topicsQuestionsCount, topicQuestionsCount)
	}
	//--------------------------------------------------------
	for i := 0; i < len(Topics); i++ {

		insert, err := db.Query(fmt.Sprintf("INSERT INTO tests (test_name, course_id, topic, questions_count,time) VALUES('%s','%s','%s','%s','%s')", testName+" ("+today.Format("2006-01-02 15:04:05")+")", courseID, Topics[i].Name, topicsQuestionsCount[i], testTime))
		if err != nil {
			panic(err)
		}
		defer insert.Close()
	}

	http.Redirect(w, r, r.Header.Get("Referer"), 302)
}

func teacherTestInfo(w http.ResponseWriter, r *http.Request) {
	testID := s.Replace(fmt.Sprint(r.URL), "/teacher/test/", "", -1)
	connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	res, err := db.Query(fmt.Sprintf("SELECT distinct student_name FROM student_answers where test_id='%s'", testID))
	if err != nil {
		panic(err)
	}

	stdnts := []string{}
	for res.Next() {
		var stdnt string
		err = res.Scan(&stdnt)
		if err != nil {
			panic(err)
		}
		stdnts = append(stdnts, stdnt)
	}
	defer res.Close()
	t, err := template.ParseFiles("templates/teacherTestInfo.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	var info Info
	info.UserName = getUserName(r)
	info.UserStatus = getUserStatus(r)
	info.UserPosition = getUserPosition(r)
	t.ExecuteTemplate(w, "teacherTestInfo", struct{ Info, Stdnts interface{} }{info, stdnts})
}

func teacherTestInfoStd(w http.ResponseWriter, r *http.Request) {
	testIdstudentName := s.Replace(fmt.Sprint(r.URL), "/teacher/test/", "", -1)
	testID := s.Replace(testIdstudentName, "/", " ", -1)
	a := s.Fields(testID)
	fmt.Println(a[0])

	connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var pp int32
	res3, err := db.Query(fmt.Sprintf("SELECT count(date) from student_answers where test_id='%s' and student_name='%s'", a[0], a[1]))
	if err != nil {
		panic(err)
	}
	for res3.Next() {
		err = res3.Scan(&pp)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("всего дат", pp)
	var ppp int32
	res2, err := db.Query(fmt.Sprintf("SELECT count(distinct date) from student_answers where test_id='%s' and student_name='%s'", a[0], a[1]))
	if err != nil {
		panic(err)
	}
	for res2.Next() {
		err = res2.Scan(&ppp)
		if err != nil {
			panic(err)
		}
	}
	//-------
	unicDates := []string{}
	res22, err := db.Query(fmt.Sprintf("SELECT distinct date from student_answers where test_id='%s' and student_name='%s'", a[0], a[1]))
	if err != nil {
		panic(err)
	}
	for res22.Next() {
		var date string
		err = res22.Scan(&date)
		if err != nil {
			panic(err)
		}
		unicDates = append(unicDates, date)
	}
	fmt.Println(unicDates)
	//-------
	fmt.Println("уникальных дат", ppp)
	var quesInTestCount float32
	quesInTestCount = float32(pp) / float32(ppp)
	fmt.Println("вопросов в тесте", quesInTestCount)
	defer res2.Close()
	type StdntAnswr struct {
		Question, StudentAnswer, Answer, StudentName, Date string
		Mark                                               uint16
	}

	rrr := [][]StdntAnswr{}
	for i := 0; i < len(unicDates); i++ {
		stdntAnswrs := []StdntAnswr{}
		res, err := db.Query(fmt.Sprintf("SELECT question, student_answer, answer, mark, date, student_name FROM student_answers where test_id='%s' and student_name='%s' and date='%s'", a[0], a[1], unicDates[i]))
		if err != nil {
			panic(err)
		}
		for res.Next() {
			var stdntAnswr StdntAnswr
			err = res.Scan(&stdntAnswr.Question, &stdntAnswr.StudentAnswer, &stdntAnswr.Answer, &stdntAnswr.Mark, &stdntAnswr.Date, &stdntAnswr.StudentName)
			if err != nil {
				panic(err)
			}
			stdntAnswrs = append(stdntAnswrs, stdntAnswr)
		}
		defer res.Close()
		rrr = append(rrr, stdntAnswrs)
	}

	fmt.Println(rrr)
	t, err := template.ParseFiles("templates/teacherTestInfoCopy.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	var info Info
	info.UserName = getUserName(r)
	info.UserStatus = getUserStatus(r)
	info.UserPosition = getUserPosition(r)
	t.ExecuteTemplate(w, "teacherTestInfoCopy", struct{ Info, StdntAnswrs, StdName interface{} }{info, rrr, a[1]})
}

// func teacherTestInfo(w http.ResponseWriter, r *http.Request) {
// 	testID := s.Replace(fmt.Sprint(r.URL), "/teacher/test/", "", -1)
// 	connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
// 	db, err := sql.Open("postgres", connStr)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()
// 	res, err := db.Query(fmt.Sprintf("SELECT question, student_answer, answer, mark, date, student_name FROM student_answers where test_id='%s' and student_name='%s'", testID, student_name))
// 	if err != nil {
// 		panic(err)
// 	}
// 	type StdntAnswr struct {
// 		Question, StudentAnswer, Answer, StudentName, Date string
// 		Mark                                               uint16
// 	}
// 	stdntAnswrs := []StdntAnswr{}
// 	for res.Next() {
// 		var stdntAnswr StdntAnswr
// 		err = res.Scan(&stdntAnswr.Question, &stdntAnswr.StudentAnswer, &stdntAnswr.Answer, &stdntAnswr.Mark, &stdntAnswr.Date, &stdntAnswr.StudentName)
// 		if err != nil {
// 			panic(err)
// 		}
// 		stdntAnswrs = append(stdntAnswrs, stdntAnswr)
// 	}
// 	defer res.Close()
// 	fmt.Println(stdntAnswrs)
// 	t, err := template.ParseFiles("templates/teacherTestInfo.html", "templates/header.html", "templates/footer.html"); if err != nil {
// 		fmt.Fprintf(w, err.Error())
// 	}
// 	var info Info
// 	info.UserName = getUserName(r)
// 	info.UserStatus = getUserStatus(r)
// 	info.UserPosition = getUserPosition(r)
// 	t.ExecuteTemplate(w, "teacherTestInfo", struct{ Info, StdntAnswrs interface{} }{info, stdntAnswrs})
// }
