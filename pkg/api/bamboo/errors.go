package bamboo

import "errors"

var (
	ErrUnauthorized = errors.New("the bamboo server replied that your token is not authorized to access this data")
)
