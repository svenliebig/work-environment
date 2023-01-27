package bamboo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/svenliebig/work-environment/pkg/utils/rest"
)

type DeployProjectForPlanResult struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *Client) DeployProjectForPlan(planKey string) ([]DeployProjectForPlanResult, error) {
	headers := make(map[string]string)
	headers["Accept"] = "application/json"

	res, err := c.get(context.TODO(), fmt.Sprintf("/deploy/project/forPlan?planKey=%s", planKey), &rest.Options{
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

	var result []DeployProjectForPlanResult

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}
