// +build mage

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	BuildDir = "dist"
	GoArch   = "amd64"
	Project  = "tfe-cli"
)

var (
	Default   = Setup
	Platforms = []string{"linux", "darwin", "windows"}
)

var simpleSemVer = regexp.MustCompile(`-(\d{2}\.\d{2}\.\d{1,2}-[a-z]+)`)

type Coverage mg.Namespace

// Setup the full environment.
func Setup() {
	sh.Run("go", "mod", "tidy")
}

// Run the full CI suite.
func Ci() {
	sh.RunV("goimports", "-d", ".")
	sh.RunV("golint", "./...")
	sh.RunV("go", "vet")
	mg.Deps(Test)
}

// Remove unwanted files in project (!DESTRUCTIVE!).
func Clean() {
	sh.Run("git", "clean", "-ffdx")
	sh.Run("git", "reset", "--hard")
}

func Test() {
	sh.RunV("go", "test", "-v", "-cover", "-coverprofile=coverage.out", "./...")
}

// View code coverage in text.
func (Coverage) Text() {
	sh.RunV("go", "tool", "cover", "-func=coverage.out")
}

// View code coverage in HTML.
func (Coverage) Html() {
	sh.Run("go", "tool", "cover", "-html=coverage.out")
}

// Build binaries for the targeted platforms.
func Dist() {
	tag := getTag()
	currentDir, _ := os.Getwd()
	for _, platform := range Platforms {
		builder(platform, GoArch, tag, filepath.Join(currentDir, BuildDir, Project), currentDir)
	}
}

// Publish a new GitHub release.
func Publish() {
	currentDir, _ := os.Getwd()
	buildDir := filepath.Join(currentDir, BuildDir)
	files := getBuiltArtifacts()
	assets := []string{}

	// Get the assets to publish.
	for _, file := range files {
		assets = append(assets, "-a")
		assets = append(assets, filepath.Join(buildDir, file.Name()))
	}

	// Prepare the command to run
	if len(assets) == 0 {
		fmt.Printf("Nothing to release in %s. Forgot to build?\n", buildDir)
		os.Exit(1)
	}

	// Prepare the release arguments.
	keeparelease_args := []string{
		"-f",
		"Changelog.md",
	}
	keeparelease_args = append(keeparelease_args, assets...)
	sh.RunV("keeparelease", keeparelease_args...)
}

// Release a new version.
func Release() {
	mg.Deps(Dist)
	mg.Deps(Publish)
}

func builder(platform, arch, tag, out, cwd string) {
	cmd := exec.Command(
		"go",
		"build",
		fmt.Sprintf("-ldflags"),
		fmt.Sprintf("-X github.com/rgreinho/tfe-cli/cmd.Version=%s", tag),
		fmt.Sprintf("-o"),
		fmt.Sprintf("%s-%s-%s-%s", out, tag, platform, arch),
	)
	cmd.Dir = cwd
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOOS=%s", platform))
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOARCH=%s", arch))
	fmt.Printf("Build command: %s\n", cmd.String())
	if stdoutStderr, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Error while running: %s\n", cmd.String())
		fmt.Printf("%s\n", stdoutStderr)
	}
}

func getTag() string {
	tag := os.Getenv("GITHUB_REF")
	if strings.HasPrefix(tag, "refs/tags/") {
		return strings.Replace(tag, "refs/tags/", "", 1)
	}

	// If the tag could not be retrieved from the environment, use Git.
	if tag == "" {
		tag, err := sh.Output("git", "describe")
		if err != nil {
			fmt.Printf("Cannot retrieve current git tag: %s.\n", err)
			os.Exit(1)
		}
		tag = tag
	}

	return tag
}

func getBuiltArtifacts() []os.FileInfo {
	currentDir, _ := os.Getwd()
	buildDir := filepath.Join(currentDir, BuildDir)
	files, err := ioutil.ReadDir(buildDir)
	if err != nil {
		fmt.Printf("Cannot find to package in %s. Forgot to build?\n", buildDir)
		os.Exit(1)
	}
	return files
}
