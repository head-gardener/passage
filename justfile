# starts client and configures its interface
run:
  #!/usr/bin/env sh
  set -ex
  {
    sleep 2
    ip a add 10.1.0.1/24 dev tun1
    ip l set dev tun1 up
  } &
  go run . -config ./example/config.yml

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
  docker run -v ./example/config.yml:/config.yml \
    --cap-add NET_ADMIN --rm -it passage passage-wrapped 10.1.0.1

# builds and loads a docker image with nix
build-docker:
  nix build .#passage-image --print-out-paths | xargs -I _ sh -c "cat _ | docker load"
