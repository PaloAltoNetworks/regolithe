package magetask

import (
	"errors"
	"os"
	"path"

	"github.com/BurntSushi/toml"
)

type constraint struct {
	Name     string
	Revision string
	Version  string
	Branch   string
}

type project struct {
	Constraint []constraint
	Override   []constraint
	Projects   []constraint
}

type versionTemplate struct {
	ProjectVersion string
	ProjectSha     string
	Libs           map[string]string
}

func parseProject(proj project) map[string]string {

	data := map[string]string{}

	// Apply constraint first
	for _, c := range proj.Constraint {
		if c.Version != "" {
			data[c.Name] = c.Version
		} else if c.Branch != "" {
			data[c.Name] = c.Branch
		} else if c.Revision != "" {
			data[c.Name] = c.Revision
		}
	}

	// Then projects
	for _, p := range proj.Projects {
		if p.Version != "" {
			data[p.Name] = p.Version
		} else if p.Branch != "" {
			data[p.Name] = p.Branch
		} else if p.Revision != "" {
			data[p.Name] = p.Revision
		}
	}

	// Then overrides
	for _, o := range proj.Override {
		if o.Version != "" {
			data[o.Name] = o.Version
		} else if o.Branch != "" {
			data[o.Name] = o.Branch
		} else if o.Revision != "" {
			data[o.Name] = o.Revision
		}
	}

	return data
}

func makeVersionFromDep(folder string, outFolder string, projectVersion string, projectSha string) error {

	if projectVersion == "" || projectSha == "" {
		return errors.New("you must set both projectVersion and projectSha")
	}

	if folder == "" {
		folder = "./"
	}
	if outFolder == "" {
		outFolder = "./internal/versions"
	}

	var proj project

	if _, err := os.Stat("./Gopkg.lock"); err == nil {
		if _, err := toml.DecodeFile(path.Join(folder, "Gopkg.lock"), &proj); err != nil {
			return err
		}
	} else {
		if _, err := toml.DecodeFile(path.Join(folder, "Gopkg.toml"), &proj); err != nil {
			return err
		}
	}

	return writeVersionsFile(
		versionTemplate{
			ProjectVersion: projectVersion,
			ProjectSha:     projectSha,
			Libs:           parseProject(proj),
		},
		outFolder,
	)
}
