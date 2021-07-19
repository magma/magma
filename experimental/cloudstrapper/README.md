
# Chapter 1

## 0. Prerequisites

   - Understand key directories and naming conventions
     - CODE_DIR: Directory hosting Magma and Cloudstrapper code that was cloned for this purpose, typically in the ~/code/magma/experimental/cloudstrapper folder
     - VARS_DIR: Directory where all variables reside, typically in the CODE_DIR/playbooks/roles/vars folder
     - WORK_DIR: Directory created by user as part Cloudstrapper deployment with the value taken from dirLocalInventory used as source. This is the directory where all working copies reside, typically in the ~/magma-experimental folder
     - All variables except the S3 bucket are to be named in camel case.
     - S3 bucket names can have only lower case letters and be globally unique

   - Identify two security keys in the Build/Gateway regions to be used for the following
     - Bootkey: Used only by Cloudstrapper instance
     - Hostkey: Used by all Gateway instances
     - Both values are specified in the defaults.yaml vars file and embedded in hosts.
     - [Optional] To generate keys through a playbook, see section 1 below.

     Note: Users who already have preferred keys to be used across their EC2 instances can use them in this environment by setting the keyBoot (for Cloudstrapper) and keyHost (for all other entities created) in defaults.yaml. If such keys do not exist or if the users prefer unique keys for the Cloudstrapper and other AWS artifacts, the aws-prerequisites playbook below will generate the keys.

   - Create inventory directory on localhost to save keys, secrets etc. This directory will be referred to as WORK_DIR and used as dirInventory in commands.
     Ex: mkdir ~/magma-experimental

   - Gather following credentials and update secrets.yaml on local machine in WORK_DIR/secrets.yaml. Use format from $CODE_DIR/playbooks/roles/vars/secrets.yaml as base.
     - AWS Access and Secret keys

   - If you are using the previously built community Magma artifacts for Orc8r and AGW,
     - Locate the Cloudstrapper AMI and version-specific AGW AMI and their respective AMI ids from AWS and make changes to cluster.yaml
     - The AGW AMI id would be set to awsAgwAmi
     - If using the Cloudstrapper tooling to provision the instance, the Cloudstrapper AMI id would be set to awsCloudstrapperAmi

  - If you are building Magma by yourself (Section 3 below) or using your own repositories,

    - Create private repo on github to host helm charts.
     This information would be used in the build.yaml (buildHelmRepo) and cluster.yaml (gitHelmRepo) files when building and deploying Magma respectively.

    - Update local secrets.yaml to include Github and Docker access information for your custom repo
     - Github username and PAT (Personal Access Token)
     - Dockerhub username and password

## 1. Run aws-essentials to setup all AWS related base components as a stack

  The aws-essentials playbook will:
  - Create boot and host keys if required using the keyCreate tag. Default is to not create keys.
  - Create security group on the default VPC
  - Create default bucket for shared storage.
  - Create default CloudFormation stack for all essential components
  Note: If the keys, security groups and the default bucket exist, this playbook can be skipped.

  Before running the playbook, validate the following variables in defaults.yaml and cluster.yaml in VARS_DIR
  - defaults.yaml
    - secgroupDefault indicates name of the default SecurityGroup that would be created and used across all EC2 instances created (including Cloudstrapper and all the AGW instances)
    - bucketDefault indicates the name of the default S3 bucket used to persist information
    - stackEssentialsDefault indicates the name of the CloudFormation stack that hosts all the artifacts
    - keyBoot refers to the key used for the Cloudstrapper instance. Use the keyCreate tag or if it's an existing key, use the right name
    - keyHost refers to the key used for all the other AWS EC2 instances spawned by the Cloudstrapper. Both keys will be generated or skipped together.
    - If two unique keys are already available, update keyBoot and keyHost to reflect the names and skip keyCreate

  - cluster.yaml
    - awsAgwRegion indicates which region would be used to create these artifacts

  - Run the following commands
    ```
    ansible-playbook aws-prerequisites.yaml -e "dirLocalInventory=<directory>" [ --tags keyCreate, essentialsCreate ]
    ```
  - Result: Created a CloudFormation stack with common security group, S3 storage, keys optionally created and .pem files stored in WORK_DIR

