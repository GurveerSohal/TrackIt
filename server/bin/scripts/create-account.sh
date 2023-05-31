#!/bin/bash

# set -x # if you want to print the commands

source .env

req_url=`echo $CREATE_ACCOUNT_URL`

curl -i \
    $req_url \
    -X POST \
    -H "Content-Type: application/json" \
    -d @dummy-account.json