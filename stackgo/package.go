package stackgo

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Package struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (p Package) Install() error {
	log.Println("Install package:", p.Name)

	out, err := exec.Command("apt-cache", "policy", "-q", p.Name).Output()

	if err != nil {
		log.Println(string(out))
		return err
	}

	for _, line := range strings.Split(string(out), "\n") {
		sline := strings.TrimSpace(line)

		if strings.HasPrefix(sline, "Installed: ") {
			if !strings.HasSuffix(sline, "(none)") {
				return nil
			}
		}
	}

	out, err = exec.Command("apt-get", "install", "-y", p.Name).Output()

	if err != nil {
		log.Println(string(out))
	}

	return err
}

type Source struct {
	Uri          string   `json:"uri"`
	Distribution string   `json:"distribution"`
	Components   []string `json:"components"`
}

func (s *Source) Entry() string {
	// TODO: Add component support
	entry := "deb " + s.Uri + " " + s.Distribution

	for _, component := range s.Components {
		entry = entry + " " + component
	}

	return entry
}

type SourceList struct {
	Filename string   `json:"filename"`
	KeyURI   string   `json:"key_uri"`
	Sources  []Source `json:"sources"`
}

func (sl *SourceList) Path() string {
	return "/etc/apt/sources.list.d/" + sl.Filename
}

// Return created, error
func (sl *SourceList) InstallSources() (bool, error) {
	log.Println("Install archive:", sl.Path())

	_, err := os.Stat(sl.Path())

	if os.IsNotExist(err) {

		var body string

		for _, source := range sl.Sources {
			body = body + source.Entry() + "\n"
		}

		err = ioutil.WriteFile(sl.Path(), []byte(body), 0644)

		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func (sl *SourceList) InstallKey() error {
	if sl.KeyURI == "" {
		return nil
	}

	// TODO: Figure out cache module
	keyPath, err := download(sl.KeyURI)

	if err != nil {
		return err
	}

	log.Println("Install key")

	out, err := exec.Command("apt-key", "add", keyPath).Output()

	if err != nil {
		log.Println(string(out))
		return err
	}

	return nil
}

func (sl *SourceList) Install() (bool, error) {
	source_added, err := sl.InstallSources()

	if err != nil {
		return source_added, err
	}

	//Fix this
	if !source_added {
		return false, nil
	}

	err = sl.InstallKey()

	if err != nil {
		return false, err
	}

	return true, nil
}

type PersonalPackageArchive struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

func (ppa *PersonalPackageArchive) Path() string {
	return fmt.Sprintf("/etc/apt/sources.list.d/%s-%s-precise.list", ppa.Owner, ppa.Name)
}

func (ppa *PersonalPackageArchive) Install() (bool, error) {
	_, err := os.Stat(ppa.Path())

	id := fmt.Sprintf("ppa:%s/%s", ppa.Owner, ppa.Name)

	log.Println("Install archive:", id)

	if os.IsNotExist(err) {
		out, err := exec.Command("add-apt-repository", "-y", id).Output()

		if err != nil {
			log.Println(string(out))
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func ParseSourceList(path string) (SourceList, error) {
	b, err := ioutil.ReadFile(path)

	if err != nil {
		return SourceList{}, err
	}

	var sources []Source

	for _, line := range strings.Split(string(b), "\n") {
		parts := strings.Split(line, " ")

		if len(parts) < 4 {
			// Incorrect source line
			continue
		}

		if parts[0] != "deb" {
			// Unsupported source line
			continue
		}

		source := Source{
			Uri:          parts[1],
			Distribution: parts[2],
			Components:   parts[:3],
		}

		sources = append(sources, source)
	}

	return SourceList{
		Filename: filepath.Base(path),
		Sources:  sources,
	}, nil
}
