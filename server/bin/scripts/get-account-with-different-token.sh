#!/bin/bash

# set -x # if you want to print the commands

source .env
source .user

req_url=`echo $CREATE_ACCOUNT_URL/$ID/`
token_header=`echo x-jwt-token: $TOKEN`

response=$(curl -s -w "%{http_code}" \
    -H "$token_header" \
    $req_url \
;)
http_code=$(tail -n1 <<< "$response")

if [[ $http_code != 200 ]] ; then
    echo "Getting account with valid token failed"
    exit 1
fi

req_url=`echo $CREATE_ACCOUNT_URL/$OTHER_ID/`

echo "http code with valid token $http_code"

response=$(curl -s -w "%{http_code}" \
    -H "$token_header" \
    $req_url \
;)
http_code=$(tail -n1 <<< "$response")

if [[ $http_code != 403 ]] ; then
    echo "Was able to get account with an invalid token"
    exit 1
fi

echo "http code with other account $http_code"