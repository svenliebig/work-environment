package vcs

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/svenliebig/work-environment/pkg/context"
)

var (
	clients = map[string]ClientFactory{}
	lock    = sync.RWMutex{}

	ErrClientAlreadyRegistered = errors.New("vcs client already registered")
	ErrNoSuchClient            = errors.New("no such vcs client")
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

func RegisterClient(name string, factory ClientFactory) error {
	lock.Lock()
	defer lock.Unlock()

	if _, ok := clients[name]; ok {
		return ErrClientAlreadyRegistered
	}

	clients[name] = factory
	return nil
}

func AvailableClients() []string {
	lock.RLock()
	defer lock.RUnlock()

	names := make([]string, len(clients))
	for name := range clients {
		names = append(names, name)
	}

	return names
}

func UseClient(ctx context.ProjectContext, t string) (Client, error) {
	lock.RLock()
	defer lock.RUnlock()

	if p, ok := clients[t]; ok {
		return p(ctx), nil
	}

	return nil, fmt.Errorf(
		"%w: tried to use client '%q', but no such client is registered. Available clients are: '%s'",
		ErrNoSuchClient,
		t,
		strings.Join(AvailableClients(), "', '"),
	)
}
