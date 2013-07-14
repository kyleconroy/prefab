package prefab

import (
	"log"
	"os"
)

type Symlink struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

func (s Symlink) Create() error {
	log.Printf("Create symlink: %s -> %s", s.Source, s.Destination)
	var source, _ = os.Readlink(s.Destination)
	if source == s.Source {
		return nil
	}
	//if err != nil {
	//	return err
	//}

	return os.Symlink(s.Source, s.Destination)
}
