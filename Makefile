# 開発用のDockerコンテナの起動を行います。ホットリロードに対応しています。
.PHONY: up
up:
	docker compose up

# 開発用のDockerコンテナのビルドと起動を行います。ホットリロードに対応しています。
.PHONY: up-build
up-build:
	docker compose up --build

# 単体テストを行います。
.PHONY: test
test:
	go test -v ./...
