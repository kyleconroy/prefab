package stackgo

import (
)

type Directory struct {
	Path string `json:"path"`
}

type Package struct {
	Name string `json:"name"`
}

type Manifest struct {
	Directories []Directory `json:"directories"`
	Packages    []Package   `json:"packages"`
}

func Analyze() (Manifest, error) {
	return Manifest{}, nil
}
