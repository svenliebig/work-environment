package bamboo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/svenliebig/work-environment/pkg/utils/rest"
)

type BranchesSearchResult struct {
	Size          int `json:"size"`
	SearchResults []struct {
		Id           string `json:"id"`
		SearchEntity struct {
			Key        string `json:"key"`
			BranchName string `json:"branchName"`
		} `json:"searchEntity"`
	} `json:"searchResults"`
}

func (c *Client) SearchBranches(planKey string, searchTerm string) (*BranchesSearchResult, error) {
	headers := make(map[string]string)
	headers["Accept"] = "application/json"

	res, err := c.Get(context.TODO(), fmt.Sprintf("/search/branches?masterPlanKey=%s&searchTerm=%s", planKey, searchTerm), &rest.Options{
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

	var result BranchesSearchResult

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
