#!/bin/bash
/usr/bin/envsubst < ./radius.config.json.template > ./radius.config.json
/usr/bin/envsubst < ./lb.config.json.template > ./lb.config.json