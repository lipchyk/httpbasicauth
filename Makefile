PHONY: test

test:
	go test -cover -coverprofile=c.out ./...
