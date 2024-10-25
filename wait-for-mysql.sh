#!/bin/bash

set -e

echo "Checking if MySQL is up and running..."
until nc -z -v -w30 $DB_HOST $DB_PORT
do
  echo "Waiting for MySQL server at $DB_HOST:$DB_PORT..."
  sleep 5
done

echo "MySQL is up - executing migrations..."

# Run migrations
migrate -path /app/migrations -database "mysql://$DB_USER:$DB_PASSWORD@tcp($DB_HOST:$DB_PORT)/$DB_NAME" up

echo "Starting the application..."
exec "$@"
