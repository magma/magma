
# Chapter 1

## 0. Prerequisites
   - Create private repo on github to host helm charts. 
     This information would be used in the build.yaml and cluster.yaml files when building and deploying Magma respectively. 

   - Identify two security keys in the Build/Gateway regions to be used for the following
     - Bootkey: Used only by Cloudstrapper instance
     - Hostkey: Used by all Gateway instances
       Both values are specified in the defaults.yaml vars file and embedded in hosts.
     [Optional ] To generate keys through a playbook, see section 1 below.
     Customers who already have preferred keys to be used across their EC2 instances can
     use them in this environment. If such keys do not exist or if the customers prefer
     unique keys for the Cloudstrapper, the playbook below will generate the keys.

   - Create inventory directory on localhost to save keys, secrets etc. This directory
     will be referred to as WORK_DIR and used as dirInventory in commands. 
     Ex: mkdir ~/magma-experimental 

   - Gather following credentials and update secrets.yaml on local machine in WORK_DIR. 
     Use format from $CODE_DIR/playbooks/roles/vars/secrets.yaml as base.
     - AWS Access and Secret keys
     - Github username and PAT (Personal Access Token)
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
    under the 'bucketDefault' variable name

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
    ansible-playbook devops-provision.yaml -e "dirLocalInventory=<directory>" 
    ```
  - Example:
    ```
    ansible-playbook devops-provision.yaml -e "dirLocalInventory=~/magma-experimental/files
    ```
  - Result: Base instance for Devops provisioned

  - devops-configure: Install ansible, golang, packages, local working directory and latest github sources
    Command:
    ```
    ansible-playbook devops-configure.yaml -i <dynamic inventory file> -e "< hostname,inventory folder> -u ubuntu --skip-tags usingGitSshKey,buildMagma,pubMagma,helm
    ```
    Example:
    ```
    ansible-playbook devops-configure.yaml -i ~/magma-experimental/files/common_instance_aws_ec2.yaml -e "devops=tag_Name_ec2MagmaDevopsCloudstrapper" -e "dirInventory=~/magma-experimental/files" -u ubuntu --skip-tags buildMagma,pubMagma,helm
    ```
  - Result: Base instance configured using packages and latest Mantle source 

  - devops-init: Snapshot instance  
    Command:
    ```
    ansible-playbook devops-init.yaml  -e "dirLocalInventory=<directory>"
    ```
    Example:
    ```
    ansible-playbook devops-init.yaml  -e "dirLocalInventory=~/magma-experimental/files" 
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
      - buildAwsAz indicates an Availability Zone within the region specified above
    All variables can be customized by making a change in the build.yaml file. Invocations
    using Dynamic Inventory would have to be changed to reflect the new labels.

  - build-provision: Setup build instance using default security group, Bootkey and Ubuntu with
    t2.xlarge. Optionally, Provision a AGW compliant image (Debian 4909 or Ubuntu 20.04) 
    ```
    ansible-playbook build-provision.yaml --tags devopsOrc8r,inventory
    ```

  - build-configure: Configure build instance by setting up necessary parameters and reading from
    dynamic inventory. The build node was provisioned with the tag Name:buildOrc8r in this example.
    ```
    ansible-playbook build-configure.yaml -i ~/magma-experimental/files/common_instance_aws_ec2.yaml -e "buildnode=tag_Name_buildOrc8r" -e "ansible_python_interpreter=/usr/bin/python3"
    ```

  - Result: Build instance created, images and helm charts published. 

