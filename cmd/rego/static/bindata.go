// Code generated by go-bindata.
// sources:
// templates/json-schema-global.gotpl
// templates/json-schema.gotpl
// templates/spec-md.gotpl
// templates/toc-md.gotpl
// specset/.gitignore
// specset/Gopkg.toml
// specset/specs/.regolithe-gen-cmd
// specset/specs/@identifiable.abs
// specset/specs/_api.info
// specset/specs/_type.mapping
// specset/specs/object.spec
// specset/specs/regolithe.ini
// specset/specs/root.spec
// specset/specs/type_mapping.ini
// DO NOT EDIT!

package static

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _templatesJsonSchemaGlobalGotpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x51\xcd\x6a\xf3\x30\x10\xbc\xfb\x29\x16\xe1\xa3\xa2\x07\x30\x7c\x87\xef\xde\xf6\x50\x43\xef\x8a\x3c\x2e\x0a\x8e\x25\xb4\x1b\x68\x10\x7a\xf7\x62\xd7\x0e\x76\x42\x69\xa1\x3a\x89\x99\x9d\x99\xfd\xc9\x15\x11\x91\x12\x2f\x03\x54\x43\x2a\x67\x32\x2f\xf6\x0c\x2a\x45\xe9\x85\xbb\xc6\x99\x0a\xc7\x13\x9c\xac\x68\x4c\x21\x22\x89\x07\xab\x86\xbe\x5c\xa6\x97\xf3\x81\x92\x1d\xdf\x41\xb5\xef\x3e\xda\x08\xa7\xa9\xe6\x08\xe7\x7b\xef\xac\xf8\x30\x52\xf3\x8f\x4c\x0b\x31\xed\x16\x65\x3a\x94\xb2\x73\xf1\xfd\xcd\x62\xe2\xf4\x84\x61\xec\x68\x57\x46\xf5\x60\x05\x2c\x6f\x48\xbc\x78\xef\xd3\xcc\xd3\xcc\xff\x17\x49\xfe\x78\x11\xf0\x5a\xb9\x8d\x9b\xa6\xbe\x93\x3d\x87\x0e\x83\x79\x05\xcb\xba\x8d\xed\x94\xb3\xa8\x4e\xe8\x97\x95\xfd\x24\x36\x27\x0e\xa3\xba\xc9\x8b\xfe\x55\x74\xb8\x24\x87\x6f\xe3\xd7\xb3\xd8\x94\xec\x55\xe9\x3d\xe9\x05\x67\x7e\xd0\xfc\xb1\xed\xb9\xf5\xea\xf1\x77\x77\x99\x52\x95\xea\x33\x00\x00\xff\xff\xef\xc9\xe4\x09\x57\x02\x00\x00")

func templatesJsonSchemaGlobalGotplBytes() ([]byte, error) {
	return bindataRead(
		_templatesJsonSchemaGlobalGotpl,
		"templates/json-schema-global.gotpl",
	)
}

