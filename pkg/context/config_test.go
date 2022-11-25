package context

import (
	"testing"
)

func TestContext_GetConfiguration(t *testing.T) {
	t.Run("should find a configuration", func(t *testing.T) {
		c := &Context{
			Path: "/Users/sven.liebig/workspace/repositories/isbj/redesign/kita",
		}

		_, err := c.GetConfiguration()

		if err != nil {
			t.Errorf("Context.GetConfiguration() error = %v", err)
			return
		}
	})
}
