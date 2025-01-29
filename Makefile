build:
	docker build -t forum -f Dockerfile .
run-img:
	docker run --name=forum -p 8081:8080 -d forum
run:
	go run cmd/main.go
stop:
	docker stop forum
init:
	mkdir -p certs

	openssl req -x509 -newkey rsa:4096 -keyout certs/localhost.key -out certs/localhost.crt -days 365 -nodes -subj "/CN=localhost"
	echo 'export SSL_CERT_PATH="certs/localhost.crt"' >> ~/.bashrc
	echo 'export SSL_KEY_PATH="certs/localhost.key"' >> ~/.bashrc
	@echo "✅✅✅"
