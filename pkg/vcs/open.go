package vcs

import (
	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/utils/browser"
)

func Open(ctx context.ProjectContext) error {
	vcs, err := ctx.GetVCS()

	if err != nil {
		return err
	}

	client, err := UseClient(ctx, vcs.Type)

	if err != nil {
		return err
	}

	url, err := client.WebURL()

	if err != nil {
		return err
	}

	err = browser.Open(url)

	if err != nil {
		return err
	}

	return nil

}
