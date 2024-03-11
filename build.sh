#!/bin/bash

SECONDS=0

msg () {
    echo -e "\n******* $1 *******\n"
}

if [ -d polls/ ]; then
    cd polls/
fi

source .env
cmd=""

if [[ $SERVER_ENV == "production" ]]; then 
    msg "Pulling from github"    
    git pull
    cmd="sudo"
fi


msg "Building image"

$cmd docker compose build

msg "Stopping containers"

$cmd docker compose down --remove-orphans

msg "Starting containers"

$cmd docker compose up -d

msg "Removing stale images"

$cmd docker image prune -f

msg "Finished in $SECONDS seconds"



