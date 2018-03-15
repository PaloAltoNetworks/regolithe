package magetask

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/blang/semver"
	"github.com/magefile/mage/sh"
	"golang.org/x/sync/errgroup"
)

var projectName string

func init() {

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	projectName = path.Base(wd)
}

// SetProjectName sets the name of the project.
func SetProjectName(name string) {
	projectName = name
}

// WriteVersion creates the version file if needed.
func WriteVersion() error {

	return WriteVersionIn("")
}

// GetSemver gets the semantic version of the repository
func GetSemver(branch string) (sver string, err error) {

	versions, err := sh.Output("git", "tag", "--sort", "version:refname", "--merged", branch)
	if err != nil {
		return "", err
	}

	sver = "0.0.0"
	last, _ := semver.New(sver)
	for _, v := range strings.Split(versions, "\n") {
		curr, err := semver.New(v)
		if err != nil {
			continue
		}
		if last.Compare(*curr) < 0 {
			last = curr
			sver = v
		}
	}

	return
}

// WriteVersionIn creates the version file if needed in the given folder.
func WriteVersionIn(out string) error {

	projectBranch, err := sh.Output("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return err
	}

	projectVersion, err := GetSemver(projectBranch)
	if err != nil {
		return err
	}

	projectSha, err := sh.Output("git", "rev-parse", "HEAD")
	if err != nil {
		return err
	}

	if _, err := os.Stat("./Gopkg.toml"); err == nil {
		if err := makeVersionFromDep("", out, projectVersion, projectSha); err != nil {
			return err
		}
	}

	fmt.Println("complete: versions file generation")
	return nil
}

// Lint runs the linters.
func Lint() error {

	if err := run(
		nil,
		"gometalinter",
		"--exclude",
		"bindata.go",
		"--exclude",
		"vendor",
		"--vendor",
		"--disable-all",
		"--enable",
		"vet",
		"--enable",
		"vetshadow",
		"--enable",
		"golint",
		"--enable",
		"ineffassign",
		"--enable",
		"goconst",
		"--enable",
		"errcheck",
		"--enable",
		"varcheck",
		"--enable",
		"structcheck",
		"--enable",
		"gosimple",
		"--enable",
		"misspell",
		"--enable",
		"deadcode",
		"--enable",
		"staticcheck",
		"--deadline",
		"5m",
		"--tests",
		"./...",
	); err != nil {
		return err
	}

	fmt.Println("complete: linters")
	return nil
}

// Test runs unit tests.
func Test() error {
	return TestWith(true, true)
}

// Build builds the project for the current machine.
func Build() error {

	env := map[string]string{
		"CGO_ENABLED": "0",
	}

	if err := run(env, "go", "build", "-a", "-installsuffix", "cgo"); err != nil {
		return err
	}

	fmt.Println("complete: default build")
	return nil
}

// BuildLinux builds the project for Linux.
func BuildLinux() error {

	env := map[string]string{
		"CGO_ENABLED": "0",
		"GOOS":        "linux",
		"GOARCH":      "amd64",
	}

	if err := run(env, "go", "build", "-a", "-installsuffix", "cgo"); err != nil {
		return err
	}

	fmt.Println("complete: linux build")
	return nil
}

// BuildWindows builds the project for Windows.
func BuildWindows() error {

	env := map[string]string{
		"CGO_ENABLED": "0",
		"GOOS":        "windows",
		"GOARCH":      "amd64",
	}

	if err := run(env, "go", "build", "-a", "-installsuffix", "cgo"); err != nil {
		return err
	}

	fmt.Println("complete: windows build")
	return nil
}

// BuildDarwin builds the project for macOS.
func BuildDarwin() error {

	env := map[string]string{
		"CGO_ENABLED": "0",
		"GOOS":        "darwin",
		"GOARCH":      "amd64",
	}

	if err := run(env, "go", "build", "-a", "-installsuffix", "cgo"); err != nil {
		return err
	}

	fmt.Println("complete: darwin build")
	return nil
}

