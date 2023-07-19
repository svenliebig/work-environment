package git

import (
	"fmt"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
)

type ProjectStatus struct {
	modified []string
	added    []string
	deleted  []string
}

func Status(path string) (*ProjectStatus, error) {
	r, err := getRepository(path)

	if err != nil {
		return nil, err
	}

	w, err := r.Worktree()

	if err != nil {
		return nil, err
	}

	p, err := gitignore.ReadPatterns(w.Filesystem, []string{})

	if err != nil {
		return nil, err
	}

	w.Excludes = p

	status, err := w.Status()

	if err != nil {
		return nil, err
	}

	ps := &ProjectStatus{}

	for key, s := range status {
		switch s.Staging {
		case gogit.Modified, gogit.Renamed, gogit.Copied, gogit.Unmodified, gogit.UpdatedButUnmerged:
			ps.modified = append(ps.modified, key)
		case gogit.Deleted:
			ps.deleted = append(ps.deleted, key)
		case gogit.Added, gogit.Untracked:
			ps.added = append(ps.added, key)
		}
	}

	return ps, nil
}

func (ps *ProjectStatus) String() string {
	s := ""

	if len(ps.added) > 0 {
		s += fmt.Sprintf("+%d", len(ps.added))
	}

	if len(ps.modified) > 0 {
		s += fmt.Sprintf("~%d", len(ps.modified))
	}

	if len(ps.deleted) > 0 {
		s += fmt.Sprintf("-%d", len(ps.deleted))
	}

	return s
}

func (ps *ProjectStatus) Dirty() bool {
	return len(ps.added) > 0 || len(ps.modified) > 0 || len(ps.deleted) > 0
}
