#! /bin/bash
source .env

curl -i -X POST \
    -d "{\"token\" : \"dfddkfdj\"}" \
    http://localhost:8080/api/token/verify