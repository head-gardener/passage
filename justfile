default:
  @just --list

# starts client and configures its interface
run args = "":
  sudo -E just run-no-sudo

# same as run but doesn't ask for sudo
run-no-sudo args = "":
  exec go run ./cmd/passage -config ./examples/node1.yml {{ args }}

# formats and checks
pre-commit: format test check-packaging check

# runs go tests
test:
  go test ./...

# runs go property-based tests
test-props:
  seq $(nproc) | xargs -P $(nproc) -I _ go test ./... -run Prop -quickchecks 10000

# formats whole tree with treefmt
format:
  treefmt

# verfies that everything is packaged correctly and starts
check-packaging:
  #!/usr/bin/env sh
  set -exo pipefail
  nix flake check --no-build
  just run-no-sudo -help
  sleep 2
  nix build .#passage --print-out-paths | xargs -I _ sh -c "_/bin/passage -help"
  just run-docker -help ""

# repeatedly listens for logs from docker compose
watch-logs:
  while true; do sleep 2; docker compose logs -f; done

# runs compose checks
check path = "*": build-docker
  #!/usr/bin/env sh
  set -ex
  checks="$(find test -type f -name "*{{ path }}*.yml")"
  echo "--- checks: $checks ---"
  for c in "$checks"; do
    docker compose -f "$c" up --build -d
    docker compose -f "$c" logs -f &
    trap 'echo "$checks" | xargs -I _ docker compose -f _ down; exit' EXIT
    sleep 2
    cmd="$(sed -nr 's/#! (.*)/\1/p' "$c")"
    c="$c" sh -c "set -exo pipefail; $cmd"
    docker compose -f "$c" down
  done

# runs docker image with defaults for testing
run-docker args="" docker_args="-it": build-docker
  #!/usr/bin/env sh
  set -ex
  docker network create --subnet=172.255.0.0/24 passage_test
  trap 'docker network rm passage_test; exit' EXIT
  docker run -v ./examples/node2.yml:/config.yml \
    --name passage_test --net passage_test --ip 172.255.0.2 \
    --cap-add NET_ADMIN --rm {{ docker_args }} \
    passage passage-wrapped {{ args }}

# builds and loads a docker image with nix
build-docker:
  #!/usr/bin/env sh
  set -exo pipefail
  nix build .#passage-image --print-out-paths | xargs -I _ sh -c "cat _ | docker load"
