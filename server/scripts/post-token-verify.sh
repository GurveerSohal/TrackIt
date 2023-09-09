#! /bin/bash
source .env

curl -i -X POST \
    -d "{\"token\" : \"$TOKEN\"}" \
    http://localhost:8080/api/token/verify/