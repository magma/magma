#! /bin/sh

set -e # exit on any error

SCRIPTDIR="$(dirname "$(realpath "$0")")"

if ! "$SCRIPTDIR/downloadSwagger.sh";
then
    HGROOT=$(hg root)
    DOCKERPATH="$HGROOT/fbcode/fbc/symphony/integration"
    if [ ! -d "$DOCKERPATH" ]; then
        echo "Directory $DOCKERPATH not found; is fbcode checked out?"
    fi
    echo
    echo "Please run via: "
    echo "  cd $DOCKERPATH && docker-compose exec platform-server yarn --cwd /app/packages/fbcnms-magma-api gen"
    exit 1
fi

"$SCRIPTDIR/generateAPIFromSwagger.sh"
