#!/bin/bash
set -e

_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OPENSSL_CONFIG="../openssl.cnf"
CERT_DIRECTORY="$_dir/certs"
CRLDP_BASE="http://fake-crl-service:9007/"


if [ -d "$CERT_DIRECTORY" ]
then
  echo "Directory $CERT_DIRECTORY exists." 
  echo "Assuming that certificates are already in place."
  echo "To run this script remove $CERT_DIRECTORY"
  exit 0
fi

echo "Using CRL distribution points: $CRLDP_BASE<ca_name>.crl"

# Setup: build intermediate directories.
mkdir -p "$CERT_DIRECTORY"
cd "$CERT_DIRECTORY"

rm -rf crl/
mkdir -p private
mkdir -p crl


declare -ra CA_NAMES=(
  cbsd_ca
  non_cbrs_root_ca
  non_cbrs_root_signed_cbsd_ca
  non_cbrs_root_signed_oper_ca
  non_cbrs_root_signed_sas_ca
  proxy_ca
  revoked_cbsd_ca
  revoked_proxy_ca
  revoked_sas_ca
  root_ca
  root-ecc_ca
  sas_ca
  sas-ecc_ca
  unrecognized_root_ca
)

# Create an empty index.txt database for each CA.
rm -rf db/
for ca in "${CA_NAMES[@]}"; do
  mkdir -p "db/$ca"
  touch "db/$ca/index.txt"
  echo -n "unique_subject = no" > "db/$ca/index.txt.attr"
done

# Runs openssl using the database for a particular CA.
# $1 should match an entry from the CA_NAMES array.
function openssl_db {
  OPENSSL_CNF_CA_DIR="db/$1" OPENSSL_CNF_CRLDP="$CRLDP_BASE$1.crl" \
      openssl "${@:2}" \
      -config $OPENSSL_CONFIG \
      -cert "$1.cert" -keyfile "private/$1.key"
}

function gen_cbsd_cert {
  # Called with:
  # $1 = device name (For example device_a)
  # $2 = fcc_id (For example test_fcc_id_a)
  # $3 = serial number (For example test_serial_number_a)
  echo "Generating cert $1 for device with fcc_id=$2 sn=$3"
  openssl req -new -newkey rsa:2048 -nodes \
      -config $OPENSSL_CONFIG \
      -out "$1.csr" -keyout "$1.key" \
      -subj "/C=US/O=Wireless Innovation Forum/OU=WInnForum CBSD Certificate/CN=$2:$3"
  echo "Signing cert $1 for device with fcc_id=$2 sn=$3"
  openssl_db cbsd_ca ca -in "$1.csr" \
      -out "$1.cert" \
      -policy policy_anything -extensions "cbsd_req_$1_sign" \
      -config $OPENSSL_CONFIG \
      -batch -notext -create_serial -utf8 -days 1185 -md sha384
}

