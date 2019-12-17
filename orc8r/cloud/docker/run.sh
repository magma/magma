#! /bin/bash

if [ "$1" == "-all" ]; then
    docker-compose -f docker-compose.yml -f docker-compose.override.yml -f docker-compose.metrics.yml up -d
else
    docker-compose up -d
fi
