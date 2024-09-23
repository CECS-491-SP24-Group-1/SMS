#!/bin/bash
docker build -f ./dockerfiles/app.Dockerfile -t sms:v1 --progress=plain .
#docker compose up
#docker run -it --rm sms:v1