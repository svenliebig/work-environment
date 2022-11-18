package core

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
