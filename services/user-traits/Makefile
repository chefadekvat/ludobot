.PHONY: user-traits

testsuite: user-traits
	docker compose -f testsuite/docker/docker-compose.yaml build
	docker compose -f testsuite/docker/docker-compose.yaml run testsuite pytest

user-traits: codegen
	go build -o ./build/user-traits cmd/main.go

codegen: mod sqlgen
	go tool oapi-codegen -config config/oapi/config.yaml -o ./gen/api/api.go ./docs/api.yaml

sqlgen:
	go run ${SQLGEN_ROOT}/main.go ./sql ./gen/sql

mod:
	go mod download

test: codegen
	go test -v ./...

clean:
	-rm -r ./build/ ./gen/