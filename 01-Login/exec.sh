#!/usr/bin/env -S bash -ex
docker build -t auth0-golang-web-app .
mkdir -p /tmp/store
docker run --env-file .env -p 4242:3000 -v /tmp/store/:/gotem -it auth0-golang-web-app 
