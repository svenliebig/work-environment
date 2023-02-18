package cd

import (
	"errors"
	"fmt"
	"sync"

	"github.com/svenliebig/work-environment/pkg/context"
)

type ClientInfo struct {
	Identifier string
	Type       string
	URL        string
	Version    string
	ProjectId  int
}

type Environment struct {
	Name string
	Id   int
}

type DeployResult struct {
	Id              int
	Version         string
	DeploymentState string
	Finished        int
}

type Client interface {
	Open() error
	Info() (*ClientInfo, error)
	Environments() ([]*Environment, error)
	DeployResult(environmentId int) (*DeployResult, error)
}

var (
	clients = make(map[string]ClientProvider)
	lock    = sync.RWMutex{}

	// @comm i would like to create a complete new error with that, but
	// i remember we did not want this last time..
	// @answ error wrapping
	ErrClientAlreadyRegistered = errors.New("client already registered")
	ErrNoSuchClient            = errors.New("no such client")

	ErrBuildResultNotFound = errors.New("was not able to find a build result")
)

type ClientProvider func(ctx context.ProjectContext) Client

func RegisterClient(citype string, p ClientProvider) error {
	lock.Lock()
	defer lock.Unlock()

	_, ok := clients[citype]
	if ok {
		return ErrClientAlreadyRegistered
	} else {
		clients[citype] = p
		return nil
	}
}

func UseClient(ctx context.ProjectContext, citype string) (Client, error) {
	lock.RLock()
	defer lock.RUnlock()

	if p, ok := clients[citype]; ok {
		return p(ctx), nil
	}
	return nil, fmt.Errorf("%w: tried to use client %q but not available", ErrNoSuchClient, citype)
}
