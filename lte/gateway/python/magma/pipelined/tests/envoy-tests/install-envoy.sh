sudo apt-get install -y \
	apt-transport-https \
	ca-certificates \
	curl \
	gnupg2 \
	software-properties-common

sudo curl -sL 'https://getenvoy.io/gpg' | sudo apt-key add -

sudo add-apt-repository "deb [arch=amd64] https://dl.bintray.com/tetrate/getenvoy-deb \
	$(lsb_release -cs) \
	stable"

sudo apt-get update && sudo apt-get install -y getenvoy-envoy

sudo mkdir /etc/envoy/
# sudo nano /etc/envoy/envoy.yaml

sudo nano  /etc/systemd/system/envoy.service


