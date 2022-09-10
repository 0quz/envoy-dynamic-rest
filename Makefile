server:
	go run *.go

build:
	go build -o bin/server *.go

d.up:
	docker-compose up 

d.down:
	docker-compose down

d.up.build:
	docker-compose --build up