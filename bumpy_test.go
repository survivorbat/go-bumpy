package bumpy

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func fatalIf(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err.Error())
	}
}

func setupRepo(t *testing.T, directory string, tags []string) *git.Repository {
	err := os.Mkdir(directory, 0755)
	fatalIf(t, err)

	repo, err := git.PlainInit(directory, false)
	fatalIf(t, err)

	fullPath := path.Join(directory, "README.md")

	err = os.WriteFile(fullPath, []byte("Hello world"), 0644)
	fatalIf(t, err)

	tree, err := repo.Worktree()
	fatalIf(t, err)

	status, err := tree.Status()
	fatalIf(t, err)

	for filePath := range status {
		_, err = tree.Add(filePath)
		fatalIf(t, err)
	}

	commit, err := tree.Commit("Initial commit", &git.CommitOptions{})
	fatalIf(t, err)

	_, err = repo.CommitObject(commit)
	fatalIf(t, err)

	for _, tag := range tags {
		_, err = repo.CreateTag(tag, commit, nil)
		fatalIf(t, err)
	}

	return repo
}

func TestBump_ReturnsExpectedVersionWithModuleFile(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		existing    []string
		moduleName  string
		versionBump BumpType
		expected    string
	}{
		"no tags, no module, patch": {
			existing:    []string{},
			expected:    "v0.0.1",
			versionBump: BumpTypePatch,
		},
		"no tags, a module, patch": {
			existing:    []string{},
			moduleName:  "github.com/foobar/vaz",
			expected:    "v0.0.1",
			versionBump: BumpTypePatch,
		},
		"no tags, a module version, patch": {
			existing:    []string{},
			moduleName:  "github.com/survivorbat/vv-bumpy/v5",
			expected:    "v5.0.1",
			versionBump: BumpTypePatch,
		},
		"multiple tags, no module version, patch": {
			existing:    []string{"v1.0.0", "v2.0.1", "v3.0.2", "v4.0.3"},
			expected:    "v4.0.4",
			versionBump: BumpTypePatch,
		},
		"multiple tags, a module version, patch": {
			existing:    []string{"v1.0.0", "v2.0.1", "v3.0.2", "v4.0.3", "v5.0.1", "v6.4.2"},
			moduleName:  "github.com/survivorbat/go-bumpy/v5",
			expected:    "v5.0.2",
			versionBump: BumpTypePatch,
		},
		"multiple tags, a module version, minor": {
			existing:    []string{"v1.0.0", "v2.0.1", "v3.0.2", "v4.0.3", "v5.0.1", "v6.4.2"},
			moduleName:  "github.com/survivorbat/go-bumpy/v5",
			expected:    "v5.1.0",
			versionBump: BumpTypeMinor,
		},
		"multiple tags, a new module version, patch": {
			existing:    []string{"v1.0.0", "v2.0.1", "v3.0.2", "v4.0.3", "v5.0.1", "v6.4.2"},
			moduleName:  "github.com/survivorbat/go-bumpy/v7",
			expected:    "v7.0.1",
			versionBump: BumpTypePatch,
		},
		"multiple tags, a new module version, minor": {
			existing:    []string{"v1.0.0", "v2.0.1", "v3.0.2", "v4.0.3", "v5.0.1", "v6.4.2"},
			moduleName:  "github.com/survivorbat/go-bumpy/v7",
			expected:    "v7.1.0",
			versionBump: BumpTypeMinor,
		},
	}

	for name, testData := range tests {
		testData := testData
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			directory := path.Join(t.TempDir(), "test")
			repo := setupRepo(t, directory, testData.existing)

			if testData.moduleName != "" {
				moduleContents := fmt.Sprintf("module %s\n\ngo 1.19\n", testData.moduleName)
				err := os.WriteFile(path.Join(directory, "go.mod"), []byte(moduleContents), 0644)
				fatalIf(t, err)
			}

			// Act
			err := Bump(directory, testData.versionBump)

			// Assert
			assert.NoError(t, err)

			result, err := repo.Tag(testData.expected)
			if assert.NoError(t, err) {
				assert.Equal(t, testData.expected, result.Name().Short())
			}
		})
	}
}