// BuildFor builds a a cli for the given platform.
func BuildFor(platform string, buildFunc func() error) error {

	if err := os.MkdirAll("./build/"+platform, 0755); err != nil {
		return err
	}

	if err := buildFunc(); err != nil {
		return err
	}

	return sh.Run("mv", projectName, "build/"+platform+"/"+projectName)
}

// Package packages the project for docker build.
func Package() error {

	return PackageFrom(projectName)
}

// PackageFrom packages the given binary for docker build
func PackageFrom(path string) error {

	if err := os.MkdirAll("docker/app", 0755); err != nil {
		return err
	}

	if err := run(nil, "cp", "-a", path, "docker/app"); err != nil {
		return err
	}

	fmt.Println("complete: docker packaging")
	return nil
}

// Container creates the docker container.
func Container() error {

	image := fmt.Sprintf("gcr.io/aporetodev/%s", projectName)

	if tag := os.Getenv("DOMINGO_DOCKER_TAG"); tag != "" {
		image = image + ":" + tag
	}

	if err := os.MkdirAll("docker/app", 0755); err != nil {
		return err
	}

	if err := os.Chdir("docker"); err != nil {
		return err
	}

	defer os.Chdir("..") // nolint

	out, err := sh.Output("docker", "build", "-t", image, ".")
	if err != nil {
		fmt.Println(out)
		return err
	}

	fmt.Println("complete: docker build", image)
	return nil
}

// PackageCACerts retrieves the package the CA for docker.
func PackageCACerts() error {

	if err := run(nil, "go", "get", "-u", "github.com/agl/extract-nss-root-certs"); err != nil {
		return err
	}

	if err := run(nil, "curl", "-s", "https://hg.mozilla.org/mozilla-central/raw-file/tip/security/nss/lib/ckfw/builtins/certdata.txt", "-o", "certdata.txt"); err != nil {
		return err
	}

	if err := os.MkdirAll("docker/app", 0755); err != nil {
		return err
	}

	out, err := sh.Output("extract-nss-root-certs")
	if err != nil {
		return nil
	}

	if err := ioutil.WriteFile("docker/app/ca-certificates.pem", []byte(out), 0644); err != nil {
		return err
	}

	if err := os.Remove("certdata.txt"); err != nil {
		return err
	}

	fmt.Println("complete: ca packaging")
	return nil
}

func run(env map[string]string, cmd string, args ...string) error {

	out, err := sh.OutputWith(env, cmd, args...)
	if out != "" {
		fmt.Println(out)
	}

	if err != nil {
		return fmt.Errorf("Unable to run command `%s %s`: %s", cmd, strings.Join(args, " "), err)
	}

	return nil
}

// collectCoverages collects all coverage reports.
func collectCoverages() error {

	f, err := os.OpenFile("coverage.txt", os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close() // nolint
	if err != nil {
		return err
	}

	cfiles, err := filepath.Glob("*.cover")
	if err != nil {
		return err
	}

	for _, cfile := range cfiles {

		cf, err := os.Open(cfile)
		if err != nil {
			return err
		}

		defer cf.Close() // nolint

		if _, err := io.Copy(f, cf); err != nil {
			return err
		}

		os.Remove(cfile) // nolint
	}

	return nil
}

// getTestArgs provides arguments based on test options.
func getTestArgs(race bool, cover bool, p string) []string {
	args := []string{"test"}
	if race {
		args = append(args, "-race")
	}
	if cover {
		args = append(args, "-cover", "-coverprofile=cov-"+path.Base(p)+".cover", "-covermode=atomic")
	}
	return append(args, p)
}

// TestWith runs unit tests without race.
func TestWith(race bool, cover bool) error {

	out, err := sh.Output("go", "list", "./...")
	if err != nil {
		return err
	}

	var g errgroup.Group

	packages := strings.Split(out, "\n")
	for _, p := range packages {
		g.Go(func(p string) func() error {
			return func() error {
				args := getTestArgs(race, cover, p)
				return run(nil, "go", args...)
			}
		}(p))
	}

	if err = g.Wait(); err != nil {
		return err
	}

	if cover {
		if err := collectCoverages(); err != nil {
			return err
		}
	}

	fmt.Println("complete: tests")
	return nil
}
