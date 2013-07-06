package stackgo

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
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

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (u User) Create() error {
	log.Println("Create user:", u.Name)
	return nil
}

type Service struct {
	Name string `json:"name"`
}

func (s Service) Create() error {
	log.Println("Start service:", s.Name)
	out, err := exec.Command("service", s.Name, "start").Output()

	if err != nil {
		log.Println(string(out))
		return err
	}

	return nil
}

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

type Manifest struct {
	SourceLists     []SourceList             `json:"source_lists"`
	Packages        []Package                `json:"packages"`
	Templates       []Template               `json:"templates"`
	PackageArchives []PersonalPackageArchive `json:"personal_package_archives"`
	Users           []User                   `json:"users"`
	Services        []Service                `json:"services"`
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

func Analyze() (Manifest, error) {
	path := "/etc/apt/sources.list.d"

	_, err := ioutil.ReadDir(path)

	if err != nil {
		return Manifest{}, err
	}

	return Manifest{}, nil
}

func (m Manifest) Begin() error {
	err := os.MkdirAll("/var/stackgo", 0777)

	if err != nil {
		return err
	}

	info, err := os.Stat("/var/stackgo/apt-update")

	if os.IsNotExist(err) {
		_, err = os.Create("/var/stackgo/apt-update")

		if err != nil {
			return err
		}

		log.Println("Run `apt-get update`")
		out, err := exec.Command("apt-get", "update").Output()

		if err != nil {
			log.Println(string(out))
			return err
		}

		return nil
	}

	// If the ModTime on this file is older than a week, rerun
	if info.ModTime().Before(time.Now().AddDate(0, 0, -7)) {

		log.Println("Run `apt-get update`")
		out, err := exec.Command("apt-get", "update").Output()

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
		out, err := exec.Command("apt-get", "update").Output()

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

	for _, tmpl := range m.Templates {
		err := tmpl.Create()

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

	return nil
}

func (m *Manifest) Add(other Manifest) {
	m.SourceLists = append(m.SourceLists, other.SourceLists...)
	m.Packages = append(m.Packages, other.Packages...)
	m.Templates = append(m.Templates, other.Templates...)
	m.PackageArchives = append(m.PackageArchives, other.PackageArchives...)
	m.Users = append(m.Users, other.Users...)
	m.Services = append(m.Services, other.Services...)
}
