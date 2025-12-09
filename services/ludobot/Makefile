.PHONY: bot

bot: mod sql
	@go build -o build/ludobot cmd/ludobot/main.go

mod:
	@go mod download

sql:
	@go run cmd/sqlgen/main.go

test:
	@go test -v ./...

clean:
	@rm -r build/