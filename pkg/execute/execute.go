package execute

import (
	"fmt"

	"github.com/svenliebig/work-environment/pkg/context"
)

func Do(ctx context.BaseContext) error {
	fmt.Println("do nothing")
	return nil
}
