package main

import (
	"fmt"
)

//Info .....
type Info struct {
	UserName, UserStatus, UserPosition string
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

func types() {
	fmt.Println("types")
}