## 4. Control Plane/Cloud Services

  The control- roles deploy and configure orc8r in a target region.

  - Create control plane elements: Provision, Configure and Init
    Observe the variables set in cluster.yaml

    Make any custom changes to main.tf here before initializing. If you would like to persist changes
    across re-installs, make changes to the main.tf.j2 Jinja2 template file directly so that the custom
    configuration be used across every terraform init.

  - Requires: secrets.yaml in the dirInventory folder. Use the sample file in roles/vars/secrets.yaml
  - Before beginning Deployment process, check variables to ensure deployment is customized.
    cluster.yaml : 
      - orc8rClusterName: Locally identifiable cluster name
      - orc8rDomainName: DNS name for Orc8r
      - orc8rLabel: Label to look for in container repository
      - orc8rVersion: What version of Orc8r is being deployed
      - gitHelmRepo: Repo which holds helm charts
      - awsOrc8rRegion: Region where Orc8r would run
  - Orchestrator : Deploy orchestrator 
  ```
    ansible-playbook orc8r.yaml [ --skip-tags deploy-orc8r ]
  ```

  Note: First time installs might want to skip using Terraform from within Ansible to make sure the
  new build works as expected. When using a stable build, the tag does not have to be skipped. If this
  tag is skipped, proceed with the following set of commands.

  - Change to local directory
    ```
    cd ~/magma-experimental/<orc8rClusterName defined in roles/vars/cluster.yaml>
    ``` 
  - Run terraform commands manually to provision Cloud resources, load secrets and deploy magma artifacts
    ```
    terraform apply -target=module.orc8r
    terraform apply -target=module.orc8r-app.null_resource.orc8r_seed_secrets
    terraform apply
    ```

  - Result: Orchestrator certificates created, Terraform files initialized, Orchestrator deployed
    via Terraform

  - Validate Orchestrator deployment by following the verification steps in https://magma.github.io/magma/docs/orc8r/deploy_install

## 5. Data Plane
  The agw playbooks instantiate a site and configure gateways in it. The design includes two key
  variables - SiteName (to uniquely identify an edge site) and GatewayName (to uniquely identify
  a gateway within an edge site). A single site can host multiple gateways and each gateway is
  associated to one site. While creating the gateway, the SiteName and GatewayName variable are
  used. 

  Default variable files are available in the varSite<SiteName>.yaml for sites and 
  varGateway<GatewayName>.yaml for individual gateways from roles/agw-infra/vars/ directory. Newer
  gateways and sites can also be added following the same format.

  Prerequisites: AGW AMI available and the value specified in 'awsAgwAmi' in cluster.yaml file. If
  you do not have an AGW AMI file already, please refer to Section 5.1 below to generate an AMI.
 
  Tunables: (in var/opts)
    cluster.yaml:
      agwAgwRegion: Region where Gateway is deployed
      awsAgwAz: Availability zone within aforesaid Region
      awsAgwAmi: AMI ID of AGW AMI in the region, used to deploy in-region or edge gateways 
      orc8rDomainName: Domain name to be used to attach this gateway
    defaults.yaml:
      dirSecretsLocal: Local directory with rootCA.pem file
 
  - Provision the underlying infrastructure and the gateways

    agw-provision: provisions a site with VPC, subnets, gateway, routing tables and one AGW Command:
    ```
    ansible-playbook agw-provision.yaml -e "idSite=<SiteName>" -e "idGw=<GatewayIdentifier>" [ --tags createNet,createGw,attachIface ]
    ```
    Example:
    ```
    ansible-playbook agw-provision.yaml -e "idSite=MenloPark" -e "idGw=AgwA" [ --tags createNet createGw attachIface ]
    ``` 

  - A site needs to be added only once. After a site is up, multiple gatways can be individually 
    provisioned by skipping the createNet tag as laid out below
    ```
    ansible-playbook agw-provision.yaml -e "idSite=<SiteName>" -e "idGw=<GatewayIdentifier>" --tags createGw,attachIface
    ``` 

    Example:
    ```
    ansible-playbook agw-provision.yaml -e "idSite=MenloPark" -e "idGw=AgwB" --tags createGw,attachIface
    ``` 

  - After the gateway has been provisioned, configure the gateway to attach it to an Orchestrator 
    instance. The orchstrator information is picked up from 'cluster.yaml' file

    agw-configure: configures the AGW to include controller information Command:
    ```
    ansible-playbook agw-configure.yaml -i <DynamicInventoryFile> -e "agw=tag_Name_<SiteId><GatewayId>" -u admin
    ```
    Example:

    ```
    ansible-playbook agw-configure.yaml -i ~/magma-experimental/files/common_instance_aws_ec2.yaml -e "agw=tag_Name_MenloParkAgwA" -u admin
    ```
    Proceed to add gateway to Orchestrator using NMS. 

  - Result: AWS components created, AGW provisioned, configured and connected to Orchestrator

