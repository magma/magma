#!/bin/bash
# Needs export AWS_SECRET_ACCESS_KEY and AWS_ACCESS_KEY_ID
packer build -force \
-var "subnet=subnet-38f8705f" \
-var "vpc=vpc-ad28ceca" \
$1
