VERSION="$1"
if [ -z "$VERSION" ]; then
  echo "Please provide a specific version"
  exit 0
fi

FULL_VERSION=$(apt-cache madison magma | grep $VERSION |  cut -d "|" -f 2 |  tr -s " ")
FULL_VERSION=${FULL_VERSION%% }
FULL_VERSION=${FULL_VERSION## }


apt update

if [ -z "$FULL_VERSION" ]; then
  echo "$VERSION doesn't exist using latest"
  sudo apt install magma
else
  echo "We found $FULL_VERSION and we're ready to install it"
  sudo apt install magma="$FULL_VERSION"
fi
