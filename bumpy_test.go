package bumpy

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/stretchr/testify/assert"
)

func fatalIf(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err.Error())
	}
}

func setupRepo(t *testing.T, directory string, tags []string) *git.Repository {
	t.Helper()

	err := os.Mkdir(directory, 0o755)
	fatalIf(t, err)

	repo, err := git.PlainInit(directory, false)
	fatalIf(t, err)

	fullPath := path.Join(directory, "README.md")

	err = os.WriteFile(fullPath, []byte("Hello world"), 0o600)
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
			expected:    "v0.0.0",
			versionBump: BumpTypePatch,
		},
		"no tags, a module, patch": {
			existing:    []string{},
			moduleName:  "github.com/foobar/vaz",
			expected:    "v0.0.0",
			versionBump: BumpTypePatch,
		},
		"no tags, a module version, patch": {
			existing:    []string{},
			moduleName:  "github.com/survivorbat/vv-bumpy/v5",
			expected:    "v5.0.0",
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
			expected:    "v7.0.0",
			versionBump: BumpTypePatch,
		},
		"multiple tags, a new module version, minor": {
			existing:    []string{"v1.0.0", "v2.0.1", "v3.0.2", "v4.0.3", "v5.0.1", "v6.4.2"},
			moduleName:  "github.com/survivorbat/go-bumpy/v7",
			expected:    "v7.0.0",
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
				err := os.WriteFile(path.Join(directory, "go.mod"), []byte(moduleContents), 0o600)
				fatalIf(t, err)
			}

			bumpConfig := BumpConfig{
				Directory:       directory,
				ModuleDirectory: directory,
				Type:            testData.versionBump,
			}

			// Act
			result, err := Bump(bumpConfig)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, testData.expected, result)

			_, err = repo.Tag(testData.expected)
			assert.NoError(t, err)
		})
	}
}

func TestBump_ReturnsExpectedVersionWithModuleFileInADifferentDirectory(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		existing    []string
		moduleName  string
		versionBump BumpType
		expected    string
	}{
		"no tags, no module, patch": {
			existing:    []string{},
			expected:    "v0.0.0",
			versionBump: BumpTypePatch,
		},
		"no tags, a module, patch": {
			existing:    []string{},
			moduleName:  "github.com/foobar/vaz",
			expected:    "v0.0.0",
			versionBump: BumpTypePatch,
		},
		"no tags, a module version, patch": {
			existing:    []string{},
			moduleName:  "github.com/survivorbat/vv-bumpy/v5",
			expected:    "v5.0.0",
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
			expected:    "v7.0.0",
			versionBump: BumpTypePatch,
		},
		"multiple tags, a new module version, minor": {
			existing:    []string{"v1.0.0", "v2.0.1", "v3.0.2", "v4.0.3", "v5.0.1", "v6.4.2"},
			moduleName:  "github.com/survivorbat/go-bumpy/v7",
			expected:    "v7.0.0",
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

			module := path.Join(t.TempDir(), "module")
			err := os.Mkdir(module, 0o755)
			fatalIf(t, err)

			if testData.moduleName != "" {
				moduleContents := fmt.Sprintf("module %s\n\ngo 1.19\n", testData.moduleName)
				err := os.WriteFile(path.Join(module, "go.mod"), []byte(moduleContents), 0o600)
				fatalIf(t, err)
			}

			bumpConfig := BumpConfig{
				Directory:       directory,
				ModuleDirectory: module,
				Type:            testData.versionBump,
			}

			// Act
			result, err := Bump(bumpConfig)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, testData.expected, result)

			_, err = repo.Tag(testData.expected)
			assert.NoError(t, err)
		})
	}
}

