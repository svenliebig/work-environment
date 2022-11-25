package ci

import (
	"testing"

	"github.com/svenliebig/work-environment/pkg/context"
)

func TestList(t *testing.T) {

	t.Run("should use the tabwriter", func(t *testing.T) {
		err := List(&context.Context{Path: "/Users/sven.liebig/workspace/repositories/isbj/commons/ansible-paas"})

		if err != nil {
			t.Fatal(err)
		}
	})
}
