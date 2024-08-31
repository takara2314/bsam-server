# 開発用のDockerコンテナの起動を行います。ホットリロードに対応しています。
.PHONY: up
up:
	docker compose up

# 開発用のDockerコンテナのビルドと起動を行います。ホットリロードに対応しています。
.PHONY: up-build
up-build:
	docker compose up --build

# Lintを行います。
.PHONY: lint
lint:
	docker compose run --rm lint run --config .golangci.yaml

# 簡易的なLintを行います。
.PHONY: lint-easy
lint-easy:
	docker compose run --rm lint run --config .golangci.easy.yaml

# 単体テストを行います。
.PHONY: test
test:
	go test -v ./...

# E2Eテストを行います。
.PHONY: test_e2e
test_e2e:
	go test -timeout 10s -v ./e2e/...

# E2Eテストを行います。キャッシュを利用しません。
.PHONY: test_e2e-no-cache
test_e2e-no-cache:
	go test -timeout 10s -count=1 -v ./e2e/...

# stg環境の初期化のためのTerraformを init します
.PHONY: tf-stg-init-init
tf-stg-init-init:
	terraform -chdir=terraform/stg/init init

# stg環境の初期化のためのTerraformを plan します
.PHONY: tf-stg-init-plan
tf-stg-init-plan:
	terraform -chdir=terraform/stg/init plan

# stg環境の初期化のためのTerraformを apply します
.PHONY: tf-stg-init-apply
tf-stg-init-apply:
	terraform -chdir=terraform/stg/init apply
