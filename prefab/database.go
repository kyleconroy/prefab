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
	log.Println("Extensions: ", d.Extensions)

	_, err := exec.Command("sudo", "-u", "postgres", "createdb", d.Name).CombinedOutput()

	if err != nil {
		log.Println("Database exists")
	}

	for _, extension := range d.Extensions {
		sql_command := fmt.Sprintf("create extension %s", extension.Name)
		_, err := exec.Command("psql", d.Name, "-c", sql_command).CombinedOutput()
		if err != nil {
			log.Println("Extension exists")
		}
	}

	return nil
}
