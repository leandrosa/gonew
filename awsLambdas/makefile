localstack-up:
	docker compose -f docker/docker-compose-1.yaml up -d
	echo "access https://app.localstack.cloud"

localstack-down:
	docker compose -f docker/docker-compose-1.yaml down

build:
	./scripts/build.sh

publish:
	./scripts/publish.sh hello-go.zip

publish-delete:
	./scripts/publish.sh hello-go.zip --delete

publish-all:
	./scripts/publish.sh --all