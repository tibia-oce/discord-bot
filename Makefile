build:
	go build -o ./bin/application ./src/main.go

run:
	go run ./src/main.go

test:
	go test ./src/... -v

clean:
	rm -rf ./bin

compose:
	docker-compose -f docker/docker-compose.yml down
	docker-compose -f docker/docker-compose.yml up --build
