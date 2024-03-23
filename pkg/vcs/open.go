package vcs

import (
	"fmt"

	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/utils/browser"
)

type OpenParameters struct {
	// opens the pull request page instead of the repository page
	PullRequest bool
}

func Open(ctx context.ProjectContext, params OpenParameters) error {
	vcs, err := ctx.GetVCS()

	if err != nil {
		return err
	}

	client, err := UseClient(ctx, vcs)

	if err != nil {
		return err
	}

	var url string
	if params.PullRequest {
		url, err = client.PullRequestWebURL()
	} else {
		url, err = client.WebURL()
	}

	fmt.Println(url)

	if err != nil {
		return err
	}

	err = browser.Open(url)

	if err != nil {
		return err
	}

	return nil

}
