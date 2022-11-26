package we

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/mattn/go-zglob"
	"github.com/svenliebig/work-environment/pkg/core"
	"github.com/svenliebig/work-environment/pkg/utils/git"
)

type InitializeOptions struct {
	Override bool
}

// TODO beautify
func Do(p string, o *InitializeOptions) error {
	configPath := filepath.Join(p, core.ConfigurationFileName)

	if o == nil || !o.Override {
		fi, err := os.Stat(configPath)

		if err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Fatal(err)
		}

		if fi != nil {
			fmt.Printf("configuration does already exist in:\n  %q\n\nif you want to override the existing configuration run:\n  we init --override\n", configPath)
			return nil
		}
	}

	projects, err := scanForProjects(p)

	if err != nil {
		log.Fatal(err)
	}

	// TODO warnings for same identifier

	if o.Override {
		err = os.Remove(configPath)

		if err != nil {
			log.Fatal(err)
		}

		err := write(configPath, projects)

		if err != nil {
			return err
		}

		fmt.Printf("Overwritten: Saved %d projects to your work-environment in %q.\n", len(projects), configPath)
	} else {
		err = write(configPath, projects)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Saved %d projects to your work-environment in %q.\n", len(projects), configPath)
	}

	return nil
}

func write(configPath string, projects []*core.Project) error {
	result, err := json.MarshalIndent(&core.Configuration{Projects: projects}, "", "  ")

	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(configPath, result, 0644)

	if err != nil {
		return err
	}

	return nil
}

func scanForProjects(p string) ([]*core.Project, error) {
	x := path.Join(p, "**", ".git")

	// TODO ignore node_modules and other heavy things we don't want
	dirs, err := zglob.Glob(x)

	if err != nil {
		return nil, err
	}

	projects := make([]*core.Project, len(dirs))
	for i, dir := range dirs {
		projectPath, _ := filepath.Split(dir)
		projectPath = filepath.Clean(projectPath)
		identifier := filepath.Base(projectPath)

		remoteUrl, err := git.RepositoryGetRemoteOriginUrl(projectPath)

		if err != nil {
			fmt.Printf("ERR while trying to create project for %q:\n%s", identifier, err)
		}

		projects[i] = &core.Project{
			Identifier: identifier,
			Path:       projectPath,
			Git: &core.ProjectGit{
				RemoteUrl: remoteUrl,
			},
		}
	}

	return projects, nil
}
