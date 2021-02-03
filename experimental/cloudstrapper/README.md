
# Chapter 1

## 0. Prerequisites
   - Create private repo on github to host helm charts. 
     This information would be used in the build.yaml and cluster.yaml files when building and deploying Magma respectively. 

   - Identify two security keys in the Build/Gateway regions to be used for the following
     - Bootkey: Used only by Cloudstrapper instance
     - Hostkey: Used by all Gateway instances
       Both values are specified in the defaults.yaml vars file and embedded in hosts.
     To generate keys through a playbook, see section 1 below.

   - Create inventory directory on localhost to save keys, secrets etc. This directory
     will be referred to as WORK_DIR and used as dirInventory in commands.

   - Gather following credentials and update secrets.yaml on local machine in WORK_DIR. 
     Use format from $CODE_DIR/playbooks/roles/vars/secrets.yaml as base.
     - AWS Access and Secret keys
     - Github username and PAT
     - Dockerhub username and password 

   - Understand key directories
     - CODE_DIR: Directory hosting Magma code, typically in the ~/code/magma/experimental/cloudstrapper folder
     - VARS_DIR: Directory where all variables reside, typically in the CODE_DIR/playbooks/roles/vars folder
     - WORK_DIR: Directory where all working copies reside, typically n the ~/magma-experimental
       folder

## 1. Run aws-essentials to setup all AWS related components as a stack

  The aws-essentials playbook will:
  - Create boot and host keys if required using the keyCreate tag. Default is to not create keys.
  - Create security group on the default VPC
  - Create default bucket for shared storage. Ensure bucket does not exist by checking defaults.yaml

  - Command:
    ```
    ansible-playbook aws-prerequisites.yaml -e 'awsTargetRegion=<< AWS Region >>' -e "dirInventory=<directory>" [ --tags keyCreate ]
    ```
  - Example:
    ```
    ansible-playbook aws-prerequisites.yaml -e 'awsTargetRegion=us-east-1' -e "dirInventory=~/magma-experimental/files" 
    ``` 
  - Result: Created stackMantleEssentials with common security group, S3 storage

### 1.1 For users who do not have access to a Cloudstrapper AMI: Optional CI/CD
  The devops playbooks are used to initialize a default instance, configure it to act as a Cloudstrapper and generate an
  AMI that can be used as Cloudstrapper AMI and either published in the Marketplace as a public or community AMI or 
  retained locally.


  - devops-provision: Setup instance using default security group, Bootkey and Ubuntu 
  - Command:
    ```
    ansible-playbook devops-provision.yaml -e "dirInventory=<directory>" 
    ```
  - Example:
    ```
    ansible-playbook devops-provision.yaml -e "dirInventory=~/magma-experimental/files
    ```
  - Result: Base instance for Devops provisioned

  - devops-configure: Install ansible, golang, packages, local working directory and latest github sources
    Command:
    ```
    ansible-playbook devops-configure.yaml -i <dynamic inventory file> -e "< hostname,inventory folder> -u ubuntu --skip-tags usingGitSshKey,buildMagma,pubMagma,helm
    ```
    Example:
    ```
    ansible-playbook devops-configure.yaml -i ~/magma-experimental/files/common_instance_aws_ec2.yaml -e "devops=tag_Name_ec2MagmaDevopsCloudstrapper" -e "dirInventory=~/magma-experimental/files" -u ubuntu --skip-tags usingGitSshKey,buildMagma,pubMagma,helm
    ```
  - Result: Base instance configured using packages and latest Mantle source 

  - devops-init: Snapshot instance  
    Command:
    ```
    ansible-playbook devops-init.yaml  -e "dirInventory=<directory>"
    ```
    Example:
    ```
    ansible-playbook devops-init.yaml  -e "dirInventory=~/magma-experimental/files" 
    ```
  - Result: imgMagmaCloudstrap AMI created

