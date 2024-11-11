default:
  @just --list

# starts client and configures its interface
run args = "":
  sudo -E just run-no-sudo

# same as run but doesn't ask for sudo
run-no-sudo args = "":
  #!/usr/bin/env sh
  set -ex
  {
    sleep 2
    pgrep "passage" > /dev/null
    ip a add 10.1.0.1/24 dev tun1
    ip l set dev tun1 up
  } &
  go run ./cmd/passage -config ./examples/config.yml {{ args }}

# formats and checks
pre-commit: format test check-packaging check

# runs go tests
test:
  go test ./...

# formats whole tree with treefmt
format:
  treefmt

# verfies that everything is packaged correctly and starts
check-packaging:
  nix flake check --no-build
  just run-no-sudo -help
  sleep 2
  nix build .#passage --print-out-paths | xargs -I _ sh -c "_/bin/passage -help"
  just run-docker -help ""

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
run-docker args="" interactive="-it": build-docker
  docker run -v ./examples/config.yml:/config.yml \
    --cap-add NET_ADMIN --rm {{ interactive }} \
    passage passage-wrapped 10.1.0.1 {{ args }}

# builds and loads a docker image with nix
build-docker:
  nix build .#passage-image --print-out-paths | xargs -I _ sh -c "cat _ | docker load"