function gen_corrupt_cert() {
  cp "$1" "$3"
  cp "$2" "$4"
  # cert file have first header line with 28 char: we want to change the
  # 20th cert character
  pos=48
  hex_byte=$(xxd -seek $((10#$pos)) -l 1 -ps "$4" -)
  # Modifying the byte value. If the byte character is 'z' or 'Z' or '9',
  # then it is decremented to 'y' or 'Y' or '8' respectively.
  # If the value is '+' or '/' then we set it to 'A', else the current character
  # value is incremented by 1.
  # This takes care of all the 64 characters of Base64 encoding.
  if [[ $hex_byte == "7a"  ||  $hex_byte == "5a" || $hex_byte == "39" ]]; then
    corrupted_dec_byte=$(($((16#$hex_byte)) -1))
  elif [[ $hex_byte == "2f"  ||  $hex_byte == "2b" ]]; then
    corrupted_dec_byte=65
  else
    corrupted_dec_byte=$(($((16#$hex_byte)) +1))
  fi
  # write it back
  printf "%x: %02x" $pos $corrupted_dec_byte | xxd -r - "$4"
}

# Generate root and intermediate CA certificate/key.
printf "\n\n"
echo "Generate 'root_ca' and 'root-ecc_ca' certificate/key"
openssl req -new -x509 -newkey rsa:4096 -sha384 -nodes -days 7300 \
    -extensions root_ca -config $OPENSSL_CONFIG \
    -out root_ca.cert -keyout private/root_ca.key \
    -subj "/C=US/O=WInnForum/OU=RSA Root CA0001/CN=WInnForum RSA Root CA"

openssl ecparam -genkey -out  private/root-ecc_ca.key -name secp521r1
openssl req -new -x509 -key private/root-ecc_ca.key -out root-ecc_ca.cert \
    -sha384 -nodes -days 7300 -extensions root_ca \
    -config $OPENSSL_CONFIG \
    -subj "/C=US/O=WInnForum/OU=ECC Root CA0001/CN=WInnForum ECC Root CA"

printf "\n\n"
echo "Generate 'sas_ca' and 'sas-ecc_ca' certificate/key"
openssl req -new -newkey rsa:4096 -nodes \
    -config $OPENSSL_CONFIG \
    -out sas_ca.csr -keyout private/sas_ca.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=RSA SAS Provider CA0001/CN=WInnForum RSA SAS Provider CA"
openssl_db root_ca ca -in sas_ca.csr \
    -policy policy_anything -extensions sas_ca_sign \
    -out sas_ca.cert \
    -batch -notext -create_serial -utf8 -days 5475 -md sha384

openssl ecparam -genkey -out  private/sas-ecc_ca.key -name secp521r1
openssl req -new -nodes \
    -config $OPENSSL_CONFIG \
    -out sas-ecc_ca.csr -key private/sas-ecc_ca.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=ECC SAS Provider CA0001/CN=WInnForum ECC SAS Provider CA"
openssl_db root-ecc_ca ca \
    -in sas-ecc_ca.csr -policy policy_anything -extensions sas_ca_sign \
    -out sas-ecc_ca.cert \
    -batch -notext -create_serial -utf8 -days 5475 -md sha384

printf "\n\n"
echo "Generate 'cbsd_ca' certificate/key"
openssl req -new -newkey rsa:4096 -nodes \
    -config $OPENSSL_CONFIG \
    -out cbsd_ca.csr -keyout private/cbsd_ca.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=RSA CBSD OEM CA0001/CN=WInnForum RSA CBSD OEM CA"
openssl_db root_ca ca -in cbsd_ca.csr \
    -policy policy_anything -extensions cbsd_ca_sign \
    -out cbsd_ca.cert \
    -batch -notext -create_serial -utf8 -days 5475 -md sha384

openssl ecparam -genkey -out  private/cbsd-ecc_ca.key -name secp521r1
openssl req -new -nodes \
    -config $OPENSSL_CONFIG \
    -out cbsd-ecc_ca.csr -key private/cbsd-ecc_ca.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=ECC CBSD OEM CA0001/CN=WInnForum ECC CBSD OEM CA"
openssl_db root-ecc_ca ca \
    -in cbsd-ecc_ca.csr -policy policy_anything -extensions cbsd_ca_sign \
    -out cbsd-ecc_ca.cert \
    -batch -notext -create_serial -utf8 -days 5475 -md sha384

printf "\n\n"
echo "Generate 'proxy_ca' certificate/key"
openssl req -new -newkey rsa:4096 -nodes \
    -config $OPENSSL_CONFIG \
    -out proxy_ca.csr -keyout private/proxy_ca.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=RSA Domain Proxy CA0001/CN=WInnForum RSA Domain Proxy CA"
openssl_db root_ca ca -in proxy_ca.csr \
    -policy policy_anything -extensions oper_ca_sign \
    -out proxy_ca.cert \
    -batch -notext -create_serial -utf8 -days 5475 -md sha384

# Generate fake server certificate/key.
printf "\n\n"
echo "Generate 'server' certificate/key"
openssl req -new -newkey rsa:2048 -nodes \
    -config $OPENSSL_CONFIG \
    -out server.csr -keyout server.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=WInnForum SAS Provider Certificate/CN=fake-sas-service"
openssl_db sas_ca ca \
    -in server.csr -out server.cert \
    -policy policy_anything -extensions sas_req_sign \
    -batch -notext -create_serial -utf8 -days 1185 -md sha384

openssl req -new -newkey rsa:2048 -nodes \
    -config $OPENSSL_CONFIG \
    -out domain_proxy_server.csr -keyout domain_proxy_server.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=WInnForum SAS Provider Certificate/CN=domain-proxy"
openssl_db sas_ca ca \
    -in domain_proxy_server.csr -out domain_proxy_server.cert \
    -policy policy_anything -extensions sas_req_sign \
    -batch -notext -create_serial -utf8 -days 1185 -md sha384

openssl ecparam -genkey -out  server-ecc.key -name secp521r1
openssl req -new -nodes \
    -config $OPENSSL_CONFIG \
    -out server-ecc.csr -key server-ecc.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=WInnForum SAS Provider Certificate/CN=fake-sas-service"
openssl_db sas-ecc_ca ca \
    -in server-ecc.csr -out server-ecc.cert \
    -policy policy_anything -extensions sas_req_sign \
    -batch -notext -create_serial -utf8 -days 1185 -md sha384

# Generate sas server certificate/key.
printf "\n\n"
echo "Generate 'sas server' certificate/key"
openssl req -new -newkey rsa:2048 -nodes \
    -config $OPENSSL_CONFIG \
    -out sas.csr -keyout sas.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=WInnForum SAS Provider Certificate/CN=fake-sas-service"
openssl_db sas_ca ca -in sas.csr \
    -out sas.cert \
    -policy policy_anything -extensions sas_req_sign \
    -batch -notext -create_serial -utf8 -days 1185 -md sha384

openssl req -new -newkey rsa:2048 -nodes \
    -config $OPENSSL_CONFIG \
    -out sas_1.csr -keyout sas_1.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=WInnForum SAS Provider Certificate 001/CN=fake-sas-service"
openssl_db sas_ca ca -in sas_1.csr \
    -out sas_1.cert \
    -policy policy_anything -extensions sas_req_sign \
    -batch -notext -create_serial -utf8 -days 1185 -md sha384

# Generate normal operation device certificate/key.
printf "\n\n"
echo "Generate 'certs for devices' certificate/key"
gen_cbsd_cert device_a test_fcc_id_a test_serial_number_a
gen_cbsd_cert device_b test_fcc_id_b test_serial_number_b
gen_cbsd_cert device_c test_fcc_id_c test_serial_number_c
gen_cbsd_cert device_d test_fcc_id_d test_serial_number_d
gen_cbsd_cert device_e test_fcc_id_e test_serial_number_e
gen_cbsd_cert device_f test_fcc_id_f test_serial_number_f
gen_cbsd_cert device_g test_fcc_id_g test_serial_number_g
gen_cbsd_cert device_h test_fcc_id_h test_serial_number_h
gen_cbsd_cert device_i test_fcc_id_i test_serial_number_i
gen_cbsd_cert device_j test_fcc_id_j test_serial_number_j

printf "\n\n"
echo "Generate 'admin' certificate/key"
openssl req -new -newkey rsa:2048 -nodes \
    -config $OPENSSL_CONFIG \
    -out admin.csr -keyout admin.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=WInnForum Admin Certificate/CN=SAS admin Example"
openssl_db sas_ca ca -in admin.csr \
    -out admin.cert \
    -policy policy_anything -extensions cbsd_req_sign \
    -batch -notext -create_serial -utf8 -days 1185 -md sha384

# Generate Domain Proxy certificate/key.
printf "\n\n"
echo "Generate 'domain_proxy' certificate/key"
openssl req -new -newkey rsa:2048 -nodes \
    -config $OPENSSL_CONFIG \
    -out domain_proxy.csr -keyout domain_proxy.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=WInnForum Domain Proxy Certificate/CN=0123456789:0002"
openssl_db proxy_ca ca \
    -in domain_proxy.csr -out domain_proxy.cert \
    -policy policy_anything -extensions oper_req_sign \
    -batch -notext -create_serial -utf8 -days 1185 -md sha384

openssl req -new -newkey rsa:2048 -nodes \
    -config $OPENSSL_CONFIG \
    -out domain_proxy_1.csr -keyout domain_proxy_1.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=WInnForum Domain Proxy Certificate/CN=0123456789:0003"
openssl_db proxy_ca ca \
    -in domain_proxy_1.csr -out domain_proxy_1.cert \
    -policy policy_anything -extensions oper_req_sign \
    -batch -notext -create_serial -utf8 -days 1185 -md sha384

# Generate certificates for test case WINNF.FT.S.SCS.6 -
# Unrecognized root of trust certificate presented during registration
printf "\n\n"
echo "Generate 'unrecognized_device' certificate/key"
openssl req -new -x509 -newkey rsa:4096 -sha384 -nodes -days 7300 \
    -extensions root_ca -config $OPENSSL_CONFIG \
    -out unrecognized_root_ca.cert -keyout private/unrecognized_root_ca.key \
    -subj "/C=US/O=Generic Certification Organization/OU=www.example.org/CN=Generic RSA Root CA"

openssl req -new -newkey rsa:2048 -nodes \
    -config $OPENSSL_CONFIG \
    -out unrecognized_device.csr -keyout unrecognized_device.key \
    -subj "/C=US/O=Generic Certification Organization/OU=www.example.org/CN=Unrecognized CBSD"
openssl_db unrecognized_root_ca ca \
    -in unrecognized_device.csr \
    -out unrecognized_device.cert \
    -policy policy_anything -extensions cbsd_req_sign \
    -batch -notext -create_serial -utf8 -days 1185 -md sha384

# Certificates for test case WINNF.FT.S.SCS.7 -
# corrupted certificate, based on device_a.cert
printf "\n\n"
echo "Generate 'device_corrupted' certificate/key"
gen_corrupt_cert \
    device_a.key device_a.cert device_corrupted.key device_corrupted.cert

# Certificate for test case WINNF.FT.S.SCS.8 -
# Self-signed certificate presented during registration using the same CSR that
# was created for normal operation
printf "\n\n"
echo "Generate 'device_self_signed' certificate/key"
openssl x509 -signkey device_a.key -in device_a.csr \
    -out device_self_signed.cert \
    -req -days 1185

# Certificate for test case WINNF.FT.S.SCS.9 -
# Non-CBRS trust root signed certificate presented during registration
printf "\n\n"
echo "Generate 'non_cbrs_signed_cbsd_ca' certificate/key"
openssl req -new -x509 -newkey rsa:4096 -sha384 -nodes -days 7300 \
    -extensions root_ca -config $OPENSSL_CONFIG \
    -out non_cbrs_root_ca.cert -keyout private/non_cbrs_root_ca.key \
    -subj "/C=US/O=Non CBRS company/OU=www.example.org/CN=Non CBRS Root CA"

printf "\n\n"
echo "Generate 'non_cbrs_signed_cbsd_ca' certificate/key"
openssl req -new -newkey rsa:4096 -nodes \
    -config $OPENSSL_CONFIG \
    -out non_cbrs_root_signed_cbsd_ca.csr \
    -keyout private/non_cbrs_root_signed_cbsd_ca.key \
    -subj "/C=US/O=Non CBRS company/OU=www.example.org/CN=Non CBRS CBSD CA"
openssl_db non_cbrs_root_ca ca \
    -in non_cbrs_root_signed_cbsd_ca.csr \
    -policy policy_anything -extensions cbsd_ca_sign \
    -out non_cbrs_root_signed_cbsd_ca.cert \
    -batch -notext -create_serial -utf8 -days 5475 -md sha384

# Generate CBSD certificate signed by a intermediate CBSD CA which is signed
# by a non-CBRS root CA
openssl req -new -newkey rsa:2048 -nodes \
    -config $OPENSSL_CONFIG \
    -out non_cbrs_signed_device.csr -keyout non_cbrs_signed_device.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=WInnForum CBSD Certificate/CN=test_fcc_id:0001"
openssl_db non_cbrs_root_signed_cbsd_ca ca \
    -in non_cbrs_signed_device.csr \
    -out non_cbrs_signed_device.cert \
    -policy policy_anything -extensions cbsd_req_sign \
    -batch -notext -create_serial -utf8 -days 1185 -md sha384

# Certificate for test case WINNF.FT.S.SCS.10 -
# Certificate of wrong type presented during registration.
# Creating a wrong type certificate by reusing the server.csr and creating
# a server certificate.
printf "\n\n"
echo "Generate wrong type certificate/key"
openssl_db sas_ca ca -in server.csr \
    -out device_wrong_type.cert \
    -policy policy_anything -extensions sas_req_sign \
    -batch -notext -create_serial -utf8 -days 1185 -md sha384

# Certificate for test case WINNF.FT.S.SCS.11
printf "\n\n"
echo "Generate blacklisted client certificate/key"
openssl req -new -newkey rsa:2048 -nodes \
    -config $OPENSSL_CONFIG \
    -out device_blacklisted.csr -keyout device_blacklisted.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=WInnForum CBSD Certificate/CN=test_fcc_id:sn_0001"
openssl_db cbsd_ca ca \
    -in device_blacklisted.csr -out device_blacklisted.cert \
    -policy policy_anything -extensions cbsd_req_sign \
    -batch -notext -create_serial -utf8 -days 1185 -md sha384

# Revoke the device_blacklisted.cert for WINNF.FT.S.SCS.11
openssl_db cbsd_ca ca -revoke device_blacklisted.cert

# Certificate for test case WINNF.FT.S.SDS.11
printf "\n\n"
echo "Generate blacklisted domain proxy certificate/key"
openssl req -new -newkey rsa:2048 -nodes \
    -config $OPENSSL_CONFIG \
    -out domain_proxy_blacklisted.csr -keyout domain_proxy_blacklisted.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=WInnForum Domain Proxy Certificate/CN=0123456789:0004"
openssl_db proxy_ca ca \
    -in domain_proxy_blacklisted.csr \
    -out domain_proxy_blacklisted.cert \
    -policy policy_anything -extensions oper_req_sign \
    -batch -notext -create_serial -utf8 -days 1185 -md sha384

# Revoke the domain_proxy_blacklisted.cert for WINNF.FT.S.SDS.11
openssl_db proxy_ca ca -revoke domain_proxy_blacklisted.cert

# Certificate for test case WINNF.FT.S.SSS.11
printf "\n\n"
echo "Generate blacklisted sas certificate/key"
openssl req -new -newkey rsa:2048 -nodes \
    -config $OPENSSL_CONFIG \
    -out sas_blacklisted.csr -keyout sas_blacklisted.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=WInnForum SAS Provider Certificate/CN=fake-sas-service"
openssl_db sas_ca ca \
    -in sas_blacklisted.csr -out sas_blacklisted.cert \
    -policy policy_anything -extensions sas_req_sign \
    -batch -notext -create_serial -utf8 -days 1185 -md sha384

# Revoke the sas_blacklisted.cert for WINNF.FT.S.SSS.11
openssl_db sas_ca ca -revoke sas_blacklisted.cert

# Certificate for test case WINNF.FT.S.SCS.12 -
# Expired certificate presented during registration
printf "\n\n"
echo "Generate 'device_expired' certificate/key"
openssl req -new -newkey rsa:2048 -nodes \
    -config $OPENSSL_CONFIG \
    -out device_expired.csr -keyout device_expired.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=WInnForum CBSD Certificate/CN=test_fcc_id:0002"
openssl_db cbsd_ca ca \
    -in device_expired.csr -out device_expired.cert \
    -policy policy_anything -extensions cbsd_req_sign \
    -batch -notext -create_serial -utf8 \
    -startdate 20150214120000Z -enddate 20160214120000Z -md sha384

# Certificate for test case WINNF.FT.S.SCS.15 -
# Certificate with inapplicable fields presented during registration
printf "\n\n"
echo "Generate 'inapplicable certificate for WINNF.FT.S.SCS.15' certificate/key"
openssl_db cbsd_ca ca -in device_a.csr \
    -out device_inapplicable.cert \
    -policy policy_anything -extensions cbsd_req_inapplicable_sign \
    -batch -notext -create_serial -utf8 -days 1185 -md sha384

# Generate certificates for test case WINNF.FT.S.SDS.6 -
# Unrecognized root of trust certificate presented during registration
printf "\n\n"
echo "Generate 'unrecognized_domain_proxy' certificate/key"
openssl req -new -newkey rsa:2048 -nodes \
    -config $OPENSSL_CONFIG \
    -out unrecognized_domain_proxy.csr -keyout unrecognized_domain_proxy.key \
    -subj "/C=US/O=Generic Certification Organization/OU=www.example.org/CN=Unrecognized Domain Proxy"
openssl_db unrecognized_root_ca ca \
    -in unrecognized_domain_proxy.csr \
    -out unrecognized_domain_proxy.cert \
    -policy policy_anything -extensions oper_req_sign \
    -batch -notext -create_serial -utf8 -days 1185 -md sha384

# Certificates for test case WINN.FT.S.SDS.7 -
# corrupted certificate, based on domain_proxy.cert
printf "\n\n"
echo "Generate 'domain_proxy_corrupted' certificate/key"
gen_corrupt_cert \
    domain_proxy.key domain_proxy.cert \
    domain_proxy_corrupted.key domain_proxy_corrupted.cert

# Certificate for test case WINNF.FT.S.SDS.8 -
# Self-signed certificate presented during registration.
# Using the same CSR that was created for normal operation
printf "\n\n"
echo "Generate 'domain_proxy_self_signed' certificate/key"
openssl x509 -signkey domain_proxy.key -in domain_proxy.csr \
    -out domain_proxy_self_signed.cert \
    -req -days 1185

# Certificate for test case WINNF.FT.S.SDS.9 -
# Non-CBRS trust root signed certificate presented during registration
printf "\n\n"
echo "Generate non_cbrs_root_signed_oper_ca.csr certificate/key"
openssl req -new -newkey rsa:4096 -nodes \
    -config $OPENSSL_CONFIG \
    -out non_cbrs_root_signed_oper_ca.csr \
    -keyout private/non_cbrs_root_signed_oper_ca.key \
    -subj "/C=US/O=Non CBRS company/OU=www.example.org/CN=Non CBRS Domain Proxy CA"

openssl_db non_cbrs_root_ca ca \
    -in non_cbrs_root_signed_oper_ca.csr \
    -policy policy_anything -extensions oper_ca_sign \
    -out non_cbrs_root_signed_oper_ca.cert \
    -batch -notext -create_serial -utf8 -days 5475 -md sha384

# Generate a Domain Proxy certifcate signed by an intermediate Domain Proxy CA
# which is signed by a non-CBRS root CA
printf "\n\n"
echo "Generate domain_proxy certificate/key"
openssl req -new -newkey rsa:2048 -nodes \
    -config $OPENSSL_CONFIG \
    -out non_cbrs_signed_domain_proxy.csr \
    -keyout non_cbrs_signed_domain_proxy.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=WInnForum Domain Proxy Certificate/CN=0123456789:0005"
openssl_db non_cbrs_root_signed_oper_ca ca \
    -in non_cbrs_signed_domain_proxy.csr \
    -out non_cbrs_signed_domain_proxy.cert \
    -policy policy_anything -extensions oper_req_sign \
    -batch -notext -create_serial -utf8 -days 1185 -md sha384

# Certificate for test case WINNF.FT.S.SDS.10 -
# Certificate of wrong type presented during registration.
# Creating a wrong type certificate by reusing the server.csr and creating
# a server certificate.
printf "\n\n"
echo "Generate wrong type certificate/key"
openssl_db sas_ca ca -in server.csr \
    -out domain_proxy_wrong_type.cert \
    -policy policy_anything -extensions sas_req_sign \
    -batch -notext -create_serial -utf8 -days 1185 -md sha384

# Certificate for test case WINNF.FT.S.SDS.12 -
# Expired certificate presented during registration
printf "\n\n"
echo "Generate 'domain_proxy_expired' certificate/key"
openssl req -new -newkey rsa:2048 -nodes \
    -config $OPENSSL_CONFIG \
    -out domain_proxy_expired.csr -keyout domain_proxy_expired.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=WInnForum Domain Proxy Certificate/CN=0123456789:0006"
openssl_db proxy_ca ca \
    -in domain_proxy_expired.csr -out domain_proxy_expired.cert \
    -policy policy_anything -extensions oper_req_sign \
    -batch -notext -create_serial -utf8 \
    -startdate 20150214120000Z -enddate 20160214120000Z -md sha384

# Certificate for test case WINNF.FT.S.SDS.15 -
# Certificate with inapplicable fields presented during registration
printf "\n\n"
echo "Generate 'inapplicable certificate for WINNF.FT.S.SDS.15' certificate"
openssl_db proxy_ca ca \
    -in domain_proxy.csr -out domain_proxy_inapplicable.cert \
    -policy policy_anything -extensions oper_req_inapplicable_sign \
    -batch -notext -create_serial -utf8 -days 1185 -md sha384

# Generate certificates for test case WINNF.FT.S.SSS.6 -
# Unrecognized root of trust certificate presented during registration
printf "\n\n"
echo "Generate 'unrecognized_sas' certificate/key"
openssl req -new -newkey rsa:2048 -nodes \
    -config $OPENSSL_CONFIG \
    -out unrecognized_sas.csr -keyout unrecognized_sas.key \
    -subj "/C=US/O=Generic Certification Organization/OU=www.example.org/CN=Unrecognized SAS Provider"
openssl_db unrecognized_root_ca ca \
    -in unrecognized_sas.csr \
    -out unrecognized_sas.cert \
    -policy policy_anything -extensions sas_req_sign \
    -batch -notext -create_serial -utf8 -days 1185 -md sha384

# Certificates for test case WINN.FT.S.SSS.7 -
# corrupted certificate, based on sas.cert
printf "\n\n"
echo "Generate 'sas_corrupted' certificate/key"
gen_corrupt_cert sas.key sas.cert sas_corrupted.key sas_corrupted.cert

# Certificate for test case WINNF.FT.S.SSS.8 -
# Self-signed certificate presented during registration.
# Using the same CSR that was created for normal operation
printf "\n\n"
echo "Generate 'sas_self_signed' certificate/key"
openssl x509 -signkey sas.key -in sas.csr \
    -out sas_self_signed.cert \
    -req -days 1185

# Certificate for test case WINNF.FT.S.SSS.9 -
# Non-CBRS trust root signed certificate presented during registration
printf "\n\n"
echo "Generate non_cbrs_root_signed_sas_ca.csr certificate/key"
openssl req -new -newkey rsa:4096 -nodes \
    -config $OPENSSL_CONFIG \
    -out non_cbrs_root_signed_sas_ca.csr \
    -keyout private/non_cbrs_root_signed_sas_ca.key \
    -subj "/C=US/O=Non CBRS company/OU=www.example.org/CN=Non CBRS SAS Provider CA"

openssl_db non_cbrs_root_ca ca \
    -in non_cbrs_root_signed_sas_ca.csr \
    -policy policy_anything -extensions sas_ca_sign \
    -out non_cbrs_root_signed_sas_ca.cert \
    -batch -notext -create_serial -utf8 -days 5475 -md sha384

# Generate a SAS certificate signed by an intermediate SAS CA
# which is signed by a non-CBRS root CA
printf "\n\n"
echo "Generate sas certificate/key"
openssl req -new -newkey rsa:2048 -nodes \
    -config $OPENSSL_CONFIG \
    -out non_cbrs_signed_sas.csr -keyout non_cbrs_signed_sas.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=WInnForum SAS Provider Certificate/CN=non-cbrs.example.org"
openssl_db non_cbrs_root_signed_sas_ca ca \
    -in non_cbrs_signed_sas.csr -out non_cbrs_signed_sas.cert \
    -policy policy_anything -extensions sas_req_sign \
    -batch -notext -create_serial -utf8 -days 1185 -md sha384

# Certificate for test case WINNF.FT.S.SSS.10 -
# Certificate of wrong type presented by SAS Test Harness.
# Creating a wrong type certificate by reusing the device_a.csr and
# creating a client certificate.
printf "\n\n"
echo "Generate wrong type certificate/key"
openssl_db cbsd_ca ca -in device_a.csr \
    -out sas_wrong_type.cert \
    -policy policy_anything -extensions cbsd_req_sign \
    -batch -notext -create_serial -utf8 -days 1185 -md sha384

# Certificate for test case WINNF.FT.S.SSS.12 -
# Expired certificate presented by SAS Test Harness
printf "\n\n"
echo "Generate 'sas_expired' certificate/key"
openssl req -new -newkey rsa:2048 -nodes \
    -config $OPENSSL_CONFIG \
    -out sas_expired.csr -keyout sas_expired.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=WInnForum SAS Provider Certificate/CN=expired.example.org"
openssl_db sas_ca ca -in sas_expired.csr \
    -out sas_expired.cert \
    -policy policy_anything -extensions sas_req_sign \
    -batch -notext -create_serial -utf8 \
    -startdate 20150214120000Z -enddate 20160214120000Z -md sha384

# Certificate for test case WINNF.FT.S.SSS.15 -
# Certificate with inapplicable fields presented by SAS Test Harness
printf "\n\n"
echo "Generate 'inapplicable certificate for WINNF.FT.S.SSS.15' certificate"
openssl_db sas_ca ca -in sas.csr \
    -out sas_inapplicable.cert \
    -policy policy_anything -extensions sas_req_inapplicable_sign \
    -batch -notext -create_serial -utf8 -days 1185 -md sha384

# Certificate for test case WINNF.FT.S.SCS.16 -
# Certificate signed by a revoked CA presented during registration
printf "\n\n"
echo "Generate 'revoked_cbsd_ca' certificate/key to revoke intermediate CA"
openssl req -new -newkey rsa:4096 -nodes \
    -config $OPENSSL_CONFIG \
    -out revoked_cbsd_ca.csr -keyout private/revoked_cbsd_ca.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=RSA CBSD OEM CA0002/CN=WInnForum RSA CBSD OEM CA - Revoked"
openssl_db root_ca ca \
    -in revoked_cbsd_ca.csr \
    -policy policy_anything -extensions cbsd_ca_sign \
    -out revoked_cbsd_ca.cert \
    -batch -notext -create_serial -utf8 -days 5475 -md sha384

# Generate client certificate/key signed by valid CA.
printf "\n\n"
echo "Generate 'client' certificate/key signed by revoked_cbsd_ca"
openssl req -new -newkey rsa:2048 -nodes \
    -config $OPENSSL_CONFIG \
    -out device_cert_from_revoked_ca.csr \
    -keyout device_cert_from_revoked_ca.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=WInnForum CBSD Certificate/CN=test_fcc_id:revoked"
openssl_db revoked_cbsd_ca ca \
    -in device_cert_from_revoked_ca.csr \
    -out device_cert_from_revoked_ca.cert \
    -policy policy_anything -extensions cbsd_req_sign \
    -batch -notext -create_serial -utf8 -days 1185 -md sha384

# Revoke the CBSD CA Certificate.
printf "\n\n"
echo "Revoke intermediate CA (revoked_cbsd_ca) who already signed client certificate."
openssl_db root_ca ca -revoke revoked_cbsd_ca.cert

# Certificate for test case WINNF.FT.S.SDS.16 -
# Certificate signed by a revoked CA presented during registration.
printf "\n\n"
echo "Generate 'revoked_proxy_ca' certificate/key to revoke intermediate CA."
openssl req -new -newkey rsa:4096 -nodes \
    -config $OPENSSL_CONFIG \
    -out revoked_proxy_ca.csr -keyout private/revoked_proxy_ca.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=RSA Domain Proxy CA0002/CN=WInnForum RSA Domain Proxy CA - Revoked"
openssl_db root_ca ca \
    -in revoked_proxy_ca.csr \
    -policy policy_anything -extensions oper_ca_sign \
    -out revoked_proxy_ca.cert \
    -batch -notext -create_serial -utf8 -days 5475 -md sha384

printf "\n\n"
echo "Generate 'domain_proxy' certificate/key signed by revoked_proxy_ca"
openssl req -new -newkey rsa:2048 -nodes \
    -config $OPENSSL_CONFIG \
    -out domain_proxy_cert_from_revoked_ca.csr \
    -keyout domain_proxy_cert_from_revoked_ca.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=WInnForum Domain Proxy Certificate/CN=0123456789:0010 - Revoked"
openssl_db revoked_proxy_ca ca \
    -in domain_proxy_cert_from_revoked_ca.csr \
    -out domain_proxy_cert_from_revoked_ca.cert \
    -policy policy_anything -extensions oper_req_sign \
    -batch -notext -create_serial -utf8 -days 1185 -md sha384

# Revoke the Domain Proxy CA Certificate.
printf "\n\n"
echo "Revoke intermediate CA (revoked_proxy_ca) who already signed domain proxy certificate."
openssl_db root_ca ca -revoke revoked_proxy_ca.cert

# Certificate for test case WINNF.FT.S.SSS.16 -
# Certificate signed by a revoked CA presented during registration
printf "\n\n"
echo "Generate 'sas_ca' certificate/key to revoke intermediate CA."
openssl req -new -newkey rsa:4096 -nodes \
    -config $OPENSSL_CONFIG \
    -out revoked_sas_ca.csr -keyout private/revoked_sas_ca.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=RSA SAS Provider CA0002/CN=WInnForum RSA SAS Provider CA - Revoked"
openssl_db root_ca ca \
    -in revoked_sas_ca.csr \
    -policy policy_anything -extensions sas_ca_sign \
    -out revoked_sas_ca.cert \
    -batch -notext -create_serial -utf8 -days 5475 -md sha384

# Generate sas server certificate/key.
printf "\n\n"
echo "Generate 'sas server' certificate/key signed by revoked_sas_ca"
openssl req -new -newkey rsa:2048 -nodes \
    -config $OPENSSL_CONFIG \
    -out sas_cert_from_revoked_ca.csr -keyout sas_cert_from_revoked_ca.key \
    -subj "/C=US/O=Wireless Innovation Forum/OU=WInnForum SAS Provider Certificate/CN=fake-sas-service - Revoked"
openssl_db revoked_sas_ca ca \
    -in sas_cert_from_revoked_ca.csr \
    -out sas_cert_from_revoked_ca.cert \
    -policy policy_anything -extensions sas_req_sign \
    -batch -notext -create_serial -utf8 -days 1185 -md sha384

# Revoke the SAS CA Certificate.
openssl_db root_ca ca -revoke revoked_sas_ca.cert

# Generate trusted CA bundle.
printf "\n\n"
echo "Generate 'ca' bundle"
cat cbsd_ca.cert proxy_ca.cert sas_ca.cert root_ca.cert cbsd-ecc_ca.cert \
    sas-ecc_ca.cert root-ecc_ca.cert revoked_cbsd_ca.cert \
    revoked_sas_ca.cert revoked_proxy_ca.cert > ca.cert
# Append the github cert, use to download test dump file.
echo | \
    openssl s_client \
        -servername raw.githubusercontent.com \
        -connect raw.githubusercontent.com:443 \
        -prexit < /dev/null 2>&1 | \
    sed -n -e '/BEGIN\ CERTIFICATE/,/END\ CERTIFICATE/ p' >> ca.cert

# Generate a DER-formatted CRL for each CA, per
# https://tools.ietf.org/html/rfc5280#section-4.2.1.13
for ca in "${CA_NAMES[@]}"; do
  pemfile="crl/$ca.crl.pem"
  derfile="crl/$ca.crl"
  printf "\n\n"
  echo "Generating CRL: $pemfile"
  openssl_db "$ca" ca -gencrl -crldays 365 -out "$pemfile"
  echo "Converting $pemfile (PEM) to $derfile (DER)"
  openssl crl -inform pem -outform der -in "$pemfile" -out "$derfile"
done

# Output crl_index files for fake_sas.py.  Real SASes should obtain this
# information from client certificate CRL distribution points.
printf "%s.crl\n" cbsd_ca sas_ca proxy_ca root_ca \
    > crl/crl_index_sxs11.txt
printf "%s.crl\n" revoked_cbsd_ca revoked_sas_ca revoked_proxy_ca root_ca \
    > crl/crl_index_sxs16.txt

# Generate Domain Proxy certificate bundle
cat domain_proxy_server.cert sas_ca.cert root_ca.cert \
    > domain_proxy_bundle.cert

# Cleanup: remove all files not directly used by the testcases or other scripts.
rm ./*.csr
