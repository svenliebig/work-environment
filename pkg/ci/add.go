package ci

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/svenliebig/work-environment/pkg/config"
	"github.com/svenliebig/work-environment/pkg/core"
	"github.com/svenliebig/work-environment/pkg/utils/bamboo"
	"github.com/svenliebig/work-environment/pkg/utils/cli"
)

var (
	ErrNoCIDefined          = errors.New("no ci defined in this work environment")
	ErrRequiredIdentifierCI = errors.New("ci identifier is required when there is more than one CI defined in your work environment")
	ErrProjectNotFound      = errors.New("cannot find the project in your .work-environment/config.json, make sure it exists")
	ErrNoSuggestionFound    = errors.New("not able to find any suggestion for the project")
	ErrNoKeyProvided        = errors.New("no project key for the ci provided, try --suggest or provide a key")
)

func Add(p string, ciId string, projectId string, key string, suggest bool) error {
	// @comm i would like to have this different here...
	c, err := config.GetConfig(p)

	if err != nil {
		return fmt.Errorf("%w: error while trying to read the work-environment configuration", err)
	}

	cic, err := GetConfig(p)

	if err != nil {
		return fmt.Errorf("%w: error while trying to read the ci configuration", err)
	}

	var project *core.Project

	if projectId == "" {
		project, err = c.GetProjectByPath(p)

		if err != nil {
			if errors.Is(err, config.ErrConfigDoesNotExist) {
				return ErrProjectNotFound
			}
			return fmt.Errorf("%w: error while trying the project by path", err)
		}
	}

	ci, err := cic.GetEnvironment(ciId)

	if err != nil {
		return fmt.Errorf("%w: error while get the ci environment by ci id", err)
	}

	if suggest {
		key, err := getSuggestions(ci, project.Identifier)

		if err != nil {
			return err
		}

		for _, p := range c.Projects {
			if p.Identifier == project.Identifier {
				p.CI = &core.ProjectCI{
					Id:         ci.Identifier,
					ProjectKey: key,
				}
			}
		}

		err = c.Write()

		if err != nil {
			return err
		}

		fmt.Printf("\nAdded %q ci with project key %q to project %q.\n", ci.Identifier, key, project.Identifier)
	} else {
		if key == "" {
			return ErrNoKeyProvided
		}

		// TODO validate if key is a ci plan
	}

	return nil
}

func getSuggestions(ci *CI, identifier string) (string, error) {
	// TODO make the client abstract instead of bamboo
	c := bamboo.Client{
		BaseUrl:   ci.Url,
		AuthToken: ci.AuthToken,
	}

	sr, err := c.SearchPlans(identifier)

	if err != nil {
		return "", err
	}

	if sr.Size == 0 {
		return "", ErrNoSuggestionFound
	}

	if sr.Size == 1 {
		return sr.SearchResults[0].Id, nil
	} else {
		availabeAnswers := make([]string, sr.Size)
		availableSuggestions := ""

		for i, sr := range sr.SearchResults {
			availableSuggestions += fmt.Sprintf("\n\t[%d] %s - %s", i, sr.SearchEntity.PlanName, sr.SearchEntity.ProjectName)
			availabeAnswers[i] = fmt.Sprintf("%d", i)
		}

		q := fmt.Sprintf("\nFound multiple hits for %q in %q:%s\nplease select one:", identifier, "bamboo", availableSuggestions)

		answer := cli.Question(q, availabeAnswers)

		index, err := strconv.Atoi(answer)

		if err != nil {
			return "", err
		}

		return sr.SearchResults[index].Id, nil
	}
}
