package bamboo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/svenliebig/work-environment/pkg/utils/rest"
)

type DeploymentResultsByEnvironmentIdResult struct {
	Size    int `json:"size"`
	Results []struct {
		Id                    int    `json:"id"`
		DeploymentState       string `json:"deploymentState"`
		FinishedDate          int    `json:"finishedDate"`
		DeploymentVersionName string `json:"deploymentVersionName"`
	} `json:"results"`
}

func (c *Client) DeploymentResultsByEnvironmentId(environmentId int, maxResults int) (*DeploymentResultsByEnvironmentIdResult, error) {
	headers := make(map[string]string)
	headers["Accept"] = "application/json"

	res, err := c.get(context.TODO(), fmt.Sprintf("/deploy/environment/%d/results?max-results=%d", environmentId, maxResults), &rest.Options{
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

	var result *DeploymentResultsByEnvironmentIdResult

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}
