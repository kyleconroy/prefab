package stackgo

import (
	"log"
)

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (u User) Create() error {
	log.Println("Create user:", u.Name)
	return nil
}
