#!/bin/sh

set -e

if command -v docker-compose 2>&1 >/dev/null;then
    DOCKER_COMPOSE_COMMAND="docker-compose"
else
    DOCKER_COMPOSE_COMMAND="docker compose"
fi

if [ -z "${POSTGRES_USER}" ]; then
    echo "You need to set the POSTGRES_USER environment variable to run the application"
    exit
fi

if [ -z "${POSTGRES_PASSWORD}" ]; then
    echo "You need to set the POSTGRES_PASSWORD environment variable to run application"
    exit
fi


if [ ! -e nginx/.htpasswd ]
then
    echo "No file nginx/.htpasswd exist. You can copy .htpasswd.unsafe-admin-admin for a user with username: admin and password: admin for testing"
    exit
fi

$DOCKER_COMPOSE_COMMAND down
$DOCKER_COMPOSE_COMMAND up -d