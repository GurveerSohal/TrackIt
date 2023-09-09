#! /bin/bash

curl -i -X POST \
    -d "{\"username\" : \"someuser\", \"password\" : \"unsafepwd\"}" \
    http://localhost:8080/api/signup/