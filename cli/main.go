package main

import (
	"encoding/json"
	"github.com/stackmachine/stackgo"
	"log"
	"os"
)

func main() {
	if false {
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
	
	err := stackgo.Converge()

	if err != nil {
		log.Fatal(err)
	}
}
