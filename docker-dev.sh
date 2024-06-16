#!/bin/bash

# Run the application with hot-reloading
# set `DATABASE_HOST=event-mgmt-postgres` in your .env to use 

docker stop event-mgmt-core
docker rm event-mgmt-core
docker rmi event-mgmt-core
docker build --tag event-mgmt-core .
docker run --name event-mgmt-core --env-file .env --network eventmanagementcore_event-mgmt-network --detach --publish 127.0.0.1:8081:8081 --mount type=bind,source="$(pwd)",target=/app event-mgmt-core