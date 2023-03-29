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
	module := flag.String("module", "", "If the go module is not in the repository's root directory, you can specify the path here")
	prefix := flag.String("prefix", "", "Prefix the result tag and strip the prefix from the existing tags when searching")
	flag.Parse()

	directory := flag.Arg(0)
	if directory == "" {
		log.Fatalln("No directory specified")
	}

	if *module == "" {
		*module = directory
	}

	bumpType := bumpy.BumpTypePatch

	if *minor {
		bumpType = bumpy.BumpTypeMinor
		log.Printf("Bumping minor version in %s", directory)
	} else {
		log.Printf("Bumping patch version in %s", directory)
	}

	config := bumpy.BumpConfig{
		Prefix:          *prefix,
		Directory:       directory,
		ModuleDirectory: *module,
		Type:            bumpType,
		RemotePush:      *push,
	}

	newTag, err := bumpy.Bump(config)
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println(newTag)
}
