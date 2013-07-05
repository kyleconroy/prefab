package main

import (
	"encoding/json"
	"flag"
	"github.com/stackmachine/stackgo"
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

	manifests := []stackgo.Manifest{}

	if len(flag.Args()) == 0 {
		_, err := stackgo.Analyze()

		if err != nil {
			log.Fatal(err)
		}

		return
	}

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

	start := time.Now()

	for _, manifest := range manifests {
		err := manifest.Converge()

		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("Took %s", time.Since(start))
}
