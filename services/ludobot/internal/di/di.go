package di

import (
	"context"
	"ludobot/internal/infrastructure/arguments"
)

type Dependencies struct {
	Args    arguments.Arguments
	Context context.Context
}
