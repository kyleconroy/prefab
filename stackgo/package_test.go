package stackgo

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestSourceListPath(t *testing.T) {
	slist := SourceList{
		Filename: "foo.list",
	}

	if slist.Path() != "/etc/apt/sources.list.d/foo.list" {
		t.Fatal("Path not equal", slist.Path())
	}
}

func TestParseSourceList(t *testing.T) {
	err := os.MkdirAll("test", 0777)

	if err != nil {
		t.Fatal(err)
	}

	srclist := `deb http://apt.postgresql.org/pub/repos/apt/ precise-pgdg main
deb http://apt.postgresql.org/pub/repos/apt/ precise-pgdg universe foo
deb-src http://apt.postgresql.org/pub/repos/apt/ precise-pgdg main`

	err = ioutil.WriteFile("test/test.list", []byte(srclist), 0644)

	if err != nil {
		t.Fatal(err)
	}

	slist, err := ParseSourceList("test/test.list")

	if err != nil {
		t.Fatal(err)
	}

	if slist.Filename != "test.list" {
		t.Fatalf("Source list filename: %s not %s", slist.Filename, "test")
	}

	if len(slist.Sources) != 2 {
		t.Fatal("Didn't parse the correct number of sources")
	}

	source := slist.Sources[1]

	if source.Uri != "http://apt.postgresql.org/pub/repos/apt/" {
		t.Fatal("Incorrect source uri")
	}
}

func TestSource(t *testing.T) {
	source := Source{
		Uri:          "http://example.com",
		Distribution: "precise-foo",
		Components:   []string{"main", "foo"},
	}

	if source.Entry() != "deb http://example.com precise-foo main foo" {
		t.Fatal("Source entry incorrect: ", source.Entry())
	}
}

func TestPPA(t *testing.T) {
	ppa := PersonalPackageArchive{Name: "foo", Owner: "bar"}

	if ppa.Path() != "/etc/apt/sources.list.d/bar-foo-precise.list" {
		t.Fatal("PPA path is incorrect:", ppa.Path())
	}
}
