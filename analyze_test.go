package stackgo

import (
	"io/ioutil"
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
	slist, err := ParseSourceList("fixtures/test.list")

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

func TestTemplate(t *testing.T) {
	data := map[string]interface{}{
		"bat": "bar",
	}

	tmpl := Template{
		Path:   "test.txt",
		Source: "fixtures/template.txt",
		Data:   data,
	}

	err := tmpl.Create()

	if err != nil {
		t.Fatal(err)
	}

	contents, err := ioutil.ReadFile(tmpl.Path)

	if err != nil {
		t.Fatal(err)
	}

	if string(contents) != "value: bar\n" {
		t.Fatal("File contents:", string(contents))
	}
}
