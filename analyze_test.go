package stackgo

import ("testing")


func TestPackageRepositoryPath(t *testing.T) {
	repo := PackageRepository{
		Name: "foo", 
	}

	if repo.Path() != "/etc/apt/sources.list.d/foo.list" {
		t.Fatal("Path not equal", repo.Path())
	}
}


func TestPackageRepositorySource(t *testing.T) {
	repo := PackageRepository{
		Uri: "http://example.com",
		Distribution: "precise-foo",
		Components: []string{"main foo"},
	}

	if repo.Entry() != "deb http://example.com precise-foo main foo" {
		t.Fatal("Path not equal", repo.Path())
	}
}
