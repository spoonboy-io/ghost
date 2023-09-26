run:
	go run -race ./cmd/ghost/*.go

test:
	go test -v --cover ./...