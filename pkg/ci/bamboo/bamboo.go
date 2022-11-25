package bamboo

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/svenliebig/work-environment/pkg/ci"
	"github.com/svenliebig/work-environment/pkg/context"
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

type client struct {
	ctx *context.Context

	bambooClient *bamboo.Client
}

func (c *client) bamboo() (*bamboo.Client, error) {
	if c.bambooClient == nil {
		ci, err := c.ctx.GetCI()

		if err != nil {
			return nil, err
		}

		c.bambooClient = &bamboo.Client{
			BaseUrl:   ci.Url,
			AuthToken: ci.AuthToken,
		}
	}

	return c.bambooClient, nil
}

// GetBranchPlans implements ci.Client
func (c *client) GetBranchPlans() ([]*ci.BranchPlan, error) {
	bc, err := c.bamboo()

	if err != nil {
		return nil, err
	}

	p, err := c.ctx.GetProject()

	if err != nil {
		return nil, err
	}

	b, err := p.GetBranchName()

	if err != nil {
		return nil, err
	}

	r, err := bc.SearchBranches(p.CI.ProjectKey, strings.ReplaceAll(b, "/", "-"))

	if err != nil {
		return nil, err
	}

	ret := make([]*ci.BranchPlan, r.Size)

	for i, result := range r.SearchResults {
		ret[i] = &ci.BranchPlan{
			Key: result.SearchEntity.Key,
		}
	}

	return ret, nil
}

// GetPlanSuggestions implements ci.Client
func (c *client) GetPlanSuggestion() (string, error) {
	bc, err := c.bamboo()

	if err != nil {
		return "", err
	}

	p, err := c.ctx.GetProject()

	if err != nil {
		return "", err
	}

	sr, err := bc.SearchPlans(p.Identifier)

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

		q := fmt.Sprintf("\nFound multiple hits for %q in %q:%s\nplease select one:", p.Identifier, "bamboo", availableSuggestions)

		answer := cli.Question(q, availabeAnswers)

		index, err := strconv.Atoi(answer)

		if err != nil {
			return "", err
		}

		return sr.SearchResults[index].Id, nil
	}
}

func init() {
	ci.RegisterClient("bamboo", func(ctx *context.Context) ci.Client {
		return &client{
			ctx: ctx,
		}
	})
}
