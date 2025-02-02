all:
	echo all

build:
	mkdir -p bin
	go build -o bin/bootstrap ./cmd/bootstrap

start-server:
	docker run --name keycloak -d --restart=always \
			-p 8080:8080 -e KC_BOOTSTRAP_ADMIN_USERNAME=admin -e KC_BOOTSTRAP_ADMIN_PASSWORD=admin \
			quay.io/keycloak/keycloak:26.1 \
			start-dev

stop-server:
	docker rm -f keycloak || true
