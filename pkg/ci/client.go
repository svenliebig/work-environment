package ci

import (
	"errors"
	"fmt"
	"sync"

	"github.com/svenliebig/work-environment/pkg/context"
)

type BranchPlan struct {
	Key string
}

type BuildResult struct {
	Success     bool
	BuildNumber string
	Logs        []string
	IsBuilding  bool
	LogUrl      string
}

type Client interface {
	GetBranchPlans() ([]*BranchPlan, error)
	GetPlanSuggestion() (string, error)
	LatestBuildResult() (*BuildResult, error)
	GetCD() (int, error)
}

var (
	clients = make(map[string]ClientProvider)
	lock    = sync.RWMutex{}

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
