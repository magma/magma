#!/bin/bash

set -e
dst=$1
src=$2

echo $dst
echo $src

export ASN1C_PREFIX=Ngap_

/home/vagrant/test/asn1c/asn1c/asn1c \
    -pdu=all \
    -fcompound-names \
    -fno-include-deps \
    -findirect-choice \
    -gen-PER \
    -D $dst \
     $src	
    #r16.0/*.asn1


rm -f $dst/converter-example.mk $dst/Makefile.am.asn1convert
mv $dst/converter-example.c $dst/pdu_collection.c example
