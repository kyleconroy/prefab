package prefab

import (
	"log"
	"os"
	"syscall"
)

type Directory struct {
	Path string `json:"path"`
	Mode uint32 `json:"mode"`
}

func (d *Directory) Create() error {
	log.Println("Create directory:", d.Path)

	oldMode := syscall.Umask(000)
	err := os.MkdirAll(d.Path, os.ModeDir|0777)
	syscall.Umask(oldMode)
	return err
}
