package rest

import (
	"context"
	"fmt"
	"net/http"
)

func Post(ctx context.Context, url string, options Options) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, options.Body)

	if err != nil {
		return nil, fmt.Errorf("%w,  %s", &RestGetRequestCreationError{url}, err)
	}

	req = req.WithContext(ctx)
	setHeaders(req.Header, options.Headers)
	response, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("%w,  %s", &RestGetError{url: url, headers: req.Header}, err)
	}

	return response, nil
}
