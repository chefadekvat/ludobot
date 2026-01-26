.PHONY: games-calculator

build-testsuite: games-calculator

testsuite: build-testsuite

games-calculator: mod
	go build -o build/games-calculator cmd/main.go

codegen: mod
	go tool oapi-codegen -config ./docs/config/oapi.yaml -o ./gen/api/api.go ./docs/api/api.yaml

mod:
	go mod download

test: codegen
	go test -v ./...

clean:
	-rm -r ./build/ ./gen/