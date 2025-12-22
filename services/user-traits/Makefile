.PHONY: user-traits

user-traits: mod codegen
	go build -o build/user-traits cmd/main.go

codegen:
	go tool oapi-codegen -config ./config/oapi/config.yaml -o ./gen/server/api.go ./docs/api.yaml

mod:
	go mod download

test: mod codegen
	go test -v ./...

clean:
	-rm -r build/ gen/