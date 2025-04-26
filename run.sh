#!/bin/sh

set -e


case $APP_ENV in
    prod)
        GIN_MODE=release
        ;;
    test)
        GIN_MODE=debug
        ;;
    *)
        echo "Please set environment variable APP_ENV to either test or prod"
        exit
        ;;
esac


# Some people have different commands for docker compose
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

$DOCKER_COMPOSE_COMMAND build
$DOCKER_COMPOSE_COMMAND down

case $APP_ENV in
    prod)
        GIN_MODE=$GIN_MODE $DOCKER_COMPOSE_COMMAND up -d
        ;;
    test)
        GIN_MODE=$GIN_MODE $DOCKER_COMPOSE_COMMAND up
        ;;
 
esac