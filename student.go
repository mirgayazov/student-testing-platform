package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	s "strings"
	"time"
)

func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func studentPanel(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/studentPanel.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	var info Info
	info.UserName = getUserName(r)
	info.UserStatus = getUserStatus(r)
	info.UserPosition = getUserPosition(r)

	t.ExecuteTemplate(w, "studentPanel", info)
}

var tasks = []Task{}
var tasks2 = []Task{}
var ids = []uint16{}
var ras = []string{}

func testing(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/testing.html", "templates/header.html", "templates/footer.html")
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
	var i uint16
	i = 1
	tasks = []Task{}
	ids = []uint16{}
	for res.Next() {
		var task Task
		err = res.Scan(&task.ID, &task.Question, &task.Answer)
		if err != nil {
			panic(err)
		}
		task.Index = i
		tasks = append(tasks, task)
		ids = append(ids, task.ID)
		i = i + 1
		//создал массив id-шников осталось рандомом достать 3 штуки и по id достать вопрос + ответ
	}
	//----------------------------------------------------------------------------
	// fmt.Println(fmt.Sprintf("Длина массива idшников %d",len(ids)))
	// fmt.Println(fmt.Sprintf("[%d,%d]",int(ids[0]),int(ids[len(ids)-1])))
	// rndindxs := []int{}
	// for i := 0; i < 3; i++ {
	// 	for len(rndindxs) < i+1 {
	// 		rndindx := random(int(ids[0]),int(ids[len(ids)-1]))
	// 		for j := 0; j < len(rndindxs)+1; j++ {
	// 			if rndindx != rndindxs[j] {
	// 				rndindxs = append(rndindxs, rndindx)
	// 				fmt.Println(rndindx)
	// 			}
	// 		}
	// 	}
	// }
	//----------------------------------------------------------------------------
	tasks2 = []Task{}
	tasks2 = append(tasks2, tasks[0])
	tasks2 = append(tasks2, tasks[4])
	tasks2 = append(tasks2, tasks[5])
	tasks2[0].Index = 1
	tasks2[1].Index = 2
	tasks2[2].Index = 3
	t.ExecuteTemplate(w, "testing", tasks2)
}