### 5.1 For users who do not have access to an AGW AMI. Optional CI/CD

  - Prerequisite: Debian Stretch AMI with 4.9.0.9 kernel with AMI id. 
  - Tunables (in roles/vars): 
    build.yaml:
      buildDebianAmi: AMI ID of Base Debian Stretch 4.9.0.9 image
      buildAgwAmiName: Name of the AGW AMI, used to label AMI 
      buildGwTagName: Tag to be used for the AGW Devops instance, used to filter instance
    defaults.yaml:  
      keyHost: Name of *.pem file available from <dirInenvtory>, such as ~/magma-experimental/files 
      secGroupDefault: Name of the default security group to be used for the AGW instance

  - build-provision: Setup AGW devops build instance using default security group, Bootkey using a
    a AGW compliant image (Debian 4909). This AMI value is specified in the build.yaml file as  
    'buildDebianAmi'.
    ```
    ansible-playbook build-provision.yaml --tags devopsAgw,inventory
    ```

 - ami-configure: Configure AMI for AGW by configuring base AMI image with AGW packages and building OVS.
    ```
    ansible-playbook ami-configure.yaml -i <DynamicInventoryFile> -e "aminode=tag_Name_<buildGwTagName>" -u admin
    ```
    Example:
    ```
    ansible-playbook ami-configure.yaml -i ~/magma-experimental/files/common_instance_aws_ec2.yaml -e "aminode=tag_Name_buildAgw" -u admin
    ```

  - ami-init: Snapshot the AMI instance
    ```
    ansible-playbook ami-init.yaml
    ```

  - Result: AGW AMI created and ready to be used for AGW in-region or Snowcone deployments.

## 6. Cleanup

  Cleanup deletes all Control and Dataplane components created in the regions. Cleanup can be used
  to remove all components in one stroke (orchestrator and gateways) or delete individual
  elements within each layer (database, secrets in orchestrator, any given number of gateways).

  Cleanup uses a combination of native Ansible modules when available and AWS CLI when Ansible 
  modules are unable to force-delete resources (ex: EFS, Secrets etc.) It uses the name tag of the resources
  when available and heavily relies on the current assumption that only one orchestrator deployment
  exists per region. As newer capabilities around tagging emerge, cleanup can be used to target a single
  deployment among many for cleanup.

  tunables: 
    - cluster.yaml: 
        awsOrc8rRegion : Determines which Region hosts the Orc8r instance to be deleted
        orc8rClusterName : Local folder with terraform state [ex: ~/magma-experimental/<Name of Cluster>]

  - ansible-playbook cleanup.yaml  [ --tags *various* ]

  Available tags include: agw,eks,asg,es,rds,efs,natgw,igw,subnet,secgroup,vpc

## Known Issues, Best Practices & Expected Behavior

### Best Practices
    
    1. Although the deployment will work from any Ubuntu host, using the Cloudstrapper AMI might be
       the quickest way to get the deployment going since it includes all the necessary dependencies
       in-built. 

    2. The tool is customizable to build every desired type of installation. However, for initial
       efforts, it might be better to to use the existing default values. 

    3. Some resources are not covered under the current Cleanup playbooks. This includes Route53 
       entries, keypairs and AWS roles created by Orchestrator since there is no clear way to distinguish
       them from other resources that share the same name. 

    4. Due to a cyclical dependency in orchestrator security groups, some rules have to be manually removed
       until the [issue] (https://github.com/magma/magma/issues/5150) is fixed upstream. 

### Expected Behavior - Install

    1. Prior code base
       If a prior code base resides in the home folder, install exists with an error that code already exists.
       This is done to ensure the user is aware that a prior code base exists and needs to be moved
       before pulling the new code and not automatically have it overwritten. This is expected
       behavior. 
       
       Resolution: mv ~/magma ~/magma-backup-<identifier> to  move existing code base.
