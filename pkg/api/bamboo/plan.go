package bamboo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/svenliebig/work-environment/pkg/utils/rest"
)

type PlanResult struct {
	IsBuilding bool `json:"isBuilding"`
}

func (c *Client) Plan(key string) (*PlanResult, error) {
	headers := make(map[string]string)
	headers["Accept"] = "application/json"

	res, err := c.get(context.TODO(), fmt.Sprintf("/plan/%s", key), &rest.Options{
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

	var result PlanResult

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
