if [ -z "${POSTGRES_USER}" ]; then
    echo "You need to set the POSTGRES_USER environment variable to run the application"
    exit
fi

if [ -z "${POSTGRES_PASSWORD}" ]; then
    echo "You need to set the POSTGRES_PASSWORD environment variable to run application"
    exit
fi

docker-compose down
docker-compose up