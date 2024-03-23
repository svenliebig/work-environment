package vcs

import "github.com/svenliebig/work-environment/pkg/context"

func Info(ctx context.ProjectContext) error {
	vcse, err := ctx.GetVCS()

	if err != nil {
		return err
	}

	client, err := UseClient(ctx, vcse)

	if err != nil {
		return err
	}

	return client.Info()
}
