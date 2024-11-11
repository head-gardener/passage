# starts client and configures its interface
run:
  #! /usr/bin/env sh
  set -ex
  go run . -config ./example/config.yml &
  sleep 2
  ip a add 10.1.0.1/24 dev tun1
  ip l set dev tun1 up
  wait %1

# starts docker containers and checks their connectivity
check:
  docker compose up --build -d
  sleep 2
  docker compose exec node1 ping -c 4 10.1.0.2
  docker compose down

# init sequence for docker containers
docker-init addr:
  #! /usr/bin/env sh
  set -ex
  mkdir -p /dev/net && mknod /dev/net/tun c 10 200
  ./passage -log-level debug -device-name tun1 -listener-addr "0.0.0.0:8888" &
  sleep 2
  ip a add {{ addr }}/24 dev tun1
  ip l set dev tun1 up
  wait %1
