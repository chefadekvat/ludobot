package di

import (
	"context"
	"ludobot/internal/arguments"
)

type Dependencies struct {
	Args    arguments.Arguments
	Context context.Context
}
