package ci

import (
	"errors"
	"fmt"
	"strings"

	"github.com/svenliebig/work-environment/pkg/config"
	"github.com/svenliebig/work-environment/pkg/utils/browser"
)

func Open(p string) error {
	c, err := config.GetConfig(p)

	if err != nil {
		return fmt.Errorf("%w: error while trying to read the work-environment configuration", err)
	}

	cic, err := GetConfig(p)

	if err != nil {
		return fmt.Errorf("%w: error while trying to read the ci configuration", err)
	}

	project, err := c.GetProjectByPath(p)

	if err != nil {
		if errors.Is(err, config.ErrConfigDoesNotExist) {
			return ErrProjectNotFound
		}
		return fmt.Errorf("%w: error while trying the project by path", err)
	}

	ci, err := cic.GetEnvironment(project.CI.Id)

	if err != nil {
		return fmt.Errorf("%w: error while get the ci environment by ci id", err)
	}

	if project.CI != nil {
		b, err := project.GetBranchName()

		if err != nil {
			return err
		}

		b = strings.ReplaceAll(b, "/", "-")

		c, err := ci.GetClient()

		if err != nil {
			return err
		}

		r, err := c.SearchBranches(project.CI.ProjectKey, b)

		if err != nil {
			return err
		}

		if r.Size == 0 {
			url := fmt.Sprintf("%s/browse/%s", ci.Url, project.CI.ProjectKey)
			err = browser.Open(url)
		} else if r.Size == 1 {
			url := fmt.Sprintf("%s/browse/%s", ci.Url, r.SearchResults[0].SearchEntity.Key)
			err = browser.Open(url)
		} else {
			return fmt.Errorf("found more branches related to %q, but selecting one is not supported", b)
		}
	}

	return nil
}
