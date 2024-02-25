#!/bin/bash

SECONDS=0

msg () {
    echo -e "\n******* $1 *******\n"
}

dir="$(dirname "$0")"
source "$dir/.env"

cmd=""

if [[ $SERVER_ENV == "production" ]]; then 
    msg "Pulling from github"
    git pull
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


