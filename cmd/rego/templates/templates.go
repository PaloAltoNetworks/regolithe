package templates

import "embed"

//go:embed data
var fs embed.FS

// Get returns the template with the given name
func Get(name string) ([]byte, error) {
	return fs.ReadFile("data/" + name)
}
