package server

import (
	"user-traits/internal/di"
)

type Server struct {
	depsFactory *di.DependenciesFactory
}

func NewServer(depsFactory *di.DependenciesFactory) *Server {
	return &Server{
		depsFactory: depsFactory,
	}
}
