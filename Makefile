up:
	docker-compose up --build

down:
	docker-compose down

fmt:
	go fmt ./...

test:
	go vet ./...
	go test -v ./... -count=1