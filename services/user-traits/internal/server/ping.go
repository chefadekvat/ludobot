package server

import (
	"context"
	"user-traits/gen/api"
)

func (s *Server) GetPing(context.Context, api.GetPingRequestObject) (api.GetPingResponseObject, error) {
	return api.GetPing200TextResponse("pong"), nil
}
