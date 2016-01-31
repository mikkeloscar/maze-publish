package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"

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
	AuthToken string `json:"auth"`
}

func main() {
	var workspace = plugin.Workspace{}
	var vargs = Publish{}

	plugin.Param("workspace", &workspace)
	plugin.Param("vargs", &vargs)
	plugin.MustParse()

	var pkgs []*BuiltPkg
	err := loadBuiltPkgs(path.Join(workspace.Path, "packages.built"), &pkgs)

	uploader := Uploader{
		client: NewClient(vargs.URL),
		repo:   vargs.Repo,
	}

	err = uploader.Do(pkgs)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
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
