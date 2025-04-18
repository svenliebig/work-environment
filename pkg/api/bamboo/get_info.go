package bamboo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/svenliebig/work-environment/pkg/utils/rest"
)

type infoResult struct {
	Version string `json:"version"`
}

func (c *Client) GetInfo(ctx context.Context) (*infoResult, error) {
	url := fmt.Sprintf("%s/rest/api/latest/info", c.BaseUrl)
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.AuthToken)
	headers["Accept"] = "application/json"

	res, err := rest.Get(ctx, url, &rest.Options{
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

	var result infoResult

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
