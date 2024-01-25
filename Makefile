run:
	go run -race ./cmd/ghost/*.go

test:
	go test -v --cover ./...

build:
	go build ./cmd/ghost

release:
	@echo "Enter the release version (format vx.x.x).."; \
	read VERSION; \
	git tag -a $$VERSION -m "Releasing "$$VERSION; \
	git push origin $$VERSION

buildlin:
	env GOOS=linux GOARCH=amd64 go build -o ghost-linux ./cmd/ghost
