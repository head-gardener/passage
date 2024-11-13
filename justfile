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
  exec go run ./cmd/passage -config ./examples/node1.yml {{ args }}

# starts two clients, one local and one in docker
run-pair: build-docker
  #!/usr/bin/env sh
  set -ex
  docker network create --subnet=172.20.0.0/24 passage_test
  trap 'docker network rm passage_test; exit' EXIT
  just run-docker "" "--name passage_test --net passage_test --ip 172.20.0.2 -d"
  trap 'docker container stop passage_test; docker network rm passage_test; exit' EXIT
  just run "-config examples/node1.yml"

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
    docker compose -f "$c" logs &
    sleep 2
    cmd="$(sed -nr 's/#! (.*)/\1/p' "$c")"
    docker compose -f "$c" exec $cmd
    docker compose -f "$c" down
  done

# runs docker image with defaults for testing
run-docker args="" docker_args="-it": build-docker
  docker run -v ./examples/node2.yml:/config.yml \
    --cap-add NET_ADMIN --rm {{ docker_args }} \
    passage passage-wrapped 10.1.0.2 {{ args }}

# builds and loads a docker image with nix
build-docker:
  #!/usr/bin/env sh
  set -exo pipefail
  nix build .#passage-image --print-out-paths | xargs -I _ sh -c "cat _ | docker load"
