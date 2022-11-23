package ci

import "testing"

func TestList(t *testing.T) {

	t.Run("should use the tabwriter", func(t *testing.T) {
		List("/Users/sven.liebig/workspace/repositories/isbj/commons/ansible-paas")
	})
}
