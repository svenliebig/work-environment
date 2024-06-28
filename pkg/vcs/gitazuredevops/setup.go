package gitazuredevops

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/core"
	"github.com/svenliebig/work-environment/pkg/utils/cli"
	"github.com/svenliebig/work-environment/pkg/vcs"
)

type envConfig struct {
	Organization string
	AccessToken  string
}

func init() {
	vcs.RegisterSetup(identifier, func(ctx context.BaseContext) (string, error) {
		config := &envConfig{}

		fmt.Println(cli.Bold("\nSetup Azure DevOps"))
		fmt.Println()

		answer := cli.QuestionFree("Please enter the organization name (https://dev.azure.com/{organization}): ")
		config.Organization = strings.Trim(answer, " ")
		answer = cli.QuestionFree("Access Token: ")
		config.AccessToken = strings.Trim(answer, " ")

		if result, err := json.Marshal(config); err == nil {
			return string(result), nil
		}

		return "", nil
	})
}

func (c *client) environment() (*envConfig, error) {
	if c._env != nil {
		return c._env, nil
	}

	vcse, err := c.ctx.GetVCS()

	if err != nil {
		return nil, err
	}

	config, err := environment(vcse)

	if err != nil {
		return nil, err
	}

	c._env = config

	return config, nil
}

func environment(vcse *core.VCS) (*envConfig, error) {
	config := &envConfig{}

	if err := json.Unmarshal([]byte(vcse.Configuration), config); err != nil {
		return nil, err
	}

	return config, nil
}
