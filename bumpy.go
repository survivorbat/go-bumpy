package bumpy

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"golang.org/x/mod/modfile"
	"log"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
)

type BumpType int

const BumpTypeMinor BumpType = 1
const BumpTypePatch BumpType = 2

// getTags returns a list of semver tags in the repository, these versions are ordered
// from lowest to highest
func getTags(prefix string, repo *git.Repository) (semver.Collection, error) {
	tags, err := repo.Tags()
	if err != nil {
		log.Printf("Failed to get tags: %s\n", err.Error())
		return semver.Collection{}, err
	}

	var resultTags semver.Collection

	_ = tags.ForEach(func(tag *plumbing.Reference) error {
		tagName := tag.Name().Short()

		// Skip anything that does not contain the prefix
		if !strings.Contains(tagName, prefix) {
			return nil
		}

		tagName = strings.Replace(tagName, prefix, "", 1)

		result, err := semver.NewVersion(tagName)
		if err != nil {
			return nil
		}

		resultTags = append(resultTags, result)
		return nil
	})

	sort.Sort(resultTags)

	return resultTags, nil
}

// moduleVersionTag is a regex that matches a module version tag at the end of a
// go.mod module statement
var moduleVersionTag = regexp.MustCompile(`^v\d+$`)

// getModuleVersion returns the module version from the go.mod file in the given directory,
// if no module version is found, an empty string is returned
func getModuleVersion(directory string) (string, error) {
	expectedModPath := path.Join(directory, "go.mod")

	// Check if the go.mod file exists, if not, return an empty string, we don't log it
	// because it's not an error if there is no go.mod file
	fileContents, err := os.ReadFile(expectedModPath)
	if err != nil {
		log.Printf("Not able to find %s, ignoring", expectedModPath)
		return "", nil
	}

	// If there is a go.mod file, we are going to log the error because it's unexpected
	// and we want the user to know about it
	moduleFile, err := modfile.Parse(expectedModPath, fileContents, nil)
	if err != nil {
		log.Printf("Failed to parse %s: %s\n", expectedModPath, err.Error())
		return "", err
	}

	// Split the module statement into its parts and take the final segment
	splitPath := strings.Split(moduleFile.Module.Mod.Path, "/")
	versionTag := splitPath[len(splitPath)-1]

	// Check if the final segment is the version tag, if not, return an empty string
	if moduleVersionTag.Match([]byte(versionTag)) {
		return versionTag[1:], nil
	}

	return "", nil
}

type BumpConfig struct {
	Prefix          string
	Directory       string
	ModuleDirectory string
	Type            BumpType
	RemotePush      string
}

// Bump creates a new tag for the given repository, the version identifier is determined by:
// - The current version of the module, as defined in the go.mod file (the v2 or v3 part)
// - The latest tag in the repository
// - The Bump type (minor or patch)
// - The default of v0.0.0 if there are no tags or go.mod files
func Bump(bumpConfig BumpConfig) (string, error) {
	repo, err := git.PlainOpen(bumpConfig.Directory)
	if err != nil {
		log.Printf("Failed to open repository '%s': %s\n", bumpConfig.Directory, err.Error())
		return "", err
	}

	tags, err := getTags(bumpConfig.Prefix, repo)
	if err != nil {
		return "", err
	}

	moduleVersion, err := getModuleVersion(bumpConfig.ModuleDirectory)
	if err != nil {
		return "", err
	}

	var latestTag *semver.Version
	var shouldBump = true

	switch {
	// If we have zero data, we're starting from scratch and should use v0.0.0
	case len(tags) == 0 && moduleVersion == "":
		latestTag = semver.MustParse("v0.0.0")
		shouldBump = false

	// If we have zero data, but we have a module version, we should use that with two zeros
	case len(tags) == 0 && moduleVersion != "":
		latestTag = semver.MustParse(fmt.Sprintf("v%s.0.0", moduleVersion))
		shouldBump = false

	// If we have tags, but no module version, we should use the latest tag
	case len(tags) > 0 && moduleVersion == "":
		latestTag = tags[len(tags)-1]

	// If we have tags and a module version, we should use the latest tag with the module version
	case len(tags) > 0 && moduleVersion != "":
		var filteredTags semver.Collection
		for _, tag := range tags {
			if strings.HasPrefix(tag.String(), moduleVersion) {
				filteredTags = append(filteredTags, tag)
			}
		}

		if len(filteredTags) > 0 {
			latestTag = filteredTags[len(filteredTags)-1]
		} else {
			latestTag = semver.MustParse(fmt.Sprintf("v%s.0.0", moduleVersion))
			shouldBump = false
		}
	}

	// We should bump the version if we have a tag, or if we have an existing module version
	if shouldBump {
		switch bumpConfig.Type {
		case BumpTypeMinor:
			newTag := latestTag.IncMinor()
			latestTag = &newTag
		case BumpTypePatch:
			newTag := latestTag.IncPatch()
			latestTag = &newTag
		}
	}

	headRef, err := repo.Head()
	if err != nil {
		log.Printf("Failed to get HEAD: %s\n", err.Error())
		return "", err
	}

	newTag := fmt.Sprintf("%sv%s", bumpConfig.Prefix, latestTag.String())

	log.Printf("Creating tag %s\n", newTag)
	if _, err := repo.CreateTag(newTag, headRef.Hash(), nil); err != nil {
		log.Printf("Failed to create tag: %s\n", err.Error())
		return "", err
	}

	if bumpConfig.RemotePush != "" {
		ref := fmt.Sprintf("refs/tags/%s:refs/tags/%s", newTag, newTag)
		options := &git.PushOptions{
			RemoteName: bumpConfig.RemotePush,
			RefSpecs: []config.RefSpec{
				config.RefSpec(ref),
			},
		}

		log.Printf("Pushing tag %s to %s\n", newTag, bumpConfig.RemotePush)
		if err := repo.Push(options); err != nil {
			log.Printf("Failed to pushRemote tag: %s\n", err.Error())
		}
	}

	return newTag, nil
}
