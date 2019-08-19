#!/bin/sh
/usr/bin/envsubst < ./radius.config.json.template > ./radius.config.json
./radius
