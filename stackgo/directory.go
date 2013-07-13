package stackgo

import (
	"log"
	"os"
)

type Directory struct {
	Path string `json:"path"`
	Mode uint32 `json:"mode"`
}

func (d *Directory) Create() error {
	log.Println("Create directory:", d.Path)

	err := os.MkdirAll(d.Path, 0777)
	return err
}
