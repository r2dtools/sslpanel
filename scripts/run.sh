#!/bin/bash

set -e

echo "Starting container ..."
docker compose -f ./docker-compose.staging.yml up -d --build

CONTAINER_STATUS=$(docker inspect --format "{{json .State.Health.Status }}" sslpanelweb)
until [ $CONTAINER_STATUS == '"healthy"' ]
do
    echo "Waiting for container to start..."
    ((c++)) && ((c==5)) && echo "Could not start container" && (exit 1)
    sleep 10
    CONTAINER_STATUS=$(docker inspect --format "{{json .State.Health.Status }}" sslpanelweb)
done

echo "Adding migrations ..."
docker exec sslpanelback ./cli migrations up

echo "SSLPanel is available at http://localhost:5173"
