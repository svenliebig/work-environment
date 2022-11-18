package ci

import (
	"errors"
	"fmt"

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

	ci, err := cic.GetEnvironment("")

	if err != nil {
		return fmt.Errorf("%w: error while get the ci environment by ci id", err)
	}

	if project.CI != nil {
		// TODO CI get branch builds by masterPlanKey
		url := fmt.Sprintf("%s/browse/%s", ci.Url, project.CI.ProjectKey)
		err = browser.Open(url)

		if err != nil {
			return err
		}
	}

	return nil
}
