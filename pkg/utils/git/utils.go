package git

import (
	"errors"
	"path"

	gogit "github.com/go-git/go-git/v5"
	"github.com/svenliebig/work-environment/pkg/utils/wepath"
)

func getRepository(p string) (*gogit.Repository, error) {
	r, err := gogit.PlainOpen(p)

	if err != nil {
		if errors.Is(err, gogit.ErrRepositoryNotExists) && !wepath.IsRoot(p) {
			return getRepository(path.Join(p, ".."))
		}

		return nil, err
	}

	return r, err
}
