VERSION="$1"
MAGMA_VERSION="ci"
OS_VERSION="focal"

if [ -z "$VERSION" ]; then
  echo "Please provide a specific version"
  exit 0
fi

FULL_VERSION=$(apt-cache madison magma | grep $VERSION |  cut -d "|" -f 2 |  tr -s " ")
FULL_VERSION=${FULL_VERSION%% }
FULL_VERSION=${FULL_VERSION## }

sudo rm -rf /etc/apt/sources.list.d/facebookconnectivity_jfrog_io_artifactory_list_dev_focal.list
sudo apt install -y apt-transport-https gnupg2 wget ca-certificates
sudo wget https://artifactory.magmacore.org:443/artifactory/api/gpg/key/public -O /tmp/public
sudo apt-key add /tmp/public
sudo echo "deb https://artifactory.magmacore.org/artifactory/debian-test $OS_VERSION-$MAGMA_VERSION main" >> /etc/apt/sources.list.d/magma.list

sudo apt update
sudo unlink /etc/magma
sudo cp -R magma/lte/gateway/configs/* /etc/magma/


if [ -z "$FULL_VERSION" ]; then
  echo "$VERSION doesn't exist using latest"
  sudo apt install -y -o Dpkg::Options::="--force-confold" --force-yes  magma
else
  echo "We found $FULL_VERSION and we're ready to install it"
  sudo apt install --yes -o Dpkg::Options::="--force-confold"  magma="$FULL_VERSION"
fi

sudo service magma@magmad start 
