package we

import (
	"errors"

	"github.com/svenliebig/work-environment/pkg/context"
)

func Info() error {
	var ctx context.BaseContext
	ctx, err := context.CreateProjectContext()

	if errors.Is(err, context.ErrNoSuchProjectInDirectory) {
		ctx, err = context.CreateBaseContext()
	}

	if err != nil {
		return err
	}

	cw := ctx.Info()

	cw.Print()

	return nil
}
