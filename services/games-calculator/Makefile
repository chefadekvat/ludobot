.PHONY: games-calculator

games-calculator: mod
	go build -o build/games-calculator cmd/main.go

codegen: mod
	go tool oapi-codegen -config ./docs/config/oapi.yaml -o ./gen/server/api.go ./docs/api/api.yaml

mod:
	go mod tidy
	go mod download