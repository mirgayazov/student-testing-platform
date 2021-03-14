package main

//Info .....
type Info struct {
	UserName string `json:"username"`
	UserStatus string `json:"userstatus"`
	UserPosition string `json:"userposition"`
}

//Task .....
type Task struct {
	ID, Index uint16
	Question, Answer string
}

//Topic .....
type Topic struct {
	MaxValue int16
	Name string
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
	TeacherName, CourseName, Codeword, TeacherID string
}

//Course .....
type Course struct {
	ID int64
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
