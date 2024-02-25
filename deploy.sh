#!/bin/bash

SECONDS=0

msg () {
    echo -e "\n******* $1 *******\n"
}

source .env
cmd=""

if [[ $SERVER_ENV == "production" ]]; then 
    # git pull
    echo "hi"
    cmd="sudo"
fi


msg "Building image"

$cmd docker compose build

msg "Stopping containers"

$cmd docker compose down

msg "Starting containers"

$cmd docker compose up -d

msg "Removing stale images"

$cmd docker image prune -f


msg "Finished in $SECONDS seconds"

echo "Press Enter to exit"
read


