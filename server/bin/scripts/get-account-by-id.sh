#!/bin/bash

# set -x # if you want to print the commands

source .env

if [ "#$1" = "#" ]; then 
    echo "Provide a uuid";
    exit 1;
fi

req_url=`echo $CREATE_ACCOUNT_URL$1/`

echo $req_url;

curl -i \
    $req_url \
;