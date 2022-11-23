package core

import "github.com/svenliebig/work-environment/pkg/utils/git"

type ProjectGit struct {
	RemoteUrl string
}

type ProjectCI struct {
	Id         string
	ProjectKey string
}

type Project struct {
	Identifier string
	Path       string
	Git        *ProjectGit
	CI         *ProjectCI
}

// uses the path of the project and the git package
// to get the current branch of the project
func (p *Project) GetBranchName() (string, error) {
	return git.BranchName(p.Path)
}
