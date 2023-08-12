#! /bin/bash

curl -i -X POST \
    -d "{\"username\" : \"testAccoun\", \"password\" : \"testP*ssw0rd!\"}" \
    http://localhost:8080/api/login