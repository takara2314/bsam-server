.PHONY: up
up:
	docker compose -f compose.yaml -p bsam-server-v4 up

.PHONY: up-build
up-build:
	docker compose -f compose.yaml -p bsam-server-v4 up --build
