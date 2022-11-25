package ci

import (
	"errors"
	"sync"

	"github.com/svenliebig/work-environment/pkg/context"
)

type BranchPlan struct {
	Key string
}

type Client interface {
	GetBranchPlans() ([]*BranchPlan, error)
	GetPlanSuggestion() (string, error)
}

var (
	clients = make(map[string]ClientProvider)
	lock    = sync.RWMutex{}

	// @comm i would like to create a complete new error with that, buy
	// i remember we did not want this last time..
	ErrClientAlreadyRegistered = errors.New("client already registered")
	ErrNoSuchClient            = errors.New("no such client")
)

type ClientProvider func(ctx *context.Context) Client

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

func UseClient(ctx *context.Context, citype string) (Client, error) {
	lock.RLock()
	defer lock.RUnlock()

	if p, ok := clients[citype]; ok {
		return p(ctx), nil
	}
	return nil, ErrNoSuchClient
}
