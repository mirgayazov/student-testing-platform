package types

import (
	"fmt"
)

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
}

func types() {
	fmt.Println("types")
}