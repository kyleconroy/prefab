package main

import (
	"encoding/json"
	"flag"
	"github.com/stackmachine/prefab/prefab"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var manifest prefab.Manifest

	for _, path := range flag.Args() {
		contents, err := ioutil.ReadFile(path)

		if err != nil {
			log.Fatal(err)
		}

		var userManifest prefab.Manifest

		err = json.Unmarshal(contents, &userManifest)

		if err != nil {
			log.Fatal(err)
		}

		absPath, err := filepath.Abs(path)

		if err != nil {
			log.Fatal(err)
		}

		userManifest.FixPaths(absPath)

		manifest.Add(userManifest)
	}

	start := time.Now()

	err := manifest.Converge()

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Took %s", time.Since(start))
}
