package ci

import (
	"testing"

	"github.com/svenliebig/work-environment/pkg/context"
)

func TestList(t *testing.T) {

	t.Run("should use the tabwriter", func(t *testing.T) {
		ctx := &context.BaseContext{Cwd: "/Users/sven.liebig/workspace/repositories/isbj/commons/ansible-paas"}
		err := ctx.Validate()

		if err != nil {
			t.Fatal(err)
		}

		err = List(ctx)

		if err != nil {
			t.Fatal(err)
		}
	})
}
