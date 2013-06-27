package main

import (
	"encoding/json"
	"flag"
	"github.com/stackmachine/stackgo"
	"log"
	"os"
)

func main() {
	var command string

	flag.StringVar(&command, "command", "analyze", "Command to run")
	flag.Parse()

	if command == "analyze" {
		manifest, err := stackgo.Analyze()

		if err != nil {
			log.Fatal(err)
		}

		b, err := json.Marshal(manifest)

		if err != nil {
			log.Fatal(err)
		}

		os.Stdout.Write(b)
	}

}
