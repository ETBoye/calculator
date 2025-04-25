#!/bin/sh
if command -v docker-compose 2>&1 >/dev/null;then
    DOCKER_COMPOSE_COMMAND="docker-compose"
else
    DOCKER_COMPOSE_COMMAND="docker compose"
fi


if [ -z "${POSTGRES_USER}" ]; then
    echo "You need to set the POSTGRES_USER environment variable to run db"
    exit
fi

if [ -z "${POSTGRES_PASSWORD}" ]; then
    echo "You need to set the POSTGRES_PASSWORD environment variable to run db"
    exit
fi


cd ..
docker-compose down
docker-compose up db