package prefab

import (
	"log"
	"os/exec"
)

type Service struct {
	Name string `json:"name"`
}

func (s Service) Create() error {
	log.Println("Start service:", s.Name)
	out, err := exec.Command("service", s.Name, "status").CombinedOutput()

	if err == nil {
		log.Println("service", s.Name, "running.")
		return nil
	}

	out, err = exec.Command("service", s.Name, "start").CombinedOutput()

	if err != nil {
		log.Println(string(out))
		return err
	}

	return nil
}
