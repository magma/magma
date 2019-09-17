#!/bin/bash
cd /usr/share/yang/models || exit

pyang -f tree --strict --ietf /validate.yang
