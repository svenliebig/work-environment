package git

import (
	"fmt"
	"io"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
)

func BranchName(path string) (string, error) {
	r, err := getRepository(path)

	if err != nil {
		return "", err
	}

	h, err := r.Head()

	if err != nil {
		return "", err
	}

	if h.Name().IsBranch() {
		return h.Name().Short(), nil
	}

	return "", nil
}

func DefaultBranchName(path string) (string, error) {
	r, err := getRepository(path)

	if err != nil {
		return "", err
	}

	ref, err := r.Head()

	references, err := r.References()

	if err != nil {
		return "", err
	}

	for {
		n, err := references.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}

		if n.Type() == plumbing.SymbolicReference && n.Name().String() == "refs/remotes/origin/HEAD" {
			return strings.TrimPrefix(n.Target().Short(), "origin/"), nil
		}
	}

	if err != nil {
		return "", err
	}

	fmt.Println(ref.Name().Short())

	return ref.Name().Short(), nil
}