### 1.1 For users who do not have access to a Cloudstrapper AMI: Optional CI/CD

  This section is used to create the Cloudstrapper base image. When Magma is available from the Cloud provider's Marketplace, this section will be removed since the Cloudstrapper AMI will be a Marketplace artifact.

  The devops playbooks will :
  - initialize a default instance
  - configure it to act as a Cloudstrapper and
  - generate an AMI that can be used as Cloudstrapper AMI. This could either published in the Marketplace as a public or community AMI or retained locally.

  Before running the playbook, validate the following variables in VARS_DIR
   - defaults.yaml
     - devOpsCloudstrapper indicates the 'Name' tag used to identify the DevOps Cloudstrapper instance
     - primaryCloudstrapper indicates the 'Name' tag used to identify the Primary Cloudstrapper instance
     - devOpsAmi indicates the name of the AMI created for the Cloudstrapper base image
   - build.yaml
    - buildUbuntuAmi - AMI id of base Ubuntu image to be used, available in the region where Cloudstrapper is run

  Run the following commands
  - devops-provision: Setup instance using default security group, Bootkey and Ubuntu
    ```
    ansible-playbook devops-provision.yaml -e "dirLocalInventory=<directory>"
    ```
    - Example:
    - Result: Base instance for Devops provisioned

 - devops-configure: Install ansible, golang, packages, local working directory and latest github sources
     ```
     ansible-playbook devops-configure.yaml -i <dynamic inventory file> -e "devops=tag_Name_<devOpsCloudstrapper>" -e "dirLocalInventory=<inventory folder>" -u ubuntu --skip-tags buildMagma,pubMagma,pubHelm,keyManager
     ```
   - Result: Base instance configured using packages and latest Magma source

 - devops-init: Snapshot instance
    ```
    ansible-playbook devops-init.yaml  -e "dirLocalInventory=<directory>"
    ```
  - Result: DevOps AMI created with name set in devOpsAmi

## 2. Cloudstrapper Process - Marketplace experience begins for users who have access to Cloudstrapper AMI

  - Launch from instance using Bootkey, Ubuntu 20.04 and default security group
    - (or) run cloudstrapper-provision
    ```
    ansible-playbook cloudstrapper-provision.yaml  -e "dirLocalInventory=<directory>"
    ```
  - Result: Cloudstraper node with code package running now, ordered from Marketplace (or based on devOpsAmi for custom builds)
  - Copy keyHost to Cloudstrapper manually or through the playbook to use since that is the seed key for all AWS artifacts created by the Cloudstrapper
  ```
  ansible-playbook devops-configure.yaml -i <Inventory Dir> -e "devops=tag_Name_<primaryCloudstrapper> " -e "dirLocalInventory=<Local Dir>" -u ubuntu -tags keyManager
  ```
  - Login to Cloustrapper node via SSH to start Build, Control Plane and Data Plane rollouts

  - Locate WORK_DIR/magma/experimental/cloudstrapper/playbooks/vars/secrets.yaml file and fill out Secrets
    section and save it in WORK_DIR on Cloudstrapper. Optionally, change other values if required.