## 2. Cloudstrapper Process - Marketplace experience begins for users who have access to Cloudstrapper AMI

  - Launch from instance using Bootkey, Ubuntu 20.04 and default security group
    - (or) run cloudstrapper-provision
      ```
      ansible-playbook cloudstrapper-provision.yaml  -e "dirLocalInventory=~/magma-experimental/files"
      ```
  - Result: Cloudstraper node with code package running now, ordered from Marketplace

  - Login to Cloustrapper node via SSH to start Build, Control Plane and Data Plane rollouts

  - Locate ~/code/mantle/magma-on-aws/playbooks/vars/secrets.yaml file and fill out Secrets
    section and save it in WORK_DIR on Cloudstrapper. Optionally, change other values if required.

## 3. Build

  The build- playbooks provision, configure and initiate the build process before posting 
  the artifacts on identified repositories on successful build.

  - Create build elements: Provision, Configure and Init. 
     
  - Before beginning Build process, check variables to ensure deployment is customized.
    build.yaml : 
      - buildMagmaVersion indicates which version of Magma to build (v1.3, v1.4 etc)
      - buildOrc8rLabel indicates what label the images would have
      - buildHelmRepo indicates which github repo will hold Helm charts

      - buildAwsRegion indicates which region will host the build instance. 
      - buildAwsAz indicates an AZ within the region specified above

  - build-provision: Setup build instance using default security group, Bootkey and Ubuntu with
    t2.xlarge. Optionally, Provision a AGW compliant image (Debian 4909 or Ubuntu 20.04) 
    ```
    ansible-playbook build-provision.yaml
    ```

  - build-configure: Configure build instance by setting up necessary parameters and reading from
    dynamic inventory. The build node was provisioned with the tag Name:buildOrc8r in this example.
    ```
    ansible-playbook build-configure.yaml -i ~/magma-experimental/files/common_instance_aws_ec2.yaml -e "buildnode=tag_Name_buildOrc8r" -e "ansible_python_interpreter=/usr/bin/python3"
    ```

  - build-ami-configure: Configure AMI for AGW by configuring base AMI image with AGW packages and
    building OVS.
    *TODO: Remove existing snapshot and create new snapshot*
    ```
    ansible-playbook build-ami-configure.yaml -e '@vars/main.yaml' -i files/build_instance_aws_ec2.yaml -e "buildnode=tag_Name_ec2MagmaBuild" -u admin
    ```

  - Result: Build instance created, images and helm charts published. AGW AMI created.

## 4. Control Plane/Cloud Services

  The control- roles deploy and configure orc8r in a target region.

  - Create control plane elements: Provision, Configure and Init
    Observe the variables set in cluster.yaml

    Make any custom changes to main.tf here before initializing. If you would like to persist changes
    across re-installs, make changes to the main.tf.j2 Jinja2 template file directly so that the custom
    configuration be used across every terraform init.

  - Requires: secrets.yaml in the dirInventory folder. Use the sample file in roles/vars/secrets.yaml
  - Orchestrator : Deploy orchestrator 
  ```
    ansible-playbook orc8r.yaml [ --skip-tags deploy-orc8r ]
  ```

  Note: First time installs might want to skip using Terraform from within Ansible to make sure the
  new build works as expected. When using a stable build, the tag does not have to be skipped.

  - Result: Orchestrator certificates created, Terraform files initialized, Orchestrator deployed
    via Terraform

## Known Issues, Best Practices & Expected Behavior

### Best Practices
    
    1. Although the deployment will work from any Ubuntu host, using the Cloudstrapper AMI might be
       the quickest way to get the deployment going since it includes all the necessary dependencies
       in-built. 

    2. The tool is customizable to build every desired type of installation. However, for initial
       efforts, it might be better to to use the existing default values. 

### Expected Behavior - Install

    1. Prior code base
       If a prior code base resides in the home folder, install exists reporting code already exists.
       This is done to ensure the user is aware that a prior code base exists and needs to be moved
       before pulling the new code and not automatically have it overwritten. This is expected
       behavior. 


