#!/bin/bash
export $(grep -v '^#' env.properties | xargs)

# Start the service in development mode
docker-compose -f ./docker-compose.dev.yaml up "$@"
