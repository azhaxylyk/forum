build:
	docker build -t forum -f Dockerfile .
run-img:
	docker run --name=forum -p 8081:8080 -d forum
run:
	go run cmd/main.go
stop:
	docker stop forum