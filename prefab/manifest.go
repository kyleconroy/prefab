package prefab

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
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
		return "", err
	}

	return f.Name(), nil
}

type Manifest struct {
	SourceLists     []SourceList             `json:"source_lists"`
	Packages        []Package                `json:"apt_packages"`
	Directories     []Directory              `json:"directories"`
	Templates       []Template               `json:"templates"`
	PackageArchives []PersonalPackageArchive `json:"personal_package_archives"`
	Tarballs        []Tarball                `json:"tarballs"`
	Users           []User                   `json:"users"`
	Symlinks        []Symlink                `json:"symlinks"`
	Services        []Service                `json:"services"`
	Databases       []Database               `json:"postgres_databases"`
}

func Analyze() (Manifest, error) {
	path := "/etc/apt/sources.list.d"

	_, err := ioutil.ReadDir(path)

	if err != nil {
		return Manifest{}, err
	}

	return Manifest{}, nil
}

func (m Manifest) Begin() error {
	err := os.MkdirAll("/var/prefab", 0777)

	if err != nil {
		return err
	}

	info, err := os.Stat("/var/prefab/apt-update")

	if os.IsNotExist(err) {
		_, err = os.Create("/var/prefab/apt-update")

		if err != nil {
			return err
		}

		log.Println("Run `apt-get update`")
		out, err := exec.Command("apt-get", "update").CombinedOutput()

		if err != nil {
			log.Println(string(out))
			return err
		}

		return nil
	}

	// If the ModTime on this file is older than a week, rerun
	if info.ModTime().Before(time.Now().AddDate(0, 0, -7)) {

		log.Println("Run `apt-get update`")
		out, err := exec.Command("apt-get", "update").CombinedOutput()

		if err != nil {
			log.Println(string(out))
			return err
		}

	}

	return nil

}

func (m Manifest) Converge() error {
	for _, user := range m.Users {
		err := user.Create()

		if err != nil {
			return err
		}
	}

	err := m.Begin()

	if err != nil {
		return err
	}

	apt_update_needed := false

	for _, slist := range m.SourceLists {
		created, err := slist.Install()

		if err != nil {
			return err
		}

		if created {
			apt_update_needed = true
		}
	}

	// If there are Personal Package Archives to install,
	// make sure that the `add-apt-repository` command is available
	if len(m.PackageArchives) > 0 {
		pkg := Package{Name: "python-software-properties"}
		err := pkg.Install()

		if err != nil {
			return err
		}
	}

	for _, ppa := range m.PackageArchives {
		created, err := ppa.Install()

		if err != nil {
			return err
		}

		if created {
			apt_update_needed = true
		}
	}

	// Replace this with notifications eventually
	if apt_update_needed {
		log.Println("Run `apt-get update`")
		out, err := exec.Command("apt-get", "update").CombinedOutput()

		if err != nil {
			log.Println(string(out))
			return err
		}
	}

	for _, pack := range m.Packages {
		err := pack.Install()

		if err != nil {
			return err
		}
	}

	for _, tarball := range m.Tarballs {
		err := tarball.Unpack()

		if err != nil {
			return err
		}
	}

	for _, dir := range m.Directories {
		err := dir.Create()

		if err != nil {
			return err
		}
	}

	for _, tmpl := range m.Templates {
		err := tmpl.Create()

		if err != nil {
			return err
		}
	}

	for _, symlink := range m.Symlinks {
		err := symlink.Create()

		if err != nil {
			return err
		}
	}

	for _, service := range m.Services {
		err := service.Create()

		if err != nil {
			return err
		}
	}

	for _, db := range m.Databases {
		err := db.Create()

		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Manifest) Add(other Manifest) {
	m.SourceLists = append(m.SourceLists, other.SourceLists...)
	m.Packages = append(m.Packages, other.Packages...)
	m.Directories = append(m.Directories, other.Directories...)
	m.Templates = append(m.Templates, other.Templates...)
	m.PackageArchives = append(m.PackageArchives, other.PackageArchives...)
	m.Tarballs = append(m.Tarballs, other.Tarballs...)
	m.Users = append(m.Users, other.Users...)
	m.Services = append(m.Services, other.Services...)
	m.Databases = append(m.Databases, other.Databases...)
	m.Symlinks = append(m.Symlinks, other.Symlinks...)
}
