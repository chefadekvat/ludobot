package server

import (
	"context"
	"errors"
	"log/slog"
	"user-traits/gen/api"
	"user-traits/internal/usecases"
)

func (s *Server) PostV1UserCreate(ctx context.Context, request api.PostV1UserCreateRequestObject) (api.PostV1UserCreateResponseObject, error) {
	deps := s.depsFactory.CreateDependencies(ctx)

	useCase, err := usecases.NewUserCreationUseCase(deps)
	if err != nil {
		return api.PostV1UserCreate409JSONResponse{
			Code:    "Internal error",
			Message: "Failed to create use case: " + err.Error(),
		}, nil
	}
	defer useCase.Close()

	err = useCase.CreateUser(request.Body.Id, int64(request.Body.Balance))

	if err != nil {
		if errors.Is(err, usecases.ErrUserExists) {
			return api.PostV1UserCreate409JSONResponse{
				Code:    "user_exists",
				Message: err.Error(),
			}, nil
		}

		slog.Error(err.Error())

		return api.PostV1UserCreate500JSONResponse{
			Code:    "internal_server_error",
			Message: "Internal server error",
		}, err
	}

	return api.PostV1UserCreate204Response{}, nil
}
