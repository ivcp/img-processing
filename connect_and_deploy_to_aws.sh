#!/bin/bash


ssh -t ubuntu@ec2-3-75-231-58.eu-central-1.compute.amazonaws.com \
 cd polls \
 bash deploy.sh
