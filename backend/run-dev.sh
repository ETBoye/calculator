
if [ -z "${POSTGRES_USER}" ]; then
    echo "You need to set the POSTGRES_USER environment variable to run db"
    exit
fi

if [ -z "${POSTGRES_PASSWORD}" ]; then
    echo "You need to set the POSTGRES_PASSWORD environment variable to run db"
    exit
fi

go run main.go