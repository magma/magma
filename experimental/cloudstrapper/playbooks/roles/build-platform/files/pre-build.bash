#!/bin/bash

#Objective: Run common tasks to setup or destroy pre-build elements
#Pre-requisite: AWS profile is configured with access key, secret key, region and
#json
#Philosophy: Use CloudFormation (cfn) when applicable

source ./global.env

function prebuild-create {
#Create EC2 private key pair.
#For everything else, use CloudFormation

aws ec2 create-key-pair --key-name $AWS_ANSIBLE_KEY --query 'KeyMaterial' --output text > $AWS_ANSIBLE_KEY.pem
chmod 600 $AWS_ANSIBLE_KEY.pem 

aws cloudformation create-stack --stack-name $AWS_CFN_PREBUILD_STACK --template-body file://$PROTO_CONFIG/cfnMagmaPreBuildProto.json
}


function prebuild-cleanup {

aws ec2 delete-key-pair --key-name $AWS_ANSIBLE_KEY
rm -f $AWS_ANSIBLE_KEY.pem

aws cloudformation delete-stack --stack-name $AWS_CFN_PREBUILD_STACK

}

#Check if AWS creds are configured
if [ $# -ne 1 ]
then
	echo "Need one argument. Create or Cleanup"
	exit 1
fi

if [ $1 = "Create" ]
then
	echo "Create" 
	prebuild-create
elif [ $1 = "Cleanup" ]
then
	echo "Cleanup"
	prebuild-cleanup
else
	echo "Error. Argument must be Create or Cleanup"
fi


