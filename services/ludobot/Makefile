.PHONY: bot

bot: mod
	@go build -o build/ludobot cmd/ludobot/main.go

mod:
	@go mod download

test:
	@go test -v ./...

clean:
	@rm -r build/