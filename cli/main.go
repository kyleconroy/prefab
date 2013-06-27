package main

import (
	"encoding/json"
	"flag"
	"github.com/stackmachine/stackgo"
	"io/ioutil"
	"log"
)

func main() {
	flag.Parse()

	manifests := []stackgo.Manifest{}

	for _, path := range flag.Args() {
		contents, err := ioutil.ReadFile(path)

		if err != nil {
			log.Fatal(err)
		}

		var manifest stackgo.Manifest

		err = json.Unmarshal(contents, &manifest)

		if err != nil {
			log.Fatal(err)
		}

		manifests = append(manifests, manifest)
	}

	for _, manifest := range manifests {
		err := manifest.Converge()

		if err != nil {
			log.Fatal(err)
		}
	}
}
