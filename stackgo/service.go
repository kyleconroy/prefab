package stackgo

import (
	"log"
	"os/exec"
	"strings"
)

type Service struct {
	Name string `json:"name"`
}

func (s Service) Create() error {
	log.Println("Start service:", s.Name)
	out, err := exec.Command("service", s.Name, "status").CombinedOutput()

	if err != nil {
		log.Println(string(out))
		return err
	}

	if strings.Contains(string(out), "start/running") {
		return nil
	}

	out, err = exec.Command("service", s.Name, "start").CombinedOutput()

	if err != nil {
		log.Println(string(out))
		return err
	}

	return nil
}
