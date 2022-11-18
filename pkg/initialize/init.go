package initialize

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/mattn/go-zglob"
	"github.com/svenliebig/work-environment/pkg/config"
	"github.com/svenliebig/work-environment/pkg/core"
	"github.com/svenliebig/work-environment/pkg/utils/git"
)

type InitializeOptions struct {
	Override bool
}

func Do(p string, o *InitializeOptions) {
	weDir := filepath.Join(p, ".work-environment")
	weConfigPath := filepath.Join(weDir, "config.json")

	if o == nil || !o.Override {
		fi, err := os.Stat(weConfigPath)

		if err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Fatal(err)
		}

		if fi != nil {
			fmt.Printf("configuration does already exist in:\n  %q\n\nif you want to override the existing configuration run:\n  we init --override\n", weConfigPath)
			return
		}
	}

	x := path.Join(p, "**", ".git")

	// TODO ignore node_modules and other heavy things we don't want
	dirs, err := zglob.Glob(x)

	if err != nil {
		log.Fatal(err)
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

	// TODO warnings for same identifier

	if o.Override {
		err = os.Remove(weConfigPath)

		if err != nil {
			log.Fatal(err)
		}

		err = config.Write(weConfigPath, projects)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Overwritten: Saved %d projects to your work-environment in %q.\n", len(projects), weConfigPath)
	} else {
		err = config.Write(weConfigPath, projects)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Saved %d projects to your work-environment in %q.\n", len(projects), weConfigPath)
	}
}
