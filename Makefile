build:
	go build -o ./bin/application ./src/main.go

run:
	go run ./src/main.go

test:
	go test ./...

clean:
	rm -rf ./bin

compose:
	docker-compose -f docker/docker-compose.yml down
	docker-compose -f docker/docker-compose.yml up --build

fmt:
	gofmt -s -w .

# go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
# export PATH=$PATH:$(go env GOPATH)/bin
# source ~/.zshrc ----OR---- source ~/.bashrc
lint:
	golangci-lint run
