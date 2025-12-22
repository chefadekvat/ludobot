.PHONY: bot

bot: mod sql
	@go build -o build/ludobot cmd/main.go

mod:
	@go mod download

sql:
	@go run ../../tools/sqlgen/main.go ./sql ./gen/sql

test:
	@go test -v ./...

clean:
	@rm -r build/