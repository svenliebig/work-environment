package bamboo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/svenliebig/work-environment/pkg/utils/rest"
)

type ResultsResult struct {
	Results struct {
		Size   int `json:"size"`
		Result []struct {
			BuildNumber    int    `json:"buildNumber"`
			BuildResultKey string `json:"buildResultKey"`
			BuildState     string `json:"buildState"`
		} `json:"result"`
	} `json:"results"`
}

func (c *Client) Results(key string, amount int) (*ResultsResult, error) {
	headers := make(map[string]string)
	headers["Accept"] = "application/json"

	res, err := c.get(context.TODO(), fmt.Sprintf("/result/%s?expand=results%%5B0%%3A%d%%5D", key, amount-1), &rest.Options{
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

	var result ResultsResult

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

type ResultResult struct {
	Stages struct {
		Stage []struct {
			State   string `json:"State"`
			Results struct {
				Result []struct {
					BuildNumber int `json:"buildNumber"`
					LogEntries  struct {
						LogEntry []struct {
							Log         string `json:"log"`
							UnstyledLog string `json:"unstyledLog"`
						} `json:"logEntry"`
					} `json:"logEntries"`
					LogFiles []string `json:"logFiles"`
				} `json:"result"`
			} `json:"results"`
		} `json:"stage"`
	} `json:"stages"`
}

func (c *Client) Result(buildResultKey string) (*ResultResult, error) {
	headers := make(map[string]string)
	headers["Accept"] = "application/json"

	res, err := c.get(context.TODO(), fmt.Sprintf("/result/%s?expand=stages.stage.results.result.logEntries%%5B-26%%3A%%5D", buildResultKey), &rest.Options{
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

	var result ResultResult

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
