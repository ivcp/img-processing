#!/bin/bash
dir="$(dirname "$0")"

ssh -t ubuntu@ec2-3-75-231-58.eu-central-1.compute.amazonaws.com \
 bash "$dir/polls/deploy.sh"
