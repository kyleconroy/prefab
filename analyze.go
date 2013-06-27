package stackgo

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type Directory struct {
	Path string `json:"path"`
}

type Package struct {
	Name string `json:"name"`
}

type Archive struct {
	Name string `json:"name"`
}

type Manifest struct {
	Directories []Directory `json:"directories"`
	Packages    []Package   `json:"packages"`
	Archives    []Archive   `json:"personal_package_archives"`
}

func Analyze() (Manifest, error) {
	return Manifest{}, nil
}

// Install Postgres
func InstallPostgres(version string) error {
	log.Println("Install postgres")

	path := "/etc/apt/sources.list.d/pgdg.list"

	_, err := os.Stat("/etc/apt/sources.list.d/pgdg.list")

	if os.IsNotExist(err) {
		line := "deb http://apt.postgresql.org/pub/repos/apt/ precise-pgdg main\n"

		err = ioutil.WriteFile(path, []byte(line), 0644)

		if err != nil {
			return err
		}

		f, err := ioutil.TempFile("", "stackmachine")

		if err != nil {
			return err
		}

		defer f.Close()

		resp, err := http.Get("http://apt.postgresql.org/pub/repos/apt/ACCC4CF8.asc")

		if err != nil {
			return err
		}

		defer resp.Body.Close()

		_, err = io.Copy(f, resp.Body)

		if err != nil {
			return err
		}

		out, err := exec.Command("apt-key", "add", f.Name()).Output()

		if err != nil {
			log.Println(string(out))
			return err
		}

		out, err = exec.Command("apt-get", "update").Output()

		if err != nil {
			log.Println(string(out))
			return err
		}

	}

	out, err := exec.Command("apt-get", "install", "-y", "postgresql-"+version).Output()

	if err != nil {
		log.Println(string(out))
		return err
	}

	return nil
}

func Converge() error {
	var manifest Manifest

	contents, err := ioutil.ReadFile("manifest.json")

	if err != nil {
		return err
	}

	err = json.Unmarshal(contents, &manifest)

	if err != nil {
		return err
	}

	err = InstallPostgres("9.2")

	if err != nil {
		return err
	}

	return nil
}
