package prefab

import (
	"log"
	"os"
	"os/exec"
)

type Tarball struct {
	Url  string `json:"url"`
	Path string `json:"path"`
}

func (t Tarball) Unpack() error {
	log.Printf("Unpacking %s to %s", t.Url, t.Path)

	finfo, err := os.Stat(t.Path)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Directory exists.")
			return nil
		}
	}

	if finfo != nil && finfo.IsDir() {
		log.Printf("Directory exists.")
		return nil
	}

	log.Printf("making dir %s", t.Path)
	err = os.MkdirAll(t.Path, 0777)
	if err != nil {
		return err
	}

	tmpfile, _ := download(t.Url)
	// Shell out for now
	_, err = exec.Command("tar", "-xzf", tmpfile, "-C", t.Path).CombinedOutput()
	return err
	// TODO: cleanup tempfiles
}
