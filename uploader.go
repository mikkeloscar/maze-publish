package main

import (
	"fmt"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
)

// Uploader is a wrapper around http client for uploading packages to a
// repository.
type Uploader struct {
	client    *client
	owner     string
	name      string
	sessionID string
}

func (u *Uploader) Do(pkgs []*BuiltPkg) error {
	if len(pkgs) == 0 {
		log.Println("No packages to upload")
		return nil
	}

	uResp, err := u.client.UploadStart(u.owner, u.name)
	if err != nil {
		return err
	}

	if uResp.SessionID == "" {
		return fmt.Errorf("invalid sessionID")
	}

	u.sessionID = uResp.SessionID

	err = u.uploadPkgs(pkgs)
	if err != nil {
		return err
	}

	err = u.client.UploadDone(u.owner, u.name, u.sessionID)

	return nil
}

func (u *Uploader) uploadPkgs(pkgs []*BuiltPkg) error {
	dl := make(chan error)

	for _, pkg := range pkgs {
		go u.uploadPkg(pkg, dl)
	}

	var errors []error

	for range pkgs {
		err := <-dl
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		msg := "errors while uploading packages\n"
		for _, err := range errors {
			msg += fmt.Sprintf("%s * %s\n", msg, err.Error())
		}
		return fmt.Errorf(msg)
	}

	return nil
}

func (u *Uploader) uploadPkg(pkg *BuiltPkg, ch chan<- error) {
	err := u.uploadFile(pkg.Package)
	ch <- err

	if err != nil {
		return
	}

	if pkg.Signature != "" {
		err = u.uploadFile(pkg.Signature)
		ch <- err
	}
}

func (u *Uploader) uploadFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	err = u.client.UploadFile(u.owner, u.name, path.Base(file), u.sessionID, f)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
