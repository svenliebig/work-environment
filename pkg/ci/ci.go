package ci

import (
	"fmt"

	"github.com/svenliebig/work-environment/pkg/utils/bamboo"
)

type CI struct {
	CiType     string
	Identifier string
	AuthToken  string
	Url        string
}

func (c *CI) GetClient() (*bamboo.Client, error) {
	if c.CiType == "bamboo" {
		return &bamboo.Client{
			BaseUrl:   c.Url,
			AuthToken: c.AuthToken,
		}, nil
	}
	return nil, fmt.Errorf("ci type %q not supported", c.CiType)
}
