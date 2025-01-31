docker-build:
	docker build -t my-forum-app .

docker-run:
	docker run -p 8080:8080 my-forum-app

docker-up: 

run:
	go run ./cmd

init:
	mkdir -p certs
    openssl req -x509 -newkey rsa:4096 -keyout certs/localhost.key -out certs/localhost.crt -days 365 -nodes -subj "/CN=localhost"
	echo 'export SSL_CERT_PATH="certs/localhost.crt"' >> ~/.bashrc
	echo 'export SSL_KEY_PATH="certs/localhost.key"' >> ~/.bashrc