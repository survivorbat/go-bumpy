package main

import (
	"flag"
	"fmt"
	"github.com/survivorbat/go-bumpy"
	"log"
)

func main() {
	log.SetFlags(0)
	minor := flag.Bool("minor", false, "Whether to bump the minor version, instead of the patch")
	push := flag.String("push", "", "Whether to push the new tag to a remote, and if so, which one")
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

	newTag, err := bumpy.Bump(directory, bumpType, *push)
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println(newTag)
}
