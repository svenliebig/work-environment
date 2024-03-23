package vcs

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/core"
)

var (
	clients = map[string]ClientFactory{}
	setups  = map[string]Setup{}
	lock    = sync.RWMutex{}

	ErrClientAlreadyRegistered = errors.New("vcs client already registered")
	ErrSetupAlreadyRegistered  = errors.New("vcs setup already registered")
	ErrNoSuchClient            = errors.New("no such vcs client")
	ErrNoSuchSetup             = errors.New("no such vcs setup")
	ErrRepositoryNotFound      = errors.New("repository not found")
)

type Client interface {
	// is called immediately after the client is attached to the project context
	// and is used to configure the client with the necessary information, the
	// returned value is saved on the property configuration of the context.
	Configure() (string, error)

	// prints information about the repository.
	Info() error

	// returns the web URL of the repository.
	WebURL() (string, error)

	// returns the web URL of the pull request.
	PullRequestWebURL() (string, error)
}

type ClientFactory func(ctx context.ProjectContext) Client
type Setup func(ctx context.BaseContext) (string, error)

func SetupClient(ctx context.BaseContext, environment core.VCS) (core.VCS, error) {
	setup, err := UseSetup(environment.Type)

	if err != nil {
		return environment, err
	}

	configuration, err := setup(ctx)

	if err != nil {
		return environment, err
	}

	environment.Configuration = configuration

	return environment, nil
}

func RegisterClient(name string, factory ClientFactory) error {
	lock.Lock()
	defer lock.Unlock()

	if _, ok := clients[name]; ok {
		return ErrClientAlreadyRegistered
	}

	clients[name] = factory
	return nil
}

func RegisterSetup(name string, setup Setup) error {
	lock.Lock()
	defer lock.Unlock()

	if _, ok := setups[name]; ok {
		return ErrSetupAlreadyRegistered
	}

	setups[name] = setup
	return nil
}

func AvailableClients() []string {
	lock.RLock()
	defer lock.RUnlock()

	names := make([]string, 0, len(clients))
	for name := range clients {
		names = append(names, name)
	}

	return names
}

func UseClient(ctx context.ProjectContext, vcse *core.VCS) (Client, error) {
	lock.RLock()
	defer lock.RUnlock()

	if p, ok := clients[vcse.Type]; ok {
		return p(ctx), nil
	}

	return nil, fmt.Errorf(
		"%w: tried to use client '%q', but no such client is registered. Available clients are: '%s'",
		ErrNoSuchClient,
		vcse.Identifier,
		strings.Join(AvailableClients(), "', '"),
	)
}

func UseSetup(t string) (Setup, error) {
	lock.RLock()
	defer lock.RUnlock()

	if setup, ok := setups[t]; ok {
		return setup, nil
	}

	return nil, fmt.Errorf(
		"%w: tried to use setup '%s', but no such setup is registered. Available setups are: '%s'",
		ErrNoSuchSetup,
		t,
		strings.Join(AvailableClients(), "', '"),
	)
}
