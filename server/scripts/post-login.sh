#! /bin/bash

curl -i -X POST \
    -d "{\"username\" : \"user1\", \"password\" : \"pwd1\"}" \
    http://localhost:8080/api/login