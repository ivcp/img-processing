#!/bin/bash
dir="$(dirname "$0")"

ssh -t ubuntu@api.polls.ovh \
 bash "$dir/polls/build.sh"
