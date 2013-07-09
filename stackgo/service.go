package stackgo

import (
	"log"
	"os/exec"
)

type Service struct {
	Name string `json:"name"`
}

func (s Service) Create() error {
	log.Println("Start service:", s.Name)
	out, err := exec.Command("service", s.Name, "start").Output()

	if err != nil {
		log.Println(string(out))
		return err
	}

	return nil
}