## 3. Build

  The build- playbooks provision, configure and initiate the build process before posting the artifacts on identified repositories on successful build. The Build commands can be launched from inside the Cloudstrapper.

  - Create build elements: Provision, Configure and Init.

  - Before beginning Build process, check variables to ensure deployment is customized.
    build.yaml :
      - buildMagmaVersion indicates which version of Magma to build (v1.5 etc)
      - buildOrc8rLabel indicates what label the images would have (1.5.0 etc)
      - buildHelmRepo indicates which github repo will hold Helm charts. Create one if it does not exist. Ensure it is empty.

      - buildAwsRegion indicates which region will host the build instance.
      - buildAwsAz indicates an Availability Zone within the region specified above
      - buildUbuntuAmi reflects the base Ubuntu AMI available in the region described in buildAwsRegion

    All variables can be customized by making a change in the build.yaml file. Invocations
    using Dynamic Inventory would have to be changed to reflect the new labels.

  Run the following commands

  - build-provision: Setup build instance using default security group, Bootkey and Ubuntu with
    t2.xlarge. Optionally, Provision a AGW compliant image (Ubuntu 20.04) for AGW build.
    ```
    ansible-playbook build-provision.yaml -e 'dirLocalInventory=<Inventory file>' --tags devopsOrc8r,inventory
    ```

  - build-configure: Configure build instance by setting up necessary parameters and reading from
    dynamic inventory. The build node was provisioned with the tag Name:buildOrc8r in this example.

    ```
    ansible-playbook build-configure.yaml -i <inventory file> -e "buildnode=tag_Name_<buildTagName>" -e "ansible_python_interpreter=/usr/bin/python3" -e "dirLocalInventory=<inventory folder absolute path>" -u ubuntu
    ```
  - Result: Build instance created, images and helm charts published.

