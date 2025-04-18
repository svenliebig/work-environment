package git

import (
	"path"
	"runtime"
	"testing"
)

func TestGetIdentifier(t *testing.T) {
	t.Run("should return the correct git identifier", func(t *testing.T) {
		_, filename, _, _ := runtime.Caller(0)
		got, err := GetIdentifier(path.Dir(filename))
		expected := "work-environment"

		if err != nil {
			t.Fatalf("should not return an error, but got %v", err)
		}

		if got != expected {
			t.Errorf("got %q, expected %q", got, expected)
		}
	})
}