func TestBump_ReturnsExpectedVersionWithModuleFileAndPrefix(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		existing    []string
		moduleName  string
		versionBump BumpType
		prefix      string
		expected    string
	}{
		"no tags, no module, patch": {
			existing:    []string{},
			expected:    "abc/v0.0.0",
			prefix:      "abc/",
			versionBump: BumpTypePatch,
		},
		"no tags, a module, patch": {
			existing:    []string{},
			moduleName:  "github.com/foobar/vaz",
			expected:    "abc/v0.0.0",
			prefix:      "abc/",
			versionBump: BumpTypePatch,
		},
		"no tags, a module version, patch": {
			existing:    []string{},
			moduleName:  "github.com/survivorbat/vv-bumpy/v5",
			expected:    "abc/v5.0.0",
			prefix:      "abc/",
			versionBump: BumpTypePatch,
		},
		"multiple tags, no module version, patch": {
			existing:    []string{"abc/v1.0.0", "abc/v2.0.1", "abc/v3.0.2", "abc/v4.0.3"},
			expected:    "abc/v4.0.4",
			prefix:      "abc/",
			versionBump: BumpTypePatch,
		},
		"multiple tags, a module version, patch": {
			existing:    []string{"abc/v1.0.0", "abc/v2.0.1", "abc/v3.0.2", "abc/v4.0.3", "abc/v5.0.1", "abc/v6.4.2"},
			moduleName:  "github.com/survivorbat/go-bumpy/v5",
			expected:    "abc/v5.0.2",
			prefix:      "abc/",
			versionBump: BumpTypePatch,
		},
		"multiple tags, a module version, minor": {
			existing:    []string{"abc/v1.0.0", "abc/v2.0.1", "abc/v3.0.2", "abc/v4.0.3", "abc/v5.0.1", "abc/v6.4.2"},
			moduleName:  "github.com/survivorbat/go-bumpy/v5",
			expected:    "abc/v5.1.0",
			prefix:      "abc/",
			versionBump: BumpTypeMinor,
		},
		"multiple tags, a module version, minor, different prefix": {
			existing:    []string{"def/v1.0.0", "def/v2.0.1", "def/v3.0.2", "def/v4.0.3", "def/v5.0.1", "def/v6.4.2"},
			moduleName:  "github.com/survivorbat/go-bumpy/v5",
			expected:    "abc/v5.0.0",
			prefix:      "abc/",
			versionBump: BumpTypeMinor,
		},
		"multiple tags, a module version, minor, missing prefixes": {
			existing:    []string{"v1.0.0", "v2.0.1", "v3.0.2", "v4.0.3", "v5.0.1", "v6.4.2"},
			moduleName:  "github.com/survivorbat/go-bumpy/v5",
			expected:    "abc/v5.0.0",
			prefix:      "abc/",
			versionBump: BumpTypeMinor,
		},
		"multiple tags, a new module version, patch": {
			existing:    []string{"abc/v1.0.0", "abc/v2.0.1", "abc/v3.0.2", "abc/v4.0.3", "abc/v5.0.1", "abc/v6.4.2"},
			moduleName:  "github.com/survivorbat/go-bumpy/v7",
			expected:    "abc/v7.0.0",
			prefix:      "abc/",
			versionBump: BumpTypePatch,
		},
		"multiple tags, a new module version, minor": {
			existing:    []string{"def/v1.0.0", "def/v2.0.1", "def/v3.0.2", "def/v4.0.3", "def/v5.0.1", "def/v6.4.2"},
			moduleName:  "github.com/survivorbat/go-bumpy/v7",
			expected:    "def/v7.0.0",
			prefix:      "def/",
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

			module := path.Join(t.TempDir(), "module")
			err := os.Mkdir(module, 0o755)
			fatalIf(t, err)

			if testData.moduleName != "" {
				moduleContents := fmt.Sprintf("module %s\n\ngo 1.19\n", testData.moduleName)
				err := os.WriteFile(path.Join(module, "go.mod"), []byte(moduleContents), 0o600)
				fatalIf(t, err)
			}

			bumpConfig := BumpConfig{
				Directory:       directory,
				ModuleDirectory: module,
				Type:            testData.versionBump,
				Prefix:          testData.prefix,
			}

			// Act
			result, err := Bump(bumpConfig)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, testData.expected, result)

			_, err = repo.Tag(testData.expected)
			assert.NoError(t, err)
		})
	}
}

func TestBump_PushesToRemoteCorrectly(t *testing.T) {
	t.Parallel()
	// Arrange

	// We use a directory as a remote!
	moduleContents := fmt.Sprintf("module %s\n\ngo 1.19\n", "github.com/survivorbat/go-bumpy/v5")

	remoteDirectory := path.Join(t.TempDir(), "remote")
	remoteRepo := setupRepo(t, remoteDirectory, []string{"v5.0.0", "v5.0.1"})
	err := os.WriteFile(path.Join(remoteDirectory, "go.mod"), []byte(moduleContents), 0o600)
	fatalIf(t, err)

	localDirectory := path.Join(t.TempDir(), "local")
	localRepo := setupRepo(t, localDirectory, []string{"v5.0.0", "v5.0.1"})
	err = os.WriteFile(path.Join(localDirectory, "go.mod"), []byte(moduleContents), 0o600)
	fatalIf(t, err)

	remoteConfig := &config.RemoteConfig{
		Name: "origin",
		URLs: []string{remoteDirectory},
	}
	_, err = localRepo.CreateRemote(remoteConfig)
	fatalIf(t, err)

	bumpConfig := BumpConfig{
		Directory:       localDirectory,
		ModuleDirectory: localDirectory,
		Type:            BumpTypeMinor,
		RemotePush:      "origin",
	}

	// Act
	result, err := Bump(bumpConfig)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "v5.1.0", result)

	// Both repos should now have the new tag
	_, err = remoteRepo.Tag("v5.1.0")
	assert.NoError(t, err)

	_, err = localRepo.Tag("v5.1.0")
	assert.NoError(t, err)
}
