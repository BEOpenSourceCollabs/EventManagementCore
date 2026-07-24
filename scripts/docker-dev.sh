#!/bin/bash

cd $(dirname ${BASH_SOURCE[0]})/../

# Run the application with hot-reloading
# set `DATABASE_HOST=event-mgmt-postgres` in your .env to use 

ENVFILE=".env"

if [ -f .env.$1 ]; then
  ENVFILE=".env.$1"
fi

if [ ! -f $ENVFILE ]; then
  echo "No .env file found. Please create one or specify an environment."
  exit 1
fi

docker stop event-mgmt-core
docker rm event-mgmt-core
docker rmi event-mgmt-core
docker build --tag event-mgmt-core .
docker run --name event-mgmt-core --env-file $ENVFILE --detach --publish 127.0.0.1:8081:8081 --mount type=bind,source="$(pwd)",target=/app event-mgmt-core