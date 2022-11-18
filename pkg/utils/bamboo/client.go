package bamboo

import (
	"context"
	"fmt"
	"net/http"

	"github.com/svenliebig/work-environment/pkg/utils/rest"
)

// creates a client to communicate with an instance of bamboo
//
// needs atleast the BaseUrl and AuthToken to work
type Client struct {
	BaseUrl   string
	AuthToken string
}

func (c *Client) Get(ctx context.Context, path string, o *rest.Options) (*http.Response, error) {
	url := fmt.Sprintf("%s/rest/api/latest%s", c.BaseUrl, path)
	o.Headers["Authorization"] = fmt.Sprintf("Basic %s", c.AuthToken)
	return rest.Get(ctx, url, o)
}
