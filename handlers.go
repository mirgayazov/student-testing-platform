package main

import (
	"fmt"
	"net/http"
	"html/template"
	"encoding/json"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
)

func tizd(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/tizd.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "tizd", nil)	
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html","templates/header.html","templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	var info Info
	info.UserName = getUserName(r)
	info.UserStatus = getUserStatus(r)
	info.UserPosition = getUserPosition(r)

	t.ExecuteTemplate(w, "index", info)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html","templates/header.html","templates/footer.html")	

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)
}

//Registration .....
func registration(w  http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/registration.html","templates/header.html","templates/footer.html")	

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	var info Info
	info.UserName = getUserName(r)
	info.UserStatus = getUserStatus(r)
	info.UserPosition = getUserPosition(r)

	message := ""
	t.ExecuteTemplate(w, "registration", struct{Info, Message interface{}}{info, message});
}

func saveUser(w  http.ResponseWriter, r *http.Request) {
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	userName := r.FormValue("user_name")
	password := []byte(r.FormValue("password"))

	if firstName == "" || lastName == "" || userName == "" || len(password) ==0 {
	} else {
		connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		res, err := db.Query(fmt.Sprintf("select count(user_name) from users where user_name='%s'",userName))
		if err != nil {
			panic(err)
		}
		var count int
		for res.Next() {
			err = res.Scan(&count)
			if err != nil {
				panic(err)
			}
		}
		defer res.Close()
		if count == 0 {
			hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
			if err != nil {
				panic(err)
			}
			
			insert, err := db.Query(fmt.Sprintf("INSERT INTO users (first_name, last_name, hash, is_active,user_name, position) VALUES('%s','%s','%s', 'false','%s','student')", firstName, lastName, hashedPassword,userName))
			if err != nil {
				panic(err)
			}
			defer insert.Close()

			http.Redirect(w, r, "/authorization", http.StatusSeeOther)
		} else {
			t, err := template.ParseFiles("templates/registration.html","templates/header.html","templates/footer.html")	
			if err != nil {
				fmt.Fprintf(w, err.Error())
			}
			var info Info
			info.UserName = getUserName(r)
			info.UserStatus = getUserStatus(r)
			info.UserPosition = getUserPosition(r)
			message := "Введенное имя пользователя ("+userName+") занято"
			t.ExecuteTemplate(w, "registration", struct{Info, Message interface{}}{info, message});
		}
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	logout(w,r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	db.Query(fmt.Sprintf("delete from users where user_name='%s'", getUserName(r)))
}

func login(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("user_name")
	password := []byte(r.FormValue("password"))

	if userName == "" || r.FormValue("password") == ""{
		message := "Заполните все поля!"
		t, err := template.ParseFiles("templates/message.html","templates/footer.html")	
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		t.ExecuteTemplate(w, "message", message)
	} else {
		connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
		defer db.Close()
		
		res, err := db.Query(fmt.Sprintf("SELECT id, hash, position FROM public.users where user_name='%s'", userName))
		if err != nil {
			panic(err)
		}

		var user User
		for res.Next() {
			err = res.Scan(&user.ID, &user.Hash, &user.Position)
			if err != nil {
				panic(err)
			}
		}

		err = bcrypt.CompareHashAndPassword(user.Hash, []byte(password))
		if err != nil {
			message := "Вы ввели неверные данные!"
			t, err := template.ParseFiles("templates/message.html","templates/footer.html")	
			if err != nil {
				fmt.Fprintf(w, err.Error())
			}
			t.ExecuteTemplate(w, "message", message)
		} else {
			setSession(userName, user.Position, w)
			_, err := db.Query(fmt.Sprintf("UPDATE users SET is_active = 'true' WHERE id = '%d'", user.ID))
			if err != nil {
				panic(err)
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}

func authorization(w  http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/authorization.html","templates/header.html","templates/footer.html")	

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	var info Info
	info.UserName = getUserName(r)
	info.UserStatus = getUserStatus(r)
	info.UserPosition = getUserPosition(r)

	t.ExecuteTemplate(w, "authorization", info)	
}

func logout(w http.ResponseWriter, r *http.Request) {
	userName := getUserName(r)
	connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	
	_, err = db.Query(fmt.Sprintf("UPDATE users SET is_active = 'false' WHERE user_name = '%s'", userName))
	if err != nil {
		panic(err)
	}
	clearSession(w)
	http.Redirect(w, r, "/", 302)
}

func about(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/about.html","templates/header.html","templates/footer.html")	

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	var info Info
	info.UserName = getUserName(r)
	info.UserStatus = getUserStatus(r)
	info.UserPosition = getUserPosition(r)

	jdata,err := json.Marshal(info)
	if err != nil{
		return 
	}
	fmt.Println(jdata)

	t.ExecuteTemplate(w, "about", struct{Info, Message interface{}}{info, string(jdata)});
	// t.ExecuteTemplate(w, "about", info)
}