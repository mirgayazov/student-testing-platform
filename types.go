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
	TeacherName, CourseName, Codeword string
}

//Course .....
type Course struct {
	ID uint16
	CourseName string
}

//Correctcodeword .....
type Correctcodeword struct {
	value string
}

//ID .....
type ID struct {
	value string
}

//Subscribers .....
type Subscribers struct {
	Values []uint8
}

//Student .....
type Student struct {
	LastName, FirstName string
}
