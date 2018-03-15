package magetask

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"

	"golang.org/x/tools/imports"
)

// writeVersionsFile writes the version file from the given templateData at the given path.
func writeVersionsFile(templateData versionTemplate, versionsFilePath string) error {

	t := template.Must(template.New("versions").Funcs(template.FuncMap{
		"cap":       strings.Title,
		"short":     path.Base,
		"hasprefix": strings.HasPrefix,
		"varname": func(v string) string {
			return strings.Title(path.Base(strings.Replace(v, "-", "", -1)))
		},
	}).Parse(versionFileTemplate))

	outFile := path.Join(versionsFilePath, "versions.go")
	_, err := os.Stat(versionsFilePath)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(versionsFilePath, 0744); err != nil {
			panic(err)
		}
	}

	_, err = os.Stat(outFile)
	if err == nil {
		if err = os.Remove(outFile); err != nil {
			panic(err)
		}
	}

	f, err := os.Create(outFile)
	if err != nil {
		panic(err)
	}
	defer f.Close() // nolint: errcheck

	buffer := &bytes.Buffer{}
	if err = t.Execute(buffer, templateData); err != nil {
		return err
	}

	data, err := imports.Process(".", buffer.Bytes(), &imports.Options{
		TabWidth:  8,
		TabIndent: true,
		Comments:  true,
		Fragment:  true,
	})
	if err != nil {
		fmt.Println(string(buffer.Bytes()))
		panic(err)
	}

	if _, err = f.Write(data); err != nil {
		return err
	}

	return nil
}
