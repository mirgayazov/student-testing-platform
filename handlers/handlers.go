package handlers

import (
	"fmt"
	"./types"
	"net/http"
	"html/template"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"./cookies"
)

var posts =[]types.Article{}

//Index .....
func Index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html","templates/header.html","templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM public.articles")
	if err != nil {
		panic(err)
	}

	posts = []types.Article{}

	for res.Next() {
		var post types.Article
		err = res.Scan(&post.ID, &post.Title,  &post.Anons,  &post.FullText)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}
	userStatus := cookies.GetUserStatus(r)
	// if userStatus =="" {
	// 	fmt.Println("false")
	// }
	// defer res.Close()
	t.ExecuteTemplate(w, "index", userStatus)
}

//Create .....
func Create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html","templates/header.html","templates/footer.html")	

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)
}

//Registration .....
func Registration(w  http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/registration.html","templates/header.html","templates/footer.html")	

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	userStatus := cookies.GetUserStatus(r)

	t.ExecuteTemplate(w, "registration", userStatus)	
}

//SaveUser .....
func SaveUser(w  http.ResponseWriter, r *http.Request) {
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	userName := r.FormValue("user_name")
	password := []byte(r.FormValue("password"))

	if firstName == "" || lastName == ""{
		fmt.Fprintf(w, "Не все данные заполнены")
	} else {
		connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
		defer db.Close()
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		
		insert, err := db.Query(fmt.Sprintf("INSERT INTO users (first_name, last_name, hash, is_active,user_name) VALUES('%s','%s','%s', 'false','%s')", firstName, lastName, hashedPassword,userName))
		if err != nil {
			panic(err)
		}
		defer insert.Close()

		http.Redirect(w, r, "/authorization", http.StatusSeeOther)
	}
}

//SaveArticle .....
func SaveArticle(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	fullText := r.FormValue("full_text")

	if title == "" || anons == "" || fullText == "" {
		fmt.Fprintf(w, "Не все данные заполнены")
	} else {
		connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		insert, err := db.Query(fmt.Sprintf("INSERT INTO articles (title, anons, full_text) VALUES('%s','%s','%s')", title, anons, fullText))
		if err != nil {
			panic(err)
		}
		defer insert.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//Login .....
func Login(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("user_name")
	password := []byte(r.FormValue("password"))


	if userName == "" || r.FormValue("password") == ""{
		fmt.Fprintf(w, "Введены неверные данные")
	} else {

		cookies.SetSession(userName, w)
		connStr := "user=kamil password=1809 dbname=golang sslmode=disable"
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
		defer db.Close()
		
		res, err := db.Query(fmt.Sprintf("SELECT id, hash FROM public.users where user_name='%s'", userName))
		if err != nil {
			panic(err)
		}

		var user types.User
		for res.Next() {
			err = res.Scan(&user.ID, &user.Hash)
			if err != nil {
				panic(err)
			}
		}

		err = bcrypt.CompareHashAndPassword(user.Hash, []byte(password))
		if err != nil {
			fmt.Fprintf(w, "Введены неверные данные")
		} else {
			_, err := db.Query(fmt.Sprintf("UPDATE users SET is_active = 'true' WHERE id = '%d'", user.ID))
			if err != nil {
				panic(err)
			}
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//Authorization .....
func Authorization(w  http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/authorization.html","templates/header.html","templates/footer.html")	

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	userStatus := cookies.GetUserStatus(r)

	t.ExecuteTemplate(w, "authorization", userStatus)	
}

//Logout .....
func Logout(w http.ResponseWriter, r *http.Request) {
	userName := cookies.GetUserName(r)
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
	cookies.ClearSession(w)
	http.Redirect(w, r, "/", 302)
}

//About .....
func About(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/about.html","templates/header.html","templates/footer.html")	

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	userStatus := cookies.GetUserStatus(r)
	t.ExecuteTemplate(w, "about", userStatus)
}

func handlers() {
	fmt.Println("handlers")
} 