---
id: version-1.0.0-packer
title: Packer Build
hide_title: true
original_id: packer
---
# Packer Build
## Intro
This directory contains the needed files to create a new amazon ami and vagrant
boxes to run the magma dev and test instances on. You will need to use this if
you want to update the base environment, for example, changing the debian or
kernel version.

If you're looking to install additional software specific to each box when
provisioning, you'll want to instead add that to the Ansible playbook.

## Usage
You'll need to have packer installed. You can get it here:
https://www.packer.io/downloads.html or through you package manager. Packer
works by creating an amazon instance or virtualbox, running some provisioning
scripts, and then saving it as an ami/box.

The .json file defines the base image, and is later provisioned by the shell
scripts.

### Vagrant
To upload to vagrant cloud, you'll need the upload token from lastpass.
Once you've made your changes, run
```
export ATLAS_TOKEN=<token_form_lastpass>
packer validate debian-stretch-virtualbox.json
```
and then
```
packer build debian-stretch-virtualbox.json
```
Packer is set up to handle installing the base OS and guest additions.

### AWS
Once you've made your changes, run
```
packer validate debian-stretch-aws.json
```
and then
```
packer build -force \
-var "aws_access_key=YOUR_ACCESS_KEY" \
-var "aws_secret_key=Your_SECRET_KEY" \
-var "subnet=YOUR_SUBNET" \
-var "vpc=YOUR_VPC" \
debian-stretch-aws.json
```

where YOUR\_SUBNET and YOUR\_VPC are existing subnets and vpcs on your aws
region. The choice of subnet and vpc won't affect the final box, they are
the subnet/vpc which the box is launched into while building. The subnet/vpc ids
should look something like: "subnet-8430fce3" and "vpc-7e99b91a".

After you run packer, it will spit out the ami id. Make sure you remember to so
you can launch instances with it. If you forget it, you can find it under the
"My AMIs" in the Choose an Amazon Machine Image step.
