package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/aporeto-inc/regolithe/spec"
)

func beautify(srcdir, dstdir string) error {

	if srcdir == "" {
		return fmt.Errorf("You must provide --src")
	}

	if dstdir == "" {
		dstdir = srcdir
	}

	if err := os.MkdirAll(dstdir, 0755); err != nil {
		return fmt.Errorf("Unable to create dest dir: %s", err)
	}

	filesInfo, err := ioutil.ReadDir(srcdir)
	if err != nil {
		return fmt.Errorf("Unable to read dir: %s", err)
	}

	for _, info := range filesInfo {

		if path.Ext(info.Name()) != ".spec" {
			continue
		}

		s, err := spec.LoadSpecification(path.Join(srcdir, info.Name()))
		if err != nil {
			return fmt.Errorf("Unable to read file '%s': %s", info.Name(), err)
		}

		if err = s.Write(dstdir); err != nil {
			return fmt.Errorf("Unable to write file '%s': %s", info.Name(), err)
		}
	}

	return nil
}
