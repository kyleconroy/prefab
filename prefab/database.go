package prefab

import (
	"fmt"
	"log"
	"os/exec"
)

type DatabaseExtension struct {
	Name string `json:"name"`
}

type Database struct {
	Name       string              `json:"name"`
	Extensions []DatabaseExtension `json:"extensions"`
}

func (d Database) Create() error {
	log.Println("Create database:", d.Name)

	_, err := exec.Command("sudo", "-u", "postgres", "createdb", d.Name).CombinedOutput()

	if err != nil {
		log.Println("Database exists")
	}

	for _, extension := range d.Extensions {
		sql_command := fmt.Sprintf("psql %s -c 'create extension %s'", d.Name, extension.Name)
		_, err := exec.Command("su", "postgres", "-c", sql_command).CombinedOutput()
		if err != nil {
			log.Println("Extension exists")
		}
	}

	return nil
}
