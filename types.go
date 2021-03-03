package main

//Info .....
type Info struct {
	UserName, UserStatus, UserPosition string
}

//Task .....
type Task struct {
	ID, Index uint16
	Question, Answer string
}

//Article ....
type Article struct {
	ID uint16
	Title, Anons, FullText string
}	

//Status ....
type Status struct {
	Status string
}

//User .....
type User struct {
	ID uint16
	Hash []byte
	Position string
}

//Request .....
type Request struct {
	ID uint16
	TeacherName, CourseName string
}