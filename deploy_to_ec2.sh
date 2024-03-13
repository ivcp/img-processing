#!/bin/bash
dir="$(dirname "$0")"

ssh -t ubuntu@polls.ovh \
 bash "$dir/polls/build.sh"
