#!/bin/bash

# set -x # if you want to print the commands

source .env
source .user

req_url=`echo $CREATE_ACCOUNT_URL/$ID/`
token_header=`echo x-jwt-token: $TOKEN`

echo $req_url;

curl -i \
    -H "$token_header" \
    $req_url \
;