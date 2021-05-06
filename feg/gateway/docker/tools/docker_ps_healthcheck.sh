#!/bin/bash

# add the path to the command if docker-compose.yml is not in the same folder
cd "$1" || exit 1

# check if docker services are really running or not
sleep 3
if ! (docker-compose ps | grep "...   Up"); then
  echo "--- ERROR: Docker services are not running ---"
  exit 2
fi

echo "--- Waiting 10 seconds... --- "
sleep 10
printf "\n--- Checking docker services --- \n"

for _ in {1..10}
do
  if (docker-compose ps | grep -q Restarting); then
    printf "\n--- ERROR: Docker services are restarting ---\n"
    docker-compose ps
    exit 3
  else
    printf "."
  fi
done

printf "\n--- Docker services are running ---\n"
