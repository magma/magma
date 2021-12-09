VERSION="$1"
if [ -z "$VERSION" ]; then
  echo "Please provide a specific version"
  exit 0
fi

FULL_VERSION=$(apt-cache madison magma | grep $VERSION |  cut -d "|" -f 2 |  tr -s " ")
FULL_VERSION=${FULL_VERSION%% }
FULL_VERSION=${FULL_VERSION## }

sudo apt update

sudo rm -rf /etc/magma

if [ -z "$FULL_VERSION" ]; then
  echo "$VERSION doesn't exist using latest"
  sudo apt install -y -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold" magma
else
  echo "We found $FULL_VERSION and we're ready to install it"
  sudo apt install --yes -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold" magma="$FULL_VERSION"
fi