func templatesJsonSchemaGlobalGotpl() (*asset, error) {
	bytes, err := templatesJsonSchemaGlobalGotplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/json-schema-global.gotpl", size: 599, mode: os.FileMode(420), modTime: time.Unix(1549567471, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesJsonSchemaGotpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xbc\x57\x4f\x6f\xe3\xb8\x0f\xbd\xe7\x53\x10\x42\x0e\xbf\xdf\x20\xc9\xdc\x0b\xe4\x30\x3b\xdb\x3d\x2c\xfa\x67\xd0\x2e\xe6\xb2\xe8\x81\x89\x99\x44\x53\x5b\xf6\x48\xf4\x34\x41\xd6\xdf\x7d\x21\x59\xb2\x65\x5b\x69\xa7\x8b\xc5\xde\x5c\xea\x3d\xf2\x91\xa2\x18\xf6\x7c\x5e\x42\x46\x3b\xa9\x08\x44\x85\xda\xd0\x92\x4f\x15\x09\x58\x36\xcd\xcc\x9e\xc9\x1d\xd0\x77\x58\x81\x60\x59\x78\xb3\x70\x88\x2b\x10\x86\xb5\x54\x7b\xb1\x98\x89\x5d\xa9\x0b\x64\x6b\xcb\x90\x69\xe9\xb0\x8e\x4e\xb9\xa1\xde\x87\x27\x24\xbd\x24\xe0\x9b\xb2\xcc\x09\xd5\x08\x1f\xac\x09\x82\x54\x4c\x7b\xd2\x23\x42\xb0\x26\x08\xbb\xbc\x44\x1e\xc1\x55\x5d\x6c\xd2\xe8\x72\xf3\x8d\xb6\xdc\xd7\xc6\x9d\x3a\xee\x5c\xd3\xce\x72\xcf\x67\x58\x41\xd3\xac\xbe\x99\x32\x08\x54\x59\x8f\xf7\xdf\xb3\xf3\x0c\x00\xe0\x7c\x86\x79\x8e\x4c\x86\xbf\x92\x36\xb2\x54\x70\xb5\x86\xd5\x63\x45\xdb\xd5\x8d\x33\x7f\x62\xd6\x72\x53\x33\x99\x00\xb0\x6c\x4b\x15\x2c\x39\xa7\x10\xd1\x51\x6e\xcb\x8c\xf2\xd5\xb5\x62\xc9\xa7\x3b\x2c\x08\x9a\x46\x2c\x3c\xd8\x67\xe6\xf5\x7b\x6b\xa5\xcb\x8a\x34\x4b\x32\xe2\x0a\x82\xa2\x25\x68\x54\x7b\x82\xb9\xcc\x8e\x36\xfc\x02\xe6\xc8\xac\x7b\x65\xd7\xc7\xaa\x34\x94\xf5\xd2\xc6\x39\x04\x89\xbe\x79\x82\x23\x6b\x5f\x84\x22\x84\x2c\x6c\x05\xac\xfb\xd5\xe7\x52\xfd\x20\xcd\x94\x05\xe5\x41\x91\x40\x75\xba\xb7\xa5\xfd\xd3\xfd\x19\x39\x56\x25\xc3\xea\x81\xbe\xd7\x52\x53\xe7\xd1\x9d\x43\x74\x95\x79\x2e\xa0\x59\x0c\xb8\x91\x00\x67\x1a\x1c\x7e\xfc\x00\xeb\xf5\x7a\x0d\x7f\x9c\x2a\x6a\xbf\x3e\x7c\xec\x52\x0a\xa0\xb9\xf5\x6f\x4b\xb2\x6d\x65\x3b\x70\x9b\x88\xfb\x1c\xe2\x5d\x15\xdc\xe1\xa7\x3c\x2f\x5f\x28\xfb\x7c\x28\xe5\x96\x4c\x2c\x42\x90\xaa\x8b\x41\x96\xd3\xeb\x68\x69\x0b\x98\x6f\xdd\x87\x8d\x9f\x74\x1b\x47\x1f\x5d\x44\x0b\x71\x57\x61\x85\x85\x7e\x8c\xd1\xee\x4e\x7c\x84\xa6\x11\x13\x4f\x63\xce\xd3\xb0\xb6\xfd\x7b\x69\x8b\x24\xe6\x74\x64\xd2\x0a\xed\x3d\x8c\x75\x81\x7d\xfa\xd5\x6f\x52\x1b\xbe\xa1\x1f\x94\xff\xa2\x71\xfb\x4c\x6c\xc6\x4d\xe1\x8a\x3a\xba\x83\x44\xa4\x5c\x1a\x1e\x44\xe9\xfa\x00\xb5\xc6\x93\x18\x77\x50\x1b\xe5\xb1\xde\x78\xff\x7d\x9b\x08\xc9\x54\xf4\xcf\x22\x52\xcc\x54\x54\xb6\xdf\x87\xc3\x72\xec\xa9\x63\x35\xaf\x76\x5e\x3a\x8d\x02\xab\x74\x16\xfe\xf9\xbe\x9e\xc6\x72\x90\x07\x66\x99\x64\x59\x2a\xcc\xbf\x4c\x5f\xfb\xbf\x96\x96\x1f\x83\x6f\xc7\x65\x5d\xd3\x3f\xa8\x88\x1d\xae\x43\xd8\x7b\x05\x77\x8e\x7f\xca\x0b\xa7\xd8\xad\xd2\xe4\xb0\xb8\xe7\x03\x69\xe8\xe7\xe9\xc5\xc1\x31\x9d\x04\xa8\xcd\xa8\x74\x15\xb2\x7d\x30\x7e\xb6\xfb\x09\xf3\x40\x7b\x3a\x56\x29\x72\x6b\x8a\xe6\xa0\x98\xc8\xbe\xa8\xe2\x16\x8f\x37\xa4\xf6\x7c\x18\x49\x28\x82\xdd\xb6\xca\x79\x02\x4e\x54\xe6\x72\x08\xa9\xd2\x21\x82\x7d\x10\xa2\x03\xbf\x2b\x04\x1e\xbf\x62\x5e\x8f\x5b\xb0\xc0\xa3\x2c\xdc\x50\x8d\x53\x68\x91\xef\xcc\x20\xe9\x5e\xaa\x89\xfb\x80\x7c\xd3\x3d\x08\x31\x7d\x4e\x4f\xb3\xd9\xb8\xb1\x6e\x89\x31\x43\xc6\x2b\xb8\x2b\x19\x10\x7e\x7f\xbc\xbf\x03\xb3\x3d\x50\x81\xf0\x4c\xa7\x97\x52\x67\xb0\xa9\x19\x9e\xa9\x62\xd8\x95\x1a\xa4\x6a\x37\x31\xfb\x43\x3c\xe9\xc2\x41\x5a\x5d\xc7\x74\x69\xd9\x77\xd6\xda\xfc\x4b\x0d\x52\xe2\xa8\xb5\x21\x03\xa8\xc0\x4d\x54\x40\x06\xd1\x91\x3a\x45\xc8\xc0\x07\x82\x0a\x35\x29\x1e\x0b\x88\x8b\x31\x12\x84\xd9\xbd\xca\x4f\x43\x41\xad\x2d\x1e\x1d\x17\x9c\xd8\x75\xe0\x7f\xd2\xdc\xc9\xdc\xfb\xfb\x95\x76\x58\xe7\xec\x6e\xe4\xff\xb1\xd3\x2c\x3a\x88\xaf\x2f\x26\xc0\x5f\x60\x57\xb8\x47\xb7\x9a\xca\xdd\x09\x2e\x66\x10\x5d\xe5\x10\x12\x8c\x3e\xec\x85\x89\xb8\xc3\xdc\x90\x87\x44\xe5\x1f\x6e\x3b\x2f\x92\x0f\xef\x5d\xbf\x02\x77\xbe\xb3\x3f\xaf\x76\x59\x70\x57\x3a\x3e\xf6\xeb\x05\x06\x7f\x6e\xd1\x1b\x37\x70\xb8\x23\x87\x18\x74\x4e\x6a\xd1\xb0\x37\xe1\xa3\xbe\xb6\x68\x44\xe2\xd6\x6d\x19\xd2\xbb\x48\x1f\x78\xb2\x24\xbe\xf9\x88\xdf\xb2\x3d\x2d\xfe\x8b\x37\xd7\x77\x88\x98\x17\x76\x4b\xb7\xfa\x7f\x76\x77\x6f\x19\x0f\x64\xf8\x02\x2b\x1c\xa5\x38\x65\xad\xb7\x74\x99\xd7\x1d\x4f\xb8\x5f\x70\xfb\x8c\xfb\x14\xcd\x9f\xc4\x0c\xbb\xe1\xa9\xcc\x4c\xda\xb6\x5f\x5d\xaf\x1d\x62\x01\x1e\xda\xff\x37\xe1\x13\x6f\x1d\x24\x67\x70\xc7\xbe\xd8\x4c\xae\x4b\xbc\xe3\xd7\x9a\xe2\x69\xd6\xcc\xfe\x0e\x00\x00\xff\xff\xda\x1d\x70\x5c\xee\x0e\x00\x00")

