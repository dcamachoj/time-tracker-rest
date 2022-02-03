#!/usr/bin/env bash

alias docker='podman'
alias db-logs='docker logs db'

function db-start() {
  docker run -d \
    --name db \
    -p 3306:3306 \
    --env-file .env \
    --rm \
    mysql:8
}

function db-stop() {
  docker stop db
}

function db-prompt() {
  podman exec -it db mysql -uroot -p'db-root'
}

echo "Loaded functions:"
echo "db-start : Starts the DB"
echo "db-stop  : Stops the DB"
echo "db-logs  : Show the DB logs"
echo "db-prompt: Executes the DB shell prompt"