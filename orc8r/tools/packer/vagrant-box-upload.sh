#!/bin/bash

# Requires jq in the path - apt install jq
# References
# https://www.vagrantup.com/docs/vagrant-cloud/api.html#creating-a-usable-box-from-scratch
# https://www.vagrantup.com/docs/vagrant-cloud/boxes/create.html

BOX_FILE=$1

USER=magmacore
BOX=$(basename $BOX_FILE | cut -d_ -f1-2 | cut -d. -f1)
VERSION="1.1.$(date +"%Y%m%d")"
BOX_PROVIDER=virtualbox
if echo $BOX_FILE | grep -q libvirt; then
  BOX_PROVIDER=libvirt
fi

#source vagrant-cloud-token
if [ -z "$VAGRANT_CLOUD_TOKEN" ]; then
  echo "VAGRANT_CLOUD_TOKEN variable is unset. Cannot continue." 1>&2
  exit 1
fi

echo -n "Creating the box for $USER/$BOX..."
curl --header "Content-Type: application/json" --header "Authorization: Bearer $VAGRANT_CLOUD_TOKEN" https://app.vagrantup.com/api/v1/boxes --data '{ "box": { "username": "'$USER'", "name": "'$BOX'" } }'
#{"errors":["Type has already been taken"],"success":false}

echo -en "\nCreate the version $VERSION..."
curl --header "Content-Type: application/json" --header "Authorization: Bearer $VAGRANT_CLOUD_TOKEN" https://app.vagrantup.com/api/v1/box/$USER/$BOX/versions --data '{ "version": { "version": "'$VERSION'" } }'
#{"errors":["Version has already been taken"],"success":false}

echo -en "\nCreating the $provider provider..."
curl --header "Content-Type: application/json" --header "Authorization: Bearer $VAGRANT_CLOUD_TOKEN" https://app.vagrantup.com/api/v1/box/$USER/$BOX/version/$VERSION/providers --data "{ \"provider\": { \"name\": \"$BOX_PROVIDER\" }  }"
#{"errors":["Metadata provider must be unique for version"],"success":false}

echo -en "\nReceiving the upload url..."
response=$(curl -s --header "Authorization: Bearer $VAGRANT_CLOUD_TOKEN" https://app.vagrantup.com/api/v1/box/$USER/$BOX/version/$VERSION/provider/$BOX_PROVIDER/upload)

upload_path=$(echo "$response" | jq .upload_path | tr -d '"')

echo -en "\nUploading ..."
curl $upload_path -o upload.log --progress-bar --request PUT --upload-file $BOX_FILE

# When the upload finishes, you can verify it worked by making this request and matching the hosted_token it returns to the previously retrieved upload token.
uploaded_path=$(curl -s https://app.vagrantup.com/api/v1/box/$USER/$BOX/version/$VERSION/provider/$BOX_PROVIDER | jq .hosted_token | tr -d '"')

if [[ $upload_path == $uploaded_path ]]; then
    echo -en "\nReleasing..."
    curl --header "Authorization: Bearer $VAGRANT_CLOUD_TOKEN" https://app.vagrantup.com/api/v1/box/$USER/$BOX/version/$VERSION/release --request PUT
    echo "Done!"
else
    echo -e "\nUpload did not work: $upload_path != $uploaded_path"
fi
