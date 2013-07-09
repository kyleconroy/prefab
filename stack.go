package main

import (
	"encoding/json"
	"flag"
	"github.com/stackmachine/stackgo/stackgo"
	"io/ioutil"
	"log"
	"os"
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

	var manifest stackgo.Manifest

	for _, path := range flag.Args() {
		contents, err := ioutil.ReadFile(path)

		if err != nil {
			log.Fatal(err)
		}

		var userManifest stackgo.Manifest

		err = json.Unmarshal(contents, &userManifest)

		if err != nil {
			log.Fatal(err)
		}

		manifest.Add(userManifest)
	}

	start := time.Now()

	err := manifest.Converge()

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Took %s", time.Since(start))
}