func checkAndSaveTest(w http.ResponseWriter, r *http.Request) {
	uas := []string{r.FormValue("ua1"), r.FormValue("ua2"), r.FormValue("ua3")}
	dbids := []string{r.FormValue("id1"), r.FormValue("id2"), r.FormValue("id3")}

	connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	ras = []string{}
	count := 0
	for i := 0; i < len(uas); i++ {
		res, err := db.Query(fmt.Sprintf("SELECT answer FROM public.questions where id=%s", dbids[i]))
		if err != nil {
			panic(err)
		}
		for res.Next() {
			var ra string
			err = res.Scan(&ra)
			if err != nil {
				panic(err)
			}
			ras = append(ras, ra)
			if ra == uas[i] {
				fmt.Println("Верный ответ, вопрос: ", i+1)
				count++
			}
		}
	}
	message := fmt.Sprintf("Вы ответили правильно на %d вопроса(ов)", count)
	t, err := template.ParseFiles("templates/message.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "message", message)
}

func courseOverview(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/student/courseOverview.html", "templates/header.html", "templates/footer.html")
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

	res, err := db.Query("SELECT id, course_name FROM courses")
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

	t.ExecuteTemplate(w, "courseOverview", struct{ Info, Course interface{} }{info, courses})
}

func findCourse(w http.ResponseWriter, r *http.Request) {
	courseOrTeacherName := r.FormValue("courseOrTeacherName")

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
	percent := "%"
	res, err := db.Query(fmt.Sprintf("SELECT id, course_name FROM courses where course_name similar to '%s%s' or teacher_name='%s'", courseOrTeacherName, percent, courseOrTeacherName))
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

	t, err := template.ParseFiles("templates/student/courseOverview.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "courseOverview", struct{ Info, Course interface{} }{info, courses})
}

func checkCodeword(w http.ResponseWriter, r *http.Request) {
	codeword := r.FormValue("codeword")
	courseID := r.FormValue("courseID")

	// --> берем из бд кодовое слово по id курса - cравниваем
	connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("student.go 217")
		panic(err)
	}
	defer db.Close()

	res, err := db.Query(fmt.Sprintf("SELECT codeword FROM courses where id='%s'", courseID))
	if err != nil {
		fmt.Println("student.go 224")
		panic(err)
	}

	var correctcodeword Correctcodeword
	for res.Next() {
		err = res.Scan(&correctcodeword.value)
		if err != nil {
			panic(err)
		}
	}
	defer res.Close()

	if correctcodeword.value == codeword {
		message := "Вы успешно подписались на курс!"
		t, err := template.ParseFiles("templates/course_subscribe_message.html", "templates/footer.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		status := "correctcodeword"
		t.ExecuteTemplate(w, "course_subscribe_message", struct{ Message, Status interface{} }{message, status})
		//добавление в подписичики курса -->
		res, err := db.Query(fmt.Sprintf("select id from users where user_name='%s'", getUserName(r)))
		if err != nil {
			panic(err)
		}
		var userID ID
		for res.Next() {
			err = res.Scan(&userID.value)
			if err != nil {
				panic(err)
			}
		}
		defer res.Close()
		insert, err := db.Query(fmt.Sprintf("INSERT INTO course_subscribers (student_id, course_id) VALUES('%s','%s')", userID.value, courseID))
		if err != nil {
			panic(err)
		}
		defer insert.Close()
		//добавление в подписичики курса <--
	} else {
		message := "Неверное кодовое слово :("
		t, err := template.ParseFiles("templates/course_subscribe_message.html", "templates/footer.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		status := "invalidcodeword"
		t.ExecuteTemplate(w, "course_subscribe_message", struct{ Message, Status interface{} }{message, status})
	}
	//берем из бд кодовое слово по id курса - cравниваем <--
}

func studentCourses(w http.ResponseWriter, r *http.Request) {
	connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	courseIDS := []string{}
	corses := []Course{}

	res, err := db.Query(fmt.Sprintf("select id from users where user_name='%s'", getUserName(r)))
	if err != nil {
		panic(err)
	}
	var userID string
	for res.Next() {
		err = res.Scan(&userID)
		if err != nil {
			panic(err)
		}
	}
	defer res.Close()

	res, err = db.Query(fmt.Sprintf("select course_id from course_subscribers where student_id='%s'", userID))
	if err != nil {
		panic(err)
	}
	var courseID string
	var course Course
	for res.Next() {
		err = res.Scan(&courseID)
		if err != nil {
			panic(err)
		}
		courseIDS = append(courseIDS, courseID)
	}
	defer res.Close()

	var courseName string
	for i := 0; i < len(courseIDS); i++ {
		res, err := db.Query(fmt.Sprintf("select course_name from courses where id='%s'", courseIDS[i]))
		if err != nil {
			panic(err)
		}
		for res.Next() {
			err = res.Scan(&courseName)
			if err != nil {
				panic(err)
			}
		}
		defer res.Close()

		parsedID, err := strconv.ParseInt(courseIDS[i], 0, 16)
		if err != nil {
			panic(err)
		}
		course.CourseName = courseName
		course.ID = parsedID

		corses = append(corses, course)
	}

	var info Info
	info.UserName = getUserName(r)
	info.UserStatus = getUserStatus(r)
	info.UserPosition = getUserPosition(r)

	t, err := template.ParseFiles("templates/header.html", "templates/studentCourses.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "studentCourses", struct{ Info, Course interface{} }{info, corses})
}

func studentCourse(w http.ResponseWriter, r *http.Request) {
	type Test struct {
		ID      uint16
		Name    string
		Results []struct {
			Mark uint16
			Date string
		}
	}
	courseID := s.Replace(fmt.Sprint(r.URL), "/student/course/", "", -1)

	connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	tests := []Test{}
	res, err := db.Query(fmt.Sprintf("select distinct test_name from tests where course_id='%s'", courseID))
	if err != nil {
		panic(err)
	}
	for res.Next() {
		var test Test
		err = res.Scan(&test.Name)
		if err != nil {
			panic(err)
		}

		res2, err := db.Query(fmt.Sprintf("select id from tests where test_name='%s'", test.Name))
		if err != nil {
			panic(err)
		}
		for res2.Next() {
			err = res2.Scan(&test.ID)
			if err != nil {
				panic(err)
			}
		}
		defer res2.Close()

		tstres, err := db.Query(fmt.Sprintf("select distinct mark, date from student_answers where test_id =%d and student_name='%s'", test.ID, getUserName(r)))
		if err != nil {
			panic(err)
		}

		for tstres.Next() {
			var mark uint16
			var date string
			err = tstres.Scan(&mark, &date)
			if err != nil {
				panic(err)
			}
			test.Results = append(test.Results, struct {
				Mark uint16
				Date string
			}{mark, date})
		}
		defer tstres.Close()

		tests = append(tests, test)
	}
	defer res.Close()

	var info Info
	info.UserName = getUserName(r)
	info.UserStatus = getUserStatus(r)
	info.UserPosition = getUserPosition(r)

	t, err := template.ParseFiles("templates/header.html", "templates/studentCourse.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "studentCourse", struct{ Info, Test interface{} }{info, tests})
}

func studentTest(w http.ResponseWriter, r *http.Request) {
	type Block struct {
		Topic, QuestionsCount string
		Questions             []struct {
			ID    uint16
			Value string
		}
	}
	testID := s.Replace(fmt.Sprint(r.URL), "/student/test/", "", -1)

	connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var testName string
	var testTime int
	res, err := db.Query(fmt.Sprintf("select test_name, time from tests where id='%s'", testID))
	if err != nil {
		panic(err)
	}
	for res.Next() {
		err = res.Scan(&testName, &testTime)
		if err != nil {
			panic(err)
		}
	}
	defer res.Close()

	blocks := []Block{}
	ids := []uint16{}
	res, err = db.Query(fmt.Sprintf("select topic, questions_count from tests where test_name='%s'", testName))
	if err != nil {
		panic(err)
	}
	for res.Next() {
		var block Block
		err = res.Scan(&block.Topic, &block.QuestionsCount)
		if err != nil {
			panic(err)
		}
		if block.QuestionsCount != "0" {
			//----------------------------------------
			ques, err := db.Query(fmt.Sprintf("select id, question from questions where topic='%s'", block.Topic))
			if err != nil {
				panic(err)
			}
			for ques.Next() {
				var id uint16
				var question string
				err = ques.Scan(&id, &question)
				if err != nil {
					panic(err)
				}
				ids = append(ids, id)
				block.Questions = append(block.Questions, struct {
					ID    uint16
					Value string
				}{id, question})
			}
			defer ques.Close()

			blocks = append(blocks, block)
		}

	}
	defer res.Close()

	var info Info
	info.UserName = getUserName(r)
	info.UserStatus = getUserStatus(r)
	info.UserPosition = getUserPosition(r)

	t, err := template.ParseFiles("templates/header.html", "templates/studentTest.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "studentTest", struct{ Info, Block, Ids, TestID, TestTime interface{} }{info, blocks, ids, testID, testTime})
}

func saveStudentTest(w http.ResponseWriter, r *http.Request) {
	Ids := r.FormValue("ids")
	Ids = s.Replace(Ids, "[", "", -1)
	Ids = s.Replace(Ids, "]", "", -1)
	IdsArr := strings.Split(string(Ids), " ")
	TestID := r.FormValue("testID")

	connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	type StudentAnswer struct {
		Mark                                            uint16
		QuestionID, Question, StudentAnswer, TrueAnswer string
	}

	studentAnswers := []StudentAnswer{}
	var mark uint16
	for _, id := range IdsArr {
		var studentAnswer StudentAnswer
		studentAnswer.StudentAnswer = r.FormValue("ans/" + id)
		studentAnswer.QuestionID = id

		var trueAns string
		var ques string
		res, err := db.Query(fmt.Sprintf("select answer, question from questions where id='%s'", studentAnswer.QuestionID))
		if err != nil {
			panic(err)
		}
		for res.Next() {
			err = res.Scan(&trueAns, &ques)
			if err != nil {
				panic(err)
			}
			studentAnswer.TrueAnswer = trueAns
			if studentAnswer.StudentAnswer == studentAnswer.TrueAnswer {
				mark++
			}
		}
		defer res.Close()
		studentAnswer.Mark = mark
		studentAnswer.Question = ques
		studentAnswers = append(studentAnswers, studentAnswer)
	}

	today := time.Now()

	for _, stdntans := range studentAnswers {
		insert, err := db.Query(fmt.Sprintf("INSERT INTO student_answers (question, student_answer, answer, mark,date, student_name, test_id) VALUES('%s','%s','%s','%d','%s','%s', '%s')", stdntans.Question, stdntans.StudentAnswer, stdntans.TrueAnswer, mark, today.Format("2006-01-02 15:04:05"), getUserName(r), TestID))
		if err != nil {
			panic(err)
		}
		defer insert.Close()
	}

	res3, err := db.Query(fmt.Sprintf("select course_id from tests where id='%s'", TestID))
	if err != nil {
		panic(err)
	}
	var courseId string
	for res3.Next() {
		err = res3.Scan(&courseId)
		if err != nil {
			panic(err)
		}
	}
	defer res3.Close()

	http.Redirect(w, r, "/student/course/"+courseId, 302)
}
