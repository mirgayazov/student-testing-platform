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