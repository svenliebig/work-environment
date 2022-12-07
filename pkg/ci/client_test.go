package ci

import (
	"fmt"
	"testing"

	"github.com/svenliebig/work-environment/pkg/context"
)

type clientImpl struct {
}

func (c clientImpl) GetBranchPlans() ([]*BranchPlan, error) {
	return make([]*BranchPlan, 0), nil
}

func (c clientImpl) GetPlanSuggestion() (string, error) {
	return "", nil
}

func (c clientImpl) LatestBuildResult() (*BuildResult, error) {
	return nil, nil
}

func p(_ *context.Context) Client {
	return clientImpl{}
}

func TestRegisterClient(t *testing.T) {
	t.Run("should register a client", func(t *testing.T) {
		if err := RegisterClient("bamboo", p); err != nil {
			t.Errorf("RegisterClient() error = %v", err)
		}
	})
}

// BenchmarkRegisterClient-8   	 4086864	       292.0 ns/op	     143 B/op	       2 allocs/op

func BenchmarkRegisterClient(b *testing.B) {
	for n := 0; n < b.N; n++ {
		RegisterClient(fmt.Sprintf("bamboo-%d", n), p)
	}
}
