package prefab

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// Require Gemfile.lock

type RubyBundle struct {
	Path string `json:"path"`
}

func (b *RubyBundle) Install() error {
	log.Println("Bundling ", b.Path)

	_, err := os.Stat("/usr/local/ruby/bin/bundle")
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Installing bundler.")
			_, _ = exec.Command("gem", "install", "bundler").CombinedOutput()
		}
	}

	pwd, _ := os.Getwd()
	os.Chdir(filepath.Dir(b.Path))
	out, err := exec.Command("bundle", "install").CombinedOutput()
	log.Println(string(out))

	if err != nil {
		return err
	}
	os.Chdir(pwd)

	return nil
}