func templatesJsonSchemaGotplBytes() ([]byte, error) {
	return bindataRead(
		_templatesJsonSchemaGotpl,
		"templates/json-schema.gotpl",
	)
}

func templatesJsonSchemaGotpl() (*asset, error) {
	bytes, err := templatesJsonSchemaGotplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/json-schema.gotpl", size: 3822, mode: os.FileMode(420), modTime: time.Unix(1549567555, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesSpecMdGotpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x52\x4b\x4e\xc3\x30\x10\xdd\xe7\x14\xa3\x86\x05\x2c\xea\x03\x54\x02\x09\xa9\xdd\x01\x0b\xa8\x58\xc7\x75\x06\x6a\x94\xc6\xc6\x33\x95\x5a\x45\xbe\x3b\xf2\x27\x4e\x52\xb1\xb3\x47\xef\xe3\xf7\x3c\x75\x5d\xc3\x30\x80\xf8\xb0\xa8\xc4\xab\x69\xb1\x13\xbb\x9e\x35\x5f\xdf\xe4\x09\xc1\xfb\x6a\x18\xe0\xae\x93\x8c\xc4\x9f\xe8\x48\x9b\x1e\x36\x8f\x19\xfe\x12\xc7\xcf\xcc\x4e\x1f\xce\x8c\x34\x02\x12\x8b\x9d\x3e\x91\x95\x0a\x17\xe2\x5b\x24\xe5\xb4\xe5\x8c\x8b\xf2\x78\x91\x27\xdb\x61\x10\x1e\x8f\x91\x72\x6b\x9c\x74\xf5\xd7\xc4\xf0\xbe\xaa\x43\x82\x5d\xba\x57\x55\xd3\x34\x3f\x64\xfa\x85\xac\xf7\x61\x1c\x46\xd8\xb7\xc5\xd4\x58\x0a\x86\xc6\xa2\x93\xe1\x35\x94\x3d\xc5\x3b\x76\x69\x70\xd4\x36\x0c\x91\x67\xc6\xc6\xd2\xe8\x59\x70\x93\x5e\xc2\xcd\x5d\x64\x29\x67\xaa\x6d\x77\xb1\x86\xb0\x9d\x7a\xfb\x37\xe7\x3a\xfa\xcd\xf8\x41\x32\x1a\x4f\xc4\xe8\xe1\x64\xff\x8d\x0b\xe4\x3a\x3f\xb1\x86\x26\x7c\x6d\xfe\x49\xb8\x4f\x19\xf0\x17\xc4\xfe\x6a\x11\x56\x78\x61\x74\xbd\xec\x56\xe0\x7d\x00\xc6\xa9\xf7\x9b\xb8\x0f\xe7\x43\xbe\x86\x44\x1d\xe1\x12\x53\x62\x3e\x34\xb9\x1a\xb1\x45\xeb\x50\x49\xc6\x98\xfe\x09\xf6\x47\x4d\x50\x5e\x05\x9a\xa0\x2d\x88\x59\x4d\xcb\x4d\xb9\x59\x8f\x61\x00\x75\x94\x4e\x2a\x46\xa7\x89\xb5\x22\x10\x63\x3d\x85\x5f\x8e\x7f\x01\x00\x00\xff\xff\xd6\x3e\xd8\x68\xcc\x02\x00\x00")

