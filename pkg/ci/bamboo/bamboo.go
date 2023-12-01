package bamboo

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/svenliebig/work-environment/pkg/api/bamboo"
	"github.com/svenliebig/work-environment/pkg/ci"
	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/utils/cli"
)

// TODO this could be in the CI package, to be part of the API
var (
	ErrNoCIDefined          = errors.New("no ci defined in this work environment")
	ErrRequiredIdentifierCI = errors.New("ci identifier is required when there is more than one CI defined in your work environment")
	ErrProjectNotFound      = errors.New("cannot find the project in your .work-environment/config.json, make sure it exists")
	ErrNoSuggestionFound    = errors.New("not able to find any suggestion for the project")
	ErrNoCDFound            = errors.New("not able to find any CD for the project")
	ErrNoKeyProvided        = errors.New("no project key for the ci provided, try --suggest or provide a key")
)

type client struct {
	ctx context.ProjectContext

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

// returns the main branch plan ci url
func (c *client) GetPlanUrl() (string, error) {
	bc, err := c.bamboo()

	if err != nil {
		return "", err
	}

	p := c.ctx.Project()

	if err != nil {
		return "", err
	}

	key := p.CI.ProjectKey

	return fmt.Sprintf("%s/browse/%s", bc.BaseUrl, key), nil
}

// returns the ci url of the current branch, or the main branch plan if none is found
func (c *client) GetBranchPlanUrl() (string, error) {
	bc, err := c.bamboo()

	if err != nil {
		return "", err
	}

	p := c.ctx.Project()

	if err != nil {
		return "", err
	}

	var key string

	branchPlans, err := c.GetBranchPlans()

	if err != nil {
		return "", err
	}

	if len(branchPlans) == 1 {
		key = branchPlans[0].Key
	} else {
		key = p.CI.ProjectKey
	}

	return fmt.Sprintf("%s/browse/%s", bc.BaseUrl, key), nil
}

// GetBranchPlans implements ci.Client
func (c *client) GetBranchPlans() ([]*ci.BranchPlan, error) {
	bc, err := c.bamboo()

	if err != nil {
		return nil, err
	}

	p := c.ctx.Project()

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

// LatestBuildResult implements ci.Client
func (c *client) LatestBuildResult() (*ci.BuildResult, error) {
	bc, err := c.bamboo()

	if err != nil {
		return nil, err
	}

	p := c.ctx.Project()

	if err != nil {
		return nil, err
	}

	var key string

	branchPlans, err := c.GetBranchPlans()

	if err != nil {
		return nil, err
	}

	if len(branchPlans) == 1 {
		key = branchPlans[0].Key
	} else {
		key = p.CI.ProjectKey
	}

	results, err := bc.Results(key, 1)

	if err != nil {
		return nil, err
	}

	if len(results.Results.Result) != 1 {
		return nil, fmt.Errorf("%w: could not find a build result for %s", ci.ErrBuildResultNotFound, p.CI.ProjectKey)
	}

	plan, err := bc.Plan(p.CI.ProjectKey)

	if err != nil {
		return nil, err
	}

	result, err := bc.Result(results.Results.Result[0].BuildResultKey)

	if err != nil {
		return nil, err
	}

	logs := make([]string, 0)

	for _, s := range result.Stages.Stage {
		if s.State == "Failed" {
			for _, l := range s.Results.Result[0].LogEntries.LogEntry {
				x := strings.TrimSuffix(fmt.Sprintf("%s\n", l.UnstyledLog), "\n")
				logs = append(logs, x)
			}
		}
	}

	return &ci.BuildResult{
		Success:     results.Results.Result[0].BuildState == "Successful",
		BuildNumber: fmt.Sprintf("%s-%d", key, results.Results.Result[0].BuildNumber),
		IsBuilding:  plan.IsBuilding,
		LogUrl:      fmt.Sprintf("%s/browse/%s/log", c.bambooClient.BaseUrl, results.Results.Result[0].BuildResultKey),
		Logs:        logs,
	}, nil
}

// GetPlanSuggestions implements ci.Client
func (c *client) GetPlanSuggestion() (string, error) {
	bc, err := c.bamboo()

	if err != nil {
		return "", err
	}

	p := c.ctx.Project()
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

// GetPlanSuggestions implements ci.Client
func (c *client) GetCD() (int, error) {
	bc, err := c.bamboo()

	if err != nil {
		return 0, err
	}

	p := c.ctx.Project()
	sr, err := bc.DeployProjectForPlan(p.CI.ProjectKey)

	if err != nil {
		return 0, err
	}

	if len(sr) == 0 {
		return 0, ErrNoCDFound
	}

	if len(sr) == 1 {
		return sr[0].Id, nil
	} else {
		availabeAnswers := make([]string, len(sr))
		availableSuggestions := ""

		for i, sr := range sr {
			availableSuggestions += fmt.Sprintf("\n\t[%d] %d - %s", i, sr.Id, sr.Name)
			availabeAnswers[i] = fmt.Sprintf("%d", i)
		}

		q := fmt.Sprintf("\nFound multiple hits for %q in %q:%s\nplease select one:", p.Identifier, "bamboo", availableSuggestions)

		answer := cli.Question(q, availabeAnswers)

		index, err := strconv.Atoi(answer)

		if err != nil {
			return 0, err
		}

		return sr[index].Id, nil
	}
}

func init() {
	ci.RegisterClient("bamboo", func(ctx context.ProjectContext) ci.Client {
		return &client{
			ctx: ctx,
		}
	})
}
