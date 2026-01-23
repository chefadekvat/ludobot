package di

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependencies struct {
	Context   context.Context
	DbPoolPtr *pgxpool.Pool
}

type DependenciesFactory struct {
	dbPool *pgxpool.Pool
}

func NewDependenciesFactory(dbPool *pgxpool.Pool) *DependenciesFactory {
	return &DependenciesFactory{
		dbPool: dbPool,
	}
}

func (f *DependenciesFactory) CreateDependencies(ctx context.Context) *Dependencies {
	return &Dependencies{
		Context:   ctx,
		DbPoolPtr: f.dbPool,
	}
}
