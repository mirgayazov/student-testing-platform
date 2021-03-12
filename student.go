package main

import (
	"fmt"
	"net/http"
	"html/template"
	"database/sql"
	"math/rand"
	"time"
	"strconv"
)

func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

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
var tasks2 =[]Task{}
var ids =[]uint16{}
var ras =[]string{}

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
		i=i+1
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
	tasks2[0].Index=1
	tasks2[1].Index=2
	tasks2[2].Index=3
	t.ExecuteTemplate(w, "testing", tasks2)
}

func checkAndSaveTest(w  http.ResponseWriter, r *http.Request) {
	uas := []string{r.FormValue("ua1"),r.FormValue("ua2"),r.FormValue("ua3")}
	dbids := []string{r.FormValue("id1"),r.FormValue("id2"),r.FormValue("id3")}

	connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	ras = []string{}
	count :=0
	for i := 0; i < len(uas); i++{
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
			if ra==uas[i] {
				fmt.Println("Верный ответ, вопрос: ",i+1)
				count++
			}
		}
    } 
	message := fmt.Sprintf("Вы ответили правильно на %d вопроса(ов)", count)
	t, err := template.ParseFiles("templates/message.html","templates/footer.html")	
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "message", message)
}

func courseOverview(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/student/courseOverview.html","templates/header.html","templates/footer.html")	
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

	t.ExecuteTemplate(w, "courseOverview", struct{Info, Course interface{}}{info, courses});
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
	percent :="%"
	res, err := db.Query(fmt.Sprintf("SELECT id, course_name FROM courses where course_name similar to '%s%s' or teacher_name='%s'",courseOrTeacherName, percent, courseOrTeacherName))
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

	t, err := template.ParseFiles("templates/student/courseOverview.html","templates/header.html","templates/footer.html")	
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "courseOverview", struct{Info, Course interface{}}{info, courses});
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
		t, err := template.ParseFiles("templates/course_subscribe_message.html","templates/footer.html")	
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		status := "correctcodeword"
		t.ExecuteTemplate(w, "course_subscribe_message", struct{Message, Status interface{}}{message, status});
		//добавление в подписичики курса -->
		res, err := db.Query(fmt.Sprintf("select id from users where user_name='%s'",getUserName(r)))
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
		db.Query(fmt.Sprintf("UPDATE courses SET subscribers = array_append((select subscribers from courses where id='%s') , '%s') WHERE id='%s';", courseID, userID.value, courseID))

		insert, err = db.Query(fmt.Sprintf("INSERT INTO student_courses (student_name, course_id) VALUES('%s','%s')", getUserName(r), courseID))
		if err != nil {
			panic(err)
		}
		defer insert.Close()
		//добавление в подписичики курса <--
	} else {
		message := "Неверное кодовое слово :("
		t, err := template.ParseFiles("templates/course_subscribe_message.html","templates/footer.html")	
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		status := "invalidcodeword"
		t.ExecuteTemplate(w, "course_subscribe_message", struct{Message, Status interface{}}{message, status});
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

	res, err := db.Query(fmt.Sprintf("select id from users where user_name='%s'",getUserName(r)))
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

	res, err = db.Query(fmt.Sprintf("select course_id from course_subscribers where student_id='%s'",userID))
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
		course.CourseName =courseName
		course.ID = parsedID
		
		corses = append(corses, course)
	}

	var info Info
	info.UserName = getUserName(r)
	info.UserStatus = getUserStatus(r)
	info.UserPosition = getUserPosition(r)

	t, err := template.ParseFiles("templates/header.html","templates/studentCourses.html","templates/footer.html")	
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "studentCourses", struct{Info, Course interface{}}{info, corses});
}