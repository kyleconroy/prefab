package stackgo

import (
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type Template struct {
	Path   string                 `json:"path"`
	Source string                 `json:"source"`
	Data   map[string]interface{} `json:"data"`
	Mode   uint64                 `json:"mode"`
}

func (t *Template) Create() error {
	log.Println("Create file:", t.Path)

	tmpl, err := template.ParseFiles(t.Source)

	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(t.Path), 0777)

	if err != nil {
		return err
	}

	handle, err := os.Create(t.Path)

	if err != nil {
		return err
	}

	return tmpl.Execute(handle, t.Data)
}
