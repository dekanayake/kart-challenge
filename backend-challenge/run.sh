#!/usr/bin/env bash
set -e

# Environment variables
export PORT=8080
export LOG_LEVEL=debug
export ENVIRONMENT=development
export COUPON_CODE_FOLDER_PATH=/Users/duminda/resume/recoded/code-task/kart-challenge/assets/sort # point to the coupon code file folder location. Make sure the contents in the files are sorted
export GIN_MODE=release # to check debug logs in http calls  change to debug


./backend-challenge