## 4. Control Plane/Cloud Services

  The control role deploys and configures orc8r in a target region. The control playbook can be deployed from within the Cloudstrapper.

  - Create control plane elements: Provision, Configure and Init
    Observe the variables set in cluster.yaml

    Make any custom changes to main.tf here before initializing. If you would like to persist changes
    across re-installs, make changes to the main.tf.j2 Jinja2 template file directly so that the custom
    configuration be used across every terraform init.

  - Clone the latest magma source (master branch) to a local directory named 'source'
  - Requires: secrets.yaml in the dirLocalInventory folder. Use the sample file in roles/vars/secrets.yaml
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
    ansible-playbook orc8r.yaml [ --skip-tags deploy-orc8r ] -e 'dirLocalInventory=<Dir> -e 'varBuildType=community/custom' -e 'varFirstInstall=true/false'
  ```

  Note: When using a stable build or a standard environment or a repeat install, the 'deploy-orc8r' tag does not have to be skipped. However, for first time installs skipping helps in identifying unknown issues to make sure the new build works as expected. Additionally, if there are any custom configuration requirements (such as modifying instance sizes or running multiple clusters within the same account requiring deploy_elasticsearch_service_linked_role to be set to False by default, skipping the deployment and making changes to main.tf is recommended.

  If this tag is skipped, proceed with the following set of commands.

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

  The agw playbooks instantiate a site and configure gateways in it. The design includes two key variables - SiteName (to uniquely identify an edge site) and GatewayName (to uniquely identify a gateway within an edge site). A single site can host multiple gateways and each gateway is associated to one site. While creating the gateway, the SiteName and GatewayName variable are  used. The AGW playbooks can be run from the Cloudstrapper.

  Default variable files are available in the varSite<SiteName>.yaml for sites and varGateway<GatewayName>.yaml for individual gateways from roles/agw-infra/vars/ directory. Newer gateways and sites can also be added following the same format.

  Prerequisites: AGW AMI available and the value specified in 'awsAgwAmi' in cluster.yaml file. If you do not have an AGW AMI file already, please refer to Section 5.1 below to generate an AMI.

  Tunables: (in roles/var/)
    - cluster.yaml:
      - agwAgwRegion: Region where Gateway is deployed
      - awsAgwAz: Availability zone within aforesaid Region
      - awsAgwAmi: AMI ID of AGW AMI in the region, used to deploy in-region or edge gateways
      - orc8rDomainName: Domain name to be used to attach this gateway
    - defaults.yaml:
      - dirSecretsLocal: Local directory with rootCA.pem file

  - Provision the underlying infrastructure and the gateways

    agw-provision: provisions a site with VPC, subnets, gateway, routing tables and one AGW Command:
    ```

    ansible-playbook agw-provision.yaml -e "idSite=<SiteName>" -e "idGw=<GatewayIdentifier>" -e "dirLocalInventory=<WORK_DIR>"[ --tags createNet,createBridge,createGw,cleanupBridge,cleanupNet ]
    ```
    - A site needs to be added only once. After a site is up, multiple gatways can be individually provisioned by skipping the createNet tag as laid out below

    ```
    ansible-playbook agw-provision.yaml -e "idSite=<SiteName>" -e "idGw=<GatewayIdentifier>" -e "dirLocalInventory=<WORK_DIR>" --tags infra,createGw --skip-tags createNet,createBridge,cleanupBridge,cleanupNet
    ```

  - After the gateway has been provisioned, configure the gateway to attach it to an Orchestrator instance. The orchstrator information is picked up from 'cluster.yaml' file

    agw-configure: configures the AGW to include controller information Command:
    ```
    ansible-playbook agw-configure.yaml -i <DynamicInventoryFile> -e "agw=tag_Name_<GatewayId>" -e "dirLocalInventory=<WORK_DIR>" -u ubuntu [ -e KMSKeyID=<KEY_ID_TO_ADD_TO_SSH_AUTHORIZED> -e sshKey=<PATH-TO-PUBLIC-KEY-TO-ADD>]
    ```

    Example:
    ```
    ansible-playbook agw-configure.yaml -i ~/magma-experimental/files/common_instance_aws_ec2.yaml -e "agw=tag_Name_AgwA" "dirLocalInventory=~/magma-experimental/files" -e "ansible_python_interpreter=/usr/bin/python3" [ -e "dirSecretsLocal=<Directory with rootCA.pem>" ] -u ubuntu
    ```
    Proceed to add gateway to Orchestrator using NMS. If you are using Cloudstrapper just to deploy AGWs, ensure rootCA.pem is available from dirSecretsLocal.

  - Result: AWS components created, AGW provisioned, configured and connected to Orchestrator

### 5.1 For users who do not have access to an AGW AMI. Optional CI/CD

  - Tunables (in roles/vars/):

    - build.yaml:
      - buildUbuntuAmi: AMI ID of Base Ubuntu 20.04 image
      - buildAgwAmiName: Name of the AGW AMI created, used to label the AMI
      - buildGwTagName: Tag to be used for the AGW Devops instance, used to filter instance for configuration
      - buildAgwVersion: Version of AGW to be built
      - buildAgwPackage: Specific package version

    - defaults.yaml:

      - keyHost: Name of *.pem file available from <dirInenvtory>, such as ~/magma-experimental/. AWS will use this - value as the key associated with all AGW instances

  - Run the following commands

    - agw-provision: provisions a site with VPC, subnets, gateway, routing tables and one AGW Command:
    ```

    ansible-playbook agw-provision.yaml -e "idSite=DevOps" -e "idGw=<buildGwTagName>" -e "dirLocalInventory=<WORK_DIR>" -e "agwDevops=1" --tags infra,inventory --skip-tags createBridge,cleanupBridge,cleanupNet
    ```

    - ami-configure: Configure AMI for AGW by configuring base AMI image with AGW packages and building OVS.
        - Add ```--skip-tag clearSSHKeys``` if you want to keep ssh keys on the instance
    ```

    ansible-playbook ami-configure.yaml -i <DynamicInventoryFile> -e "dirLocalInventory=<WORK_DIR>" -e "aminode=tag_Name_<buildGwTagName>" -e "ansible_python_interpreter=/usr/bin/python3" -u ubuntu --skip-tags clearSSHKeys

    ```
    - ami-init: Snapshot the AMI instance
    ```
    ansible-playbook ami-init.yaml -e "dirLocalInventory=<Local Inventory Dir>"
    ```

  - Result: AGW AMI created and ready to be used for AGW in-region or Snowcone deployments.

- If you would like to further customize the image with orc8r information and keys for test framework or to ship in a device, run agw-configure and do another ami-init with a different buildAgwAmiName variable.

## 6. Test Framework

  CloudStrapper's test framework allows the user to deploy an all-region version Magma with Orchestrator running in a given region and a cluster of AGWs running in another. Configure number of instances and UUIDs from the local variables file available from vars/main.yaml

  For a multi-node cluster, pre-configure the AMI to embed keys and control_proxy information. This can be done by following Sec 5.1 to build an AGW node, configure that instance using agw-configure and then taking a snapshot via ami-init. Alternatively, individual AGWs can also be configured using the clusterConfigure option.

- Create the network and Bridge node for cluster:
  ```
    - ansible-playbook agw-provision.yaml -e "dirLocalInventory=<Local Dir>" -e "idSite=<Name of site>" --tags createNet,createBridge
  ```
- To start a cluster:
  ```
    - ansible-playbook cluster-provision.yaml -e "dirLocalInventory=<Local Dir>" -e "idSite=<Name of site>" --tags clusterStart
  ```
- To configure a cluster:
  ```
    - ansible-playbook cluster-provision.yaml -e "dirLocalInventory=<Local Dir>" -e "idSite=<Name of site>" --tags clusterConfigure
  ```

- To destroy a cluster:
  ```
    - ansible-playbook cluster-provision.yaml -e "dirLocalInventory=<Local Dir>" -e "idSite=<Name of Site> --tags clusterCleanup
  ```
  - Example:
  ```
    - ansible-playbook cluster-provision.yaml -e "dirLocalInventory=~/magma-experimental" -e "idSite=MenloPark" --tags clusterCleanup
  ```

  Create a SSH configuration to the gateways through the jump node
  ```
  - ansible-playbook cluster-provision.yaml -e "dirLocalInventory=<Local Dir>" -e "agws=tag_Name_<ump_node_name>" -e "idSite=<Name of the Site>" --tags clusterJump
  ```

  - Example:
  ```
  - ansible-playbook cluster-provision.yaml -i /root/project/common_instance_aws_ec2.yaml -e "dirLocalInventory=/root/project" -e "agws=tag_Name_TestFrameworkGateway" -e "idSite=TestCluster" --tags clusterJump
  ```

## 7. Cleanup

  Cleanup deletes all Control and Dataplane components created in the regions. Cleanup can be used to remove all components in one stroke (orchestrator and gateways) or delete individual elements within each layer (database, secrets in orchestrator, any given number of gateways).

  Cleanup uses a combination of native Ansible modules when available and AWS CLI when Ansible  modules are unable to force-delete resources (ex: EFS, Secrets etc.) It uses the name tag of the resources when available and heavily relies on the current assumption that only one orchestrator deployment exists per region. As newer capabilities around tagging emerge, cleanup can be used to target a single deployment among many for cleanup.

  tunables:
    - cluster.yaml:
      - awsOrc8rRegion : Determines which Region hosts the Orc8r instance to be deleted
      - orc8rClusterName : Local folder with terraform state [ex: ~/magma-experimental/<Name of Cluster>]

  - ansible-playbook -e "dirLocalInventory=<Local Dir>" -e "{"deleteStacks": [stackName1, stackName2]}" cleanup.yaml  [ --tags *various* ]

  Available tags include: agw,eks,asg,es,rds,efs,natgw,igw,subnet,secgroup,vpc,orc8r,keys

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
       If a prior code base resides in the home folder, install exists with an error that code already exists. This is done to ensure the user is aware that a prior code base exists and needs to be moved before pulling the new code and not automatically have it overwritten. This is expected behavior.

       Resolution: mv ~/magma ~/magma-backup-<identifier> to  move existing code base.