func templatesSpecMdGotplBytes() ([]byte, error) {
	return bindataRead(
		_templatesSpecMdGotpl,
		"templates/spec-md.gotpl",
	)
}

func templatesSpecMdGotpl() (*asset, error) {
	bytes, err := templatesSpecMdGotplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/spec-md.gotpl", size: 716, mode: os.FileMode(420), modTime: time.Unix(1549403378, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesTocMdGotpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x52\x56\x56\x48\xa8\xae\x56\xd0\x73\x2f\xca\x2f\x2d\xf0\x4b\xcc\x4d\x55\xa8\xad\x4d\xe0\xe2\xaa\xae\x56\x28\xc9\x4f\x56\xd0\xd0\x0b\x4e\x2d\xd1\x0b\x2e\x48\x4d\xce\x4c\xcb\x4c\x4e\x2c\xc9\xcc\xcf\x03\x2b\x44\x52\xaf\xa9\xa0\x5b\x5b\xcb\x05\x08\x00\x00\xff\xff\x58\xb9\x92\xee\x47\x00\x00\x00")

func templatesTocMdGotplBytes() ([]byte, error) {
	return bindataRead(
		_templatesTocMdGotpl,
		"templates/toc-md.gotpl",
	)
}

func templatesTocMdGotpl() (*asset, error) {
	bytes, err := templatesTocMdGotplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/toc-md.gotpl", size: 71, mode: os.FileMode(420), modTime: time.Unix(1549405303, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _specsetGitignore = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x24\xca\x41\x0a\x84\x30\x0c\x05\xd0\xfd\x3f\x4a\x61\x72\xa1\x61\x16\x43\xf2\x0d\xc5\x9a\x48\x2d\x45\x6f\x2f\xe2\xe6\xad\x5e\x91\xfd\xfa\x6a\xda\x0f\x45\xe8\xfe\xfa\xa9\xb1\x24\x34\x8d\xce\xc0\x53\x14\xde\xc9\x51\xc3\xb7\x34\xb6\x03\x93\x61\xd9\x51\xa4\xa5\xae\xd0\x9c\xec\x7f\xa7\x8c\x73\xe0\x0e\x00\x00\xff\xff\x66\xf6\x5a\x96\x53\x00\x00\x00")

