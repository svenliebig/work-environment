package config

import (
	"testing"
)

func TestWorkEnvironmentConfig_GetProjectByPath(t *testing.T) {
	t.Run("should find a work environment when beeing directly on the project path", func(t *testing.T) {
		c, err := GetConfig("/Users/sven.liebig/workspace/repositories/isbj/redesign/tagespflege")

		if err != nil {
			t.Fatal(err)
		}

		got, err := c.GetProjectByPath("/Users/sven.liebig/workspace/repositories/isbj/redesign/tagespflege")

		if err != nil {
			t.Fatal(err)
		}

		if got.Identifier != "tagespflege" {
			t.Fatalf("expected %s to be %q", got.Identifier, "tagespflege")
		}
	})

	t.Run("should find a work environment when beeing inside of a project path", func(t *testing.T) {
		c, err := GetConfig("/Users/sven.liebig/workspace/repositories/isbj/redesign/tagespflege/gui/src")

		if err != nil {
			t.Fatal(err)
		}

		got, err := c.GetProjectByPath("/Users/sven.liebig/workspace/repositories/isbj/redesign/tagespflege")

		if err != nil {
			t.Fatal(err)
		}

		if got.Identifier != "tagespflege" {
			t.Fatalf("expected %s to be %q", got.Identifier, "tagespflege")
		}
	})
}
