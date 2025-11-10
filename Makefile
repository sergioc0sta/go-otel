build-prod-image:
	docker build -t go-otel -f Dockerfile .	

build-dev-image:
	docker build -t go-otel -f Dockerfile.dev .	

run-container:
	docker run -d -p 8080:8080 --name go-otel-app go-otel

run-access-container:
	docker exec -it go-otel-app sh

run-tests:
	go test ./... 

run-compose-up:
	docker-compose up -d
