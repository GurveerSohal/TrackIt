#! /bin/bash
source .env

curl -i -X POST \
    -d "{\"token\" : \"$TOKEN\", \"workout_number\" : 1, \"reps\" : 10, \"name\" : \"bench\"}" \
    http://localhost:8080/api/create-set