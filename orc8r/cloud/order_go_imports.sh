#!/bin/bash

# remove all blank lines in go 'imports' statements,
# then sort with goimports

if [ $# != 1 ] ; then
  echo "usage: $0 <filename>"
  exit 1
fi

sed -i '
  /^import/,/)/ {
    /^$/ d
  }
' "$1"

goimports -w -local magma/ "$1"
