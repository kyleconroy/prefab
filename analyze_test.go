package stackgo

import (
	"testing"
)

func TestPackageRepositoryPath(t *testing.T) {
	repo := PackageRepository{
		Filename: "foo.list",
	}

	if repo.Path() != "/etc/apt/sources.list.d/foo.list" {
		t.Fatal("Path not equal", repo.Path())
	}
}

func TestParseSourceList(t *testing.T) {
	packrepo, err := ParseSourceList("fixtures/test.list")

	if err != nil {
		t.Fatal(err)
	}

	if packrepo.Filename != "test.list" {
		t.Fatalf("Source list filename: %s not %s", packrepo.Filename, "test")
	}
}

func TestPackageRepositorySource(t *testing.T) {
	repo := PackageRepository{
		Uri:          "http://example.com",
		Distribution: "precise-foo",
		Components:   []string{"main", "foo"},
	}

	if repo.Entry() != "deb http://example.com precise-foo main foo" {
		t.Fatal("Path not equal", repo.Path())
	}
}
