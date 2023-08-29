#! /bin/bash

curl -i -X POST \
    -d "{\"username\" : \"user2\", \"password\" : \"pwd2\"}" \
    http://localhost:8080/api/login