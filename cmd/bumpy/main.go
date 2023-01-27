package main

import (
	"flag"
	"github.com/survivorbat/go-bumpy"
	"log"
)

func main() {
	log.SetFlags(0)
	minor := flag.Bool("minor", false, "Whether to bump the minor version, instead of the patch")
	flag.Parse()

	directory := flag.Arg(0)
	if directory == "" {
		log.Fatalln("No directory specified")
	}

	bumpType := bumpy.BumpTypePatch

	if *minor {
		bumpType = bumpy.BumpTypeMinor
		log.Printf("Bumping minor version in %s", directory)
	} else {
		log.Printf("Bumping patch version in %s", directory)
	}

	if err := bumpy.Bump(directory, bumpType); err != nil {
		log.Fatalln(err.Error())
	}
}
