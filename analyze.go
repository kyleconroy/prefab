package stackgo

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func download(uri string) (string, error) {
		f, err := ioutil.TempFile("", "download")

		if err != nil {
			return "", err
		}

		defer f.Close()

		resp, err := http.Get(uri)

		if err != nil {
			return "", err
		}

		defer resp.Body.Close()

		_, err = io.Copy(f, resp.Body)

		if err != nil {
			return "",err
		}

		return f.Name(), nil

}

type Package struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (p Package) Install() error {
	log.Println("Install package: ", p.Name)

	out, err := exec.Command("apt-get", "install", "-y", p.Name).Output()

	if err != nil {
		log.Println(out)
	}

	return err
}

type PackageRepository struct {
	Name         string   `json:"name"`
	Uri          string   `json:"uri"`
	Distribution string   `json:"distribution"`
	KeyURI       string   `json:"gpg_key_uri"`
	Components   []string `json:"components"`
}

func (pr *PackageRepository) Path() string {
	return "/etc/apt/sources.list.d/" + pr.Name + ".list"
}

func (pr *PackageRepository) Entry() string {
	// TODO: Add component support
	return fmt.Sprintf("deb %s %s main", pr.Uri, pr.Distribution)
}

// Return created, error
func (pr *PackageRepository) InstallSourceList() (bool, error) {
	_, err := os.Stat(pr.Path())

	if os.IsNotExist(err) {
		log.Println("Install source list")

		err = ioutil.WriteFile(pr.Path(), []byte(pr.Entry()), 0644)

		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}



func (pr *PackageRepository) InstallKey() error {
	// TODO: Figure out cache module
	keyPath, err := download(pr.KeyURI)

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


func (pr *PackageRepository) Install() error {
	source_added, err := pr.InstallSourceList()

	if err != nil {
		return err
	}

	//Fix this
	if !source_added  {
		return nil
	}

	err = pr.InstallKey()

	if err != nil {
		return err
	}

	out, err := exec.Command("apt-get", "update").Output()

	if err != nil {
		log.Println(string(out))
		return err
	}

	return nil
}

type Manifest struct {
	PackageRepos []PackageRepository `json:"package_repositories"`
	Packages     []Package           `json:"packages"`
}

func Analyze() (Manifest, error) {
	return Manifest{}, nil
}

func (m Manifest) Converge() error {
	for _, packrepo := range m.PackageRepos {
		err := packrepo.Install()

		if err != nil {
			return err
		}
	}

	for _, pack := range m.Packages {
		err := pack.Install()

		if err != nil {
			return err
		}
	}

	return nil
}
