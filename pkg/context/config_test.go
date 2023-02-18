package context

import (
	"testing"
)

func TestContext_GetConfiguration(t *testing.T) {
	t.Run("should find a configuration", func(t *testing.T) {
		c := &projectContext{
			Cwd: "/Users/sven.liebig/workspace/repositories/isbj/redesign/kita",
		}

		if err := c.validate(""); err != nil {
			t.Fatal(err)
		}

		_ = c.Configuration()
	})
}
