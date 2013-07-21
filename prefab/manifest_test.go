package prefab

import (
	"testing"
)

func TestManifest(t *testing.T) {
	manifest := Manifest{Templates: []Template{Template{Source: "foo"}}}
	manifest.FixPaths("foo/bar")

	if manifest.Templates[0].Source != "foo/foo" {
		t.Fatal("Incorrect source", manifest.Templates[0].Source)
	}
}
