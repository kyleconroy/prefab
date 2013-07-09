package stackgo

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestTemplate(t *testing.T) {
	data := map[string]interface{}{
		"bat": "bar",
	}

	err := os.MkdirAll("test", 0777)

	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile("test/template.txt", []byte("value: {{ .bat }}"), 0644)

	if err != nil {
		t.Fatal(err)
	}

	tmpl := Template{
		Path:   "test/template_output.txt",
		Source: "test/template.txt",
		Data:   data,
	}

	err = tmpl.Create()

	if err != nil {
		t.Fatal(err)
	}

	contents, err := ioutil.ReadFile(tmpl.Path)

	if err != nil {
		t.Fatal(err)
	}

	if string(contents) != "value: bar" {
		t.Fatal("File contents:", string(contents))
	}
}
