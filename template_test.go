package stackgo

import (
	"io/ioutil"
	"testing"
)

func TestTemplate(t *testing.T) {
	data := map[string]interface{}{
		"bat": "bar",
	}

	tmpl := Template{
		Path:   "output/test.txt",
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
