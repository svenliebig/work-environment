package vcs

import "github.com/svenliebig/work-environment/pkg/context"

func Info(ctx context.ProjectContext) error {
	p := ctx.Project()

	println(p.Identifier)
	vcse, err := ctx.Configuration().GetVCSEnvironmentById("azure-tp6")

	if err != nil {
		return err
	}

	client, err := UseClient(ctx, vcse)

	if err != nil {
		return err
	}

	return client.Info()
}