func specsetGitignoreBytes() ([]byte, error) {
	return bindataRead(
		_specsetGitignore,
		"specset/.gitignore",
	)
}

func specsetGitignore() (*asset, error) {
	bytes, err := specsetGitignoreBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "specset/.gitignore", size: 83, mode: os.FileMode(420), modTime: time.Unix(1528494787, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _specsetGopkgToml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8a\x2e\x28\x2a\xcd\x4b\x8d\xe5\x52\x50\x48\xcf\xd7\x2d\x49\x2d\x2e\x29\x56\xb0\x55\x28\x29\x2a\x4d\xe5\x52\x50\x28\xcd\x2b\x2d\x4e\x4d\xd1\x2d\x48\x4c\xce\x4e\x4c\x4f\x85\x4b\x00\x02\x00\x00\xff\xff\x34\x36\xb3\x89\x33\x00\x00\x00")

func specsetGopkgTomlBytes() ([]byte, error) {
	return bindataRead(
		_specsetGopkgToml,
		"specset/Gopkg.toml",
	)
}

func specsetGopkgToml() (*asset, error) {
	bytes, err := specsetGopkgTomlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "specset/Gopkg.toml", size: 51, mode: os.FileMode(420), modTime: time.Unix(1528494787, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _specsetSpecsRegolitheGenCmd = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x4c\x8e\x41\x4e\xc5\x30\x0c\x44\xf7\x3e\xc5\xf0\xf9\x0b\x40\x4a\x22\x0e\x00\x77\x49\x63\x37\x8d\x48\xe2\x2a\x0d\x88\x45\x0f\x8f\x5a\x5a\xca\xc2\x0b\x6b\xde\xb3\xe7\xf1\xc1\x0d\xa9\xba\xc1\x2f\x13\x95\x0f\x4e\x0d\x66\x86\xb5\x8e\x35\x10\xfb\xee\xdf\xee\x4f\x4d\xa2\x82\x35\xc0\x30\x2c\xd6\x15\xf2\x9d\x3a\x5e\x9f\x49\xc2\xa4\xb8\xdd\x37\xec\x86\xf7\xc3\xda\xe6\xb3\x48\xed\xbe\x27\xad\xb6\x30\x05\x86\xfd\xe7\x91\x64\x89\x52\x31\x6a\x66\x69\xdb\xd1\x65\x96\xb0\xc0\x28\x82\xf2\x1e\x5d\xec\xd5\x29\x6a\xf6\x35\x52\x2b\x30\x6d\x84\x75\xbf\xbb\x7b\xb1\x51\xa9\x7c\x9d\xaa\x93\x2c\xfb\xf3\xbc\x27\x7f\xdc\xe9\x1d\x18\xfd\x04\x00\x00\xff\xff\x25\x5b\x5e\xc5\xf7\x00\x00\x00")

func specsetSpecsRegolitheGenCmdBytes() ([]byte, error) {
	return bindataRead(
		_specsetSpecsRegolitheGenCmd,
		"specset/specs/.regolithe-gen-cmd",
	)
}

func specsetSpecsRegolitheGenCmd() (*asset, error) {
	bytes, err := specsetSpecsRegolitheGenCmdBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "specset/specs/.regolithe-gen-cmd", size: 247, mode: os.FileMode(420), modTime: time.Unix(1528494787, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _specsetSpecsIdentifiableAbs = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x54\xce\xb1\x0a\xc2\x40\x0c\xc6\xf1\xbd\x4f\xf1\x81\xb3\x82\xeb\x6d\x82\x8b\x4f\x21\xd7\xde\xd7\x1a\xa9\x97\x92\x4b\xc5\xbe\xbd\xf4\x06\x6d\x97\x10\x7e\xfc\x03\x39\xe0\xe2\x6e\xd2\xce\xce\xd2\xc4\xdf\x1a\x1a\xe0\x7d\x5e\xe7\x11\x39\xbe\x18\x70\xbb\x36\x00\x90\x58\x3a\x93\xc9\x45\xf3\x6a\x90\x02\x7f\x10\x92\x98\x5d\x7a\xa1\x41\xfb\x2a\xda\x3e\xd9\xf9\xa9\x1e\xf9\x32\x31\xa0\xb8\x49\x1e\x2a\xf0\x33\x69\x61\x0a\x70\x9b\x59\xa5\xb8\xda\x0e\x8c\x31\xdd\x35\x8f\xcb\xc6\xe2\xec\x3a\x30\xd3\xa2\xef\xda\x5e\x46\xa7\xc5\x76\xe4\x06\xff\x2f\x6d\x50\x2d\xed\xc2\x6f\x00\x00\x00\xff\xff\x62\x4c\x45\x33\xff\x00\x00\x00")

func specsetSpecsIdentifiableAbsBytes() ([]byte, error) {
	return bindataRead(
		_specsetSpecsIdentifiableAbs,
		"specset/specs/@identifiable.abs",
	)
}

func specsetSpecsIdentifiableAbs() (*asset, error) {
	bytes, err := specsetSpecsIdentifiableAbsBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "specset/specs/@identifiable.abs", size: 255, mode: os.FileMode(420), modTime: time.Unix(1533330802, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _specsetSpecs_apiInfo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x2a\x28\x4a\x4d\xcb\xac\xb0\x52\x48\x2c\xc8\xe4\x2a\xca\xcf\x2f\xb1\x52\x00\x91\x5c\x65\xa9\x45\xc5\x99\xf9\x79\x56\x0a\x86\x5c\x80\x00\x00\x00\xff\xff\x0c\x97\x42\xd8\x22\x00\x00\x00")

func specsetSpecs_apiInfoBytes() ([]byte, error) {
	return bindataRead(
		_specsetSpecs_apiInfo,
		"specset/specs/_api.info",
	)
}

func specsetSpecs_apiInfo() (*asset, error) {
	bytes, err := specsetSpecs_apiInfoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "specset/specs/_api.info", size: 34, mode: os.FileMode(420), modTime: time.Unix(1528494787, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _specsetSpecs_typeMapping = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x4a\xcd\x49\xcd\x4d\xcd\x2b\x49\xcc\xb1\x52\xa8\xae\xe5\x02\x04\x00\x00\xff\xff\x14\xdf\xfc\xd6\x0e\x00\x00\x00")

func specsetSpecs_typeMappingBytes() ([]byte, error) {
	return bindataRead(
		_specsetSpecs_typeMapping,
		"specset/specs/_type.mapping",
	)
}

func specsetSpecs_typeMapping() (*asset, error) {
	bytes, err := specsetSpecs_typeMappingBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "specset/specs/_type.mapping", size: 14, mode: os.FileMode(420), modTime: time.Unix(1528494787, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _specsetSpecsObjectSpec = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x90\x31\x6b\xc3\x30\x10\x85\x77\xfd\x8a\x07\x19\x32\xb5\xd0\x55\x53\x0b\x85\x4e\xa5\x4b\x3a\x17\xc5\x7a\xb6\xd5\xda\x92\xd1\x9d\x4a\xf2\xef\x8b\xec\x10\x62\x37\xcb\x21\xbe\xfb\xf4\xee\xa4\x1d\xde\x93\xe7\x60\xc6\x5a\xad\x01\x32\x45\xbf\xa2\x1b\x69\x91\x8e\xdf\x6c\x74\x61\xa9\xe4\x86\x2b\x2e\x06\x60\xd4\xa0\xe7\x0b\xfe\xb8\xe2\xc9\x35\x3f\xae\xa3\x85\x67\xeb\xca\x50\x23\x3c\xa5\xc9\x61\xd2\x90\xa2\xc5\xa1\x0f\x82\x20\xc8\x2e\xfa\x34\x5e\xf2\x1e\x0d\xd0\x51\xeb\x0e\x1b\xfd\x8d\x2a\xd0\x9e\x37\x62\x99\xbc\x53\xde\x71\x3f\xe7\xc6\x56\xf7\x1c\x78\x57\x7f\x9d\x1b\x5b\x9d\x27\x65\xf4\x52\xfd\x07\xec\x9f\x83\xaf\xef\x6c\x83\x3b\x0e\xdc\x1b\xb3\xc3\x8b\x6a\x0e\xc7\xa2\x14\xe3\xae\xc7\x6a\xff\x3e\x2d\x77\x96\x0f\xa9\xf5\xff\xc8\x43\xcf\xb9\x83\xd4\x6e\xc6\x02\x7a\x9e\x68\x21\x9a\x43\xec\x66\xc0\xd3\x94\x84\xde\x42\x73\x59\xb2\x44\x53\x5e\x81\x36\x0c\xca\x5c\x77\xbb\x81\x29\xfb\x15\xfb\x0b\x00\x00\xff\xff\x4a\xf2\xda\x7e\xe7\x01\x00\x00")

func specsetSpecsObjectSpecBytes() ([]byte, error) {
	return bindataRead(
		_specsetSpecsObjectSpec,
		"specset/specs/object.spec",
	)
}

func specsetSpecsObjectSpec() (*asset, error) {
	bytes, err := specsetSpecsObjectSpecBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "specset/specs/object.spec", size: 487, mode: os.FileMode(420), modTime: time.Unix(1533220356, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _specsetSpecsRegolitheIni = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x54\xcc\x41\x0a\x02\x31\x0c\x05\xd0\x7d\x4e\x91\x13\x54\x3d\x40\xc1\x7b\x0c\x83\x84\x1a\x67\x0a\x93\xa6\x64\x52\x4b\x6f\xef\xa2\x2e\x74\xf7\xe1\xff\xff\x16\xe3\x4d\x8f\xec\x3b\xaf\x50\x4d\x9f\x2d\xf9\xa3\x90\x30\x46\x94\x41\x35\x43\xd2\x3a\x2c\x6f\xbb\x63\xc4\xa1\xcd\x92\x4a\xa5\x32\x00\x16\x37\x2a\xe7\x4b\x4d\xd8\x56\xf8\xfb\x34\x3b\x30\x62\xef\x3d\xc8\xf8\xee\x43\x61\xbf\xcc\x96\x9a\xef\x6a\x93\x03\x16\xca\xc7\xcc\xf7\x1f\x3e\x24\x15\x78\xb3\x9d\x59\x0b\x46\xbc\x85\x2b\x7c\x02\x00\x00\xff\xff\x0e\x59\x87\x8e\xaa\x00\x00\x00")

func specsetSpecsRegolitheIniBytes() ([]byte, error) {
	return bindataRead(
		_specsetSpecsRegolitheIni,
		"specset/specs/regolithe.ini",
	)
}

func specsetSpecsRegolitheIni() (*asset, error) {
	bytes, err := specsetSpecsRegolitheIniBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "specset/specs/regolithe.ini", size: 170, mode: os.FileMode(420), modTime: time.Unix(1533330787, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _specsetSpecsRootSpec = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x8e\x31\x8e\xc2\x30\x10\x45\x7b\x9f\x62\xa4\xd4\xbb\x07\x70\xbb\xf5\x36\xb9\x00\x32\xce\x27\x18\x12\x4f\x34\x33\x01\x71\x7b\x64\x1c\x20\x20\x68\x2c\xeb\xcd\xff\x7a\xbf\xa1\x7f\xee\x30\xb8\xb1\xbc\xde\x11\x09\xd4\x36\x39\x8c\xf0\x24\xcc\x56\x09\xcf\x12\xf1\x4a\x91\x2d\xd9\x65\x61\x6d\x65\x53\x88\xc7\xd0\x3f\x33\x1d\x34\x4a\x9a\x2c\x71\xae\x8c\x78\x7b\x40\xb4\x5f\x47\xd4\xc3\x8a\xee\x2d\xd4\xc3\x94\x6c\x8f\x55\xb0\xf4\x3c\x99\xcc\x70\xae\xa1\x16\x43\x28\x51\x75\x72\xff\x79\xf7\xb3\x5e\x5d\x9b\x5f\x0d\x2d\x4c\x12\x4e\xa8\x9a\x21\xa9\x11\xef\x96\x92\x16\x5f\x14\x04\xc3\x87\xe6\xdf\xed\xa0\x14\x28\xe3\xfc\x18\x78\x0d\x00\x00\xff\xff\x22\x7b\xf8\x38\x41\x01\x00\x00")

func specsetSpecsRootSpecBytes() ([]byte, error) {
	return bindataRead(
		_specsetSpecsRootSpec,
		"specset/specs/root.spec",
	)
}

func specsetSpecsRootSpec() (*asset, error) {
	bytes, err := specsetSpecsRootSpecBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "specset/specs/root.spec", size: 321, mode: os.FileMode(420), modTime: time.Unix(1533330789, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _specsetSpecsType_mappingIni = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xe2\x02\x04\x00\x00\xff\xff\x93\x06\xd7\x32\x01\x00\x00\x00")

func specsetSpecsType_mappingIniBytes() ([]byte, error) {
	return bindataRead(
		_specsetSpecsType_mappingIni,
		"specset/specs/type_mapping.ini",
	)
}

func specsetSpecsType_mappingIni() (*asset, error) {
	bytes, err := specsetSpecsType_mappingIniBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "specset/specs/type_mapping.ini", size: 1, mode: os.FileMode(420), modTime: time.Unix(1533220356, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"templates/json-schema-global.gotpl": templatesJsonSchemaGlobalGotpl,
	"templates/json-schema.gotpl": templatesJsonSchemaGotpl,
	"templates/spec-md.gotpl": templatesSpecMdGotpl,
	"templates/toc-md.gotpl": templatesTocMdGotpl,
	"specset/.gitignore": specsetGitignore,
	"specset/Gopkg.toml": specsetGopkgToml,
	"specset/specs/.regolithe-gen-cmd": specsetSpecsRegolitheGenCmd,
	"specset/specs/@identifiable.abs": specsetSpecsIdentifiableAbs,
	"specset/specs/_api.info": specsetSpecs_apiInfo,
	"specset/specs/_type.mapping": specsetSpecs_typeMapping,
	"specset/specs/object.spec": specsetSpecsObjectSpec,
	"specset/specs/regolithe.ini": specsetSpecsRegolitheIni,
	"specset/specs/root.spec": specsetSpecsRootSpec,
	"specset/specs/type_mapping.ini": specsetSpecsType_mappingIni,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"specset": &bintree{nil, map[string]*bintree{
		".gitignore": &bintree{specsetGitignore, map[string]*bintree{}},
		"Gopkg.toml": &bintree{specsetGopkgToml, map[string]*bintree{}},
		"specs": &bintree{nil, map[string]*bintree{
			".regolithe-gen-cmd": &bintree{specsetSpecsRegolitheGenCmd, map[string]*bintree{}},
			"@identifiable.abs": &bintree{specsetSpecsIdentifiableAbs, map[string]*bintree{}},
			"_api.info": &bintree{specsetSpecs_apiInfo, map[string]*bintree{}},
			"_type.mapping": &bintree{specsetSpecs_typeMapping, map[string]*bintree{}},
			"object.spec": &bintree{specsetSpecsObjectSpec, map[string]*bintree{}},
			"regolithe.ini": &bintree{specsetSpecsRegolitheIni, map[string]*bintree{}},
			"root.spec": &bintree{specsetSpecsRootSpec, map[string]*bintree{}},
			"type_mapping.ini": &bintree{specsetSpecsType_mappingIni, map[string]*bintree{}},
		}},
	}},
	"templates": &bintree{nil, map[string]*bintree{
		"json-schema-global.gotpl": &bintree{templatesJsonSchemaGlobalGotpl, map[string]*bintree{}},
		"json-schema.gotpl": &bintree{templatesJsonSchemaGotpl, map[string]*bintree{}},
		"spec-md.gotpl": &bintree{templatesSpecMdGotpl, map[string]*bintree{}},
		"toc-md.gotpl": &bintree{templatesTocMdGotpl, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

