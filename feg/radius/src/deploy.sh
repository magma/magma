#!/bin/bash

# test server ssh
SSH_CERT_FILE=${SSH_CERT_FILE:-~/.ssh/aws2.pem}
IP_ADDRESS=${IP_ADDRESS:-54.145.218.49}
USERNAME=${USERNAME:-centos}

# server runtime configfile relative to server executable
RADIUS_ACCESS_TOKEN=${RADIUS_ACCESS_TOKEN:-$1}
CONFIG_TEMPLATE_FILE=${CONFIG_TEMPLATE_FILE:-./config/samples/radius.xwfv3.config.json.template}
CONFIG_FILE=${CONFIG_FILE:-./config/samples/radius.xwfv3.config.json}

rm ./radius ./radius.zip

echo creating and copying zip files ...
zip -q radius.zip -r .
zip -q libradius.zip -r ../../lib/go/radius/
scp -i "$SSH_CERT_FILE" "./radius.zip" "$USERNAME@$IP_ADDRESS:~/radius.zip"
scp -i "$SSH_CERT_FILE" "./libradius.zip" "$USERNAME@$IP_ADDRESS:~/libradius.zip"
rm -f radius.zip libradius.zip

ssh "$USERNAME@$IP_ADDRESS" -i "$SSH_CERT_FILE" << "EOF"
	echo extracting radius server ...
	mkdir -p ~/radius/server/
	unzip -o -q ~/radius.zip -d ~/radius/server

	echo extracting libradius ...
	mkdir -p ~/lib/go/radius/
	unzip -o -q ~/libradius.zip -d ~

	echo building ...
	cd ~/radius/server/
	/usr/local/go/bin/go build .

	echo generating config ...
	export RADIUS_ACCESS_TOKEN=$RADIUS_ACCESS_TOKEN
	envsubst < $CONFIG_TEMPLATE_FILE > $CONFIG_FILE

	echo running ...
	./radius -config $CONFIG_FILE
EOF


