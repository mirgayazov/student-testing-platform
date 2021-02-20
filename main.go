package main

import (
	"fmt"
	"net/http"
	"html/template"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/bcrypt"
)

// cookie handling

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
)


func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

//Article ....
type Article struct {
	ID uint16
	Title, Anons, FullText string
}	

//User .....
type User struct {
	ID uint16
	hash []byte
}

var posts =[]Article{}

func index(w http.ResponseWriter, r *http.Request) {
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

	posts = []Article{}

	for res.Next() {
		var post Article
		err = res.Scan(&post.ID, &post.Title,  &post.Anons,  &post.FullText)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}
	// defer res.Close()
	t.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html","templates/header.html","templates/footer.html")	

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)
}

func registration(w  http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/registration.html","templates/header.html","templates/footer.html")	

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "registration", nil)	
}

func saveUser(w  http.ResponseWriter, r *http.Request) {
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

func saveArticle(w http.ResponseWriter, r *http.Request) {
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

//LogIn .....
func login(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("user_name")
	password := []byte(r.FormValue("password"))


	if userName == "" || r.FormValue("password") == ""{
		fmt.Fprintf(w, "Введены неверные данные")
	} else {

		setSession(userName, w)

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

		var user User
		for res.Next() {
			err = res.Scan(&user.ID, &user.hash)
			if err != nil {
				panic(err)
			}
		}

		err = bcrypt.CompareHashAndPassword(user.hash, []byte(password))
		if err != nil {
			fmt.Fprintf(w, "Введены неверные данные")
		} else {
			_, err := db.Query(fmt.Sprintf("UPDATE users SET is_active = 'true' WHERE id = '%d'", user.ID))
			if err != nil {
				panic(err)
			}
		}

		//
		t, err := template.ParseFiles("templates/index.html","templates/header.html","templates/footer.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		status := "true"
		t.ExecuteTemplate(w, "index", status)
		//
		// http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func authorization(w  http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/authorization.html","templates/header.html","templates/footer.html")	

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "authorization", nil)	
}


func handleFunc() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/create", create)
	http.HandleFunc("/save_article", saveArticle)
	http.HandleFunc("/registration", registration)
	http.HandleFunc("/authorization", authorization)
	http.HandleFunc("/login", login)
	http.HandleFunc("/saveUser", saveUser)
	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("Переход -> http://localhost:8080/")
	handleFunc()
}