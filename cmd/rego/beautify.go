package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/aporeto-inc/regolithe/spec"
)

func beautifyAll(srcdir, dstdir string) error {

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

func beautifyOne(src string) (err error) {

	var s *spec.Specification

	if src != "" {
		fi, err1 := os.Stat(src)
		if err1 != nil {
			return fmt.Errorf("Unable to read file '%s': %s", src, err1)
		}

		if fi.IsDir() {
			return fmt.Errorf("Given read file '%s' is a directory", src)
		}

		if path.Ext(fi.Name()) != ".spec" {
			return fmt.Errorf("Given read file '%s' is not a spec file", src)
		}

		s, err = spec.LoadSpecification(src)
	} else {
		s, err = spec.LoadSpecificationFrom(os.Stdin)
	}

	if err != nil {
		return fmt.Errorf("Unable to read file '%s': %s", src, err)
	}

	s.Attributes = s.OriginalSortedAttributes()
	data, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
