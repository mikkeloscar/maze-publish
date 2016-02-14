package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/drone/drone-plugin-go/plugin"
)

// BuiltPkg defines a built package and optional signature file.
type BuiltPkg struct {
	Package   string `json:"package"`
	Signature string `json:"signature"`
}

// Publish defines the args passed from the config file.
type Publish struct {
	URL       string `json:"url"`
	Repo      string `json:"repo"`
	AuthToken string `json:"auth_token"`
}

func main() {
	var workspace = plugin.Workspace{}
	var vargs = Publish{}

	plugin.Param("workspace", &workspace)
	plugin.Param("vargs", &vargs)
	plugin.MustParse()

	var pkgs []*BuiltPkg
	err := loadBuiltPkgs(path.Join(workspace.Path, "drone_pkgbuild", "packages.built"), &pkgs)

	owner, name, err := splitOwnerName(vargs.Repo)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	uploader := Uploader{
		client: NewClientToken(vargs.URL, vargs.AuthToken),
		owner:  owner,
		name:   name,
	}

	err = uploader.Do(pkgs)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func splitOwnerName(repo string) (string, string, error) {
	split := strings.Split(repo, "/")
	if len(split) != 2 {
		return "", "", fmt.Errorf("invalid repo format: %s", repo)
	}
	return split[0], split[1], nil
}

type formatter struct{}

func (f *formatter) Format(entry *log.Entry) ([]byte, error) {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "[%s] %s\n", entry.Level.String(), entry.Message)
	return buf.Bytes(), nil
}

func loadBuiltPkgs(file string, out interface{}) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	return dec.Decode(out)
}
