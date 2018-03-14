package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/aporeto-inc/regolithe/spec"
)

func beautify(reader io.Reader) error {

	s := spec.NewSpecification()

	if err := s.Read(reader); err != nil {
		return fmt.Errorf("Unable to load specs: %s", err)
	}

	buf := &bytes.Buffer{}

	if err := s.Write(buf); err != nil {
		return err
	}

	fmt.Println(buf.String())

	return nil
}
