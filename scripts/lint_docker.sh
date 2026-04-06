CONTAINER_NAME=my_golangci_lint

repo_root=$(git rev-parse --show-toplevel)

docker volume inspect go-mod-cache   >/dev/null 2>&1 || docker volume create go-mod-cache
docker volume inspect golangci-cache >/dev/null 2>&1 || docker volume create golangci-cache
docker build -t $CONTAINER_NAME -f ./golangci-lint.Dockerfile .
docker run --rm -i \
  -e GOMODCACHE=/go/pkg/mod \
  -e XDG_CACHE_HOME=/cache/golangci-lint \
  -e GOCACHE=/cache/go-build \
  -v go-mod-cache:/go/pkg/mod \
  -v golangci-cache:/cache \
  -v "${repo_root}":/app:ro \
  -w /app \
  $CONTAINER_NAME \
  golangci-lint run --timeout=5m --color=always ./...
