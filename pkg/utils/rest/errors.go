package rest

import (
	"fmt"
	"net/http"
)

var (
	_ error = &RestGetError{}
	_ error = &RestGetRequestCreationError{}
)

type RestGetError struct {
	url     string
	headers http.Header
}

func (err *RestGetError) Error() string {
	return fmt.Sprintf("Error while trying to perform a rest.Get request to '%s'.", err.url)
}

type RestGetRequestCreationError struct {
	url string
}

func (err *RestGetRequestCreationError) Error() string {
	return fmt.Sprintf("Error while trying to create a rest.Get request with the url '%s'.", err.url)
}
