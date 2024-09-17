#!/bin/bash
TPL="config.template"
DST="/root/.ssh/config"

mkdir -p "/root/.ssh" && envsubst < "/$TPL" > "$DST"

eval "$(ssh-agent)"
ssh-add /root/.ssh/id_rsa
ssh -fN -f -L 60055:localhost:60055 -J jumphost agw -p "$REMOTE_PORT"
#TODO: discuss on which port traffic should be send from enodebd
ssh -fN -f -R 12345:localhost:22 -J jumphost agw -p "$REMOTE_PORT"
sleep infinity
