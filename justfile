default:
  @just --list

# starts client and configures its interface
run:
  #!/usr/bin/env sh
  set -ex
  sudo w > /dev/null
  {
    sleep 2
    sudo ip a add 10.1.0.1/24 dev tun1
    sudo ip l set dev tun1 up
  } &
  sudo -E go run ./cmd/passage -config ./examples/config.yml

# formats and checks
pre-commit: format test check

# runs go tests
test:
  go test ./...

# formats whole tree with treefmt
format:
  treefmt

# repeatedly listens for logs from docker compose
watch-logs:
  while true; do sleep 2; docker compose logs -f; done

# starts docker containers and checks their connectivity
check: build-docker
  #!/usr/bin/env sh
  set -ex
  trap 'docker compose down' EXIT SIGINT
  docker compose up --build -d
  sleep 2
  docker compose exec node1 ping -c 4 10.1.0.2

# runs docker image with defaults for testing
run-docker: build-docker
  docker run -v ./examples/config.yml:/config.yml \
    --cap-add NET_ADMIN --rm -it passage passage-wrapped 10.1.0.1

# builds and loads a docker image with nix
build-docker:
  nix build .#passage-image --print-out-paths | xargs -I _ sh -c "cat _ | docker load"
