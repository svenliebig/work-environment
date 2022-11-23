package bamboo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/svenliebig/work-environment/pkg/utils/rest"
)

type PlansSearchResult struct {
	Size          int `json:"size"`
	SearchResults []struct {
		Id           string `json:"id"`
		SearchEntity struct {
			ProjectName string `json:"projectName"`
			PlanName    string `json:"planName"`
		} `json:"searchEntity"`
	} `json:"searchResults"`
}

func (c *Client) SearchPlans(searchTerm string) (*PlansSearchResult, error) {
	headers := make(map[string]string)
	headers["Accept"] = "application/json"

	res, err := c.get(context.TODO(), fmt.Sprintf("/search/plans?searchTerm=%s", searchTerm), &rest.Options{
		Headers: headers,
	})

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		if res.StatusCode == 401 {
			return nil, ErrUnauthorized
		}

		return nil, errors.New(string(body))
	}

	var result PlansSearchResult

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
