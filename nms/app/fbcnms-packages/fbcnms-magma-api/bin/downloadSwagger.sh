#! /bin/sh

set -e

help()
{
   echo ""
   echo "Usage: $0 -c <file> -k <file> -s"
   echo "    -h : specify hostname, e.g. api-staging.magma.etagecom.io"
   echo "    -c : specify client certificate, e.g. magma.pem"
   echo "    -k : specify client private key, e.g. magma.key.pem"
   echo "Uses env API_HOST, API_CERT_FILENAME and API_PRIVATE_KEY_FILENAME "
   echo "   for host, cert, and key, if defined."
   exit 1
}

while getopts "h:c:k:" opt
do
   case "$opt" in
      h ) host="$OPTARG" ;;
      c ) cert="$OPTARG" ;;
      k ) key="$OPTARG" ;;
      ? ) help ;;
   esac
done

if [ -z "$host" ] && [ -n "$API_HOST" ]
then
    host="$API_HOST"
fi

if [ -z "$cert" ] && [ -n "$API_CERT_FILENAME" ]
then
    cert="$API_CERT_FILENAME"
fi

if [ -z "$key" ] && [ -n "$API_PRIVATE_KEY_FILENAME" ]
then
    key="$API_PRIVATE_KEY_FILENAME"
fi

if [ -z "$host" ] || [ -z "$cert" ] || [ -z "$key" ]
then
    echo "missing host, cert, or key"
    exit 1
fi

url="https://$host/apidocs/v1/swagger.yml"

echo curl "$url" --cert "$cert" --key "$key"
curl "$url" --cert "$cert" --key "$key" > swagger.yml
