compose: env
	docker-compose -f docker/docker-compose.yml down -v
	docker volume prune -f
	docker-compose -f docker/docker-compose.yml up --build -d
	docker logs bot -f

build:
	go build -o ./bin/application ./src/main.go

run:
	go run ./src/main.go

test:
	go test ./...

clean:
	rm -rf ./bin
	docker-compose -f docker/docker-compose.yml down -v
	docker system prune -a --volumes -f

fmt:
	gofmt -s -w .

# go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
# export PATH=$PATH:$(go env GOPATH)/bin
# source ~/.zshrc ----OR---- source ~/.bashrc
lint:
	golangci-lint run

env:
	@ if [ ! -f .env ]; then cp .env.example .env; fi
