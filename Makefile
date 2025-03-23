build:
	go build -o bin/app cmd/main.go 

run: build
	./bin/app

run-docker:
	docker run -d \
	  --name postgres_container \
	  -e POSTGRES_USER=postgres \
	  -e POSTGRES_PASSWORD=T0psecret \
	  -e POSTGRES_DB=postgres \
	  -p 5432:5432 \
	  postgres:16
	
