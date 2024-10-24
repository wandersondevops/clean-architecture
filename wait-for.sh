#!/bin/sh

# wait-for.sh host:port

set -e

host="$1"
shift
cmd="$@"

until nc -z "$host" 5672; do
  >&2 echo "RabbitMQ is unavailable - sleeping"
  sleep 3
done

>&2 echo "RabbitMQ is up - executing command"
exec $cmd
