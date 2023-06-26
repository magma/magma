---
id: version-1.8.0-aws_cloudstrapper
title: AWS Cloudstrapper Install
hide_title: true
original_id: aws_cloudstrapper
---

# Deploying Magma via Cloudstrapper

There are two basic options for setting up Magma Cloudstrapper within Amazon Web Services: Marketplace or via private image.

## 1) Launching Cloudstrapper from Marketplace

- Access the [AWS Marketplace](https://aws.amazon.com/marketplace) and search for “Magma Cloudstrapper”
    - Alternatively, check the [direct link](https://aws.amazon.com/marketplace/pp/prodview-wkchyk2okdnhc?qid=1627070115980&sr=0-2&ref_=srh_res_product_title)
- Click on “Continue to Subscribe” and “Continue to Configuration”
- Choose “Delivery Method”, “Software Version” and “Region.” The only default we recommend you change is the “Region"
- Click on “Continue to Launch”
    - In “Choose Action”, select “Launch from Website” (default)
    - The EC2 Instance Type dropdown will select “t2.medium” by default
    - Choose preferred values for other drop-boxes. Cloudstrapper will work fine deployed on the public subnet.
    - Under “Security Group Settings” select a security group that allows SSH traffic and any other rules that are relevant to your network.
    - Under “Key Pair Settings” select your preferred key pair.
- Click on Launch
- In order to ssh into your Cloudstrapper use the key pair .pem file and ubuntu in this format: ssh -i &lt;KeyPair&gt; ubuntu@&lt;InstanceIP&gt;
    - Example:
- ssh -i "cloudstrapper-test.pem" ubuntu@1.1.1.1
    - **NOTE:** If you receive an “WARNING: UNPROTECTED PRIVATE KEY FILE!” error while trying to SSH you will need to change the permissions on your key file by running `chmod 400` to make it read-only.

## 2) Launching Cloudstrapper from Private Images

- Navigate to the “AMIs” page to verify Cloudstrapper images show up under “Images Owned By Me”
- Select the Cloudstrapper AMI followed by the “Launch” button.
- On the “Choose Instance Type” page, select the t2.micro instance type (free) and click next.
- No changes are required to default settings on the “Configure Instance Details” page. Click next and navigate to the “Add Storage” page.
- Update the Size parameter to 32GB of space and proceed to the “Add Tags” page.
- Select the “Add Tag” button or click the link to add a name tag. In the “Value” column enter in a name for your Cloudstrapper and in the “Key” field create a key. Proceed to the “Configure Security Group” page.
- On this page, select a security group that at allows SSH traffic (default) and any other rules that are relevant to your network.
- Proceed to the “Review” page to ensure all your network details are correct followed by clicking the Launch button.
- Finally, select a pre-existing key pair that will allow access to your network and launch your instances.

## Configure Orchestrator Deployment Parameters

- Once your instance is created, click on the “Instance ID” url and select the “Connect” button
- In the SSH tab copy the command to connect to your instance and run it within your CLI
    - **NOTE:** If you receive an “WARNING: UNPROTECTED PRIVATE KEY FILE!” error while trying to SSH you will need to change the permissions on your key file by running `chmod 400` to make it read-only.
- Create a magma-dev directory and clone magma master onto it by running the following commands:

```bash
mkdir ~/magma-dev
cd ~/magma-dev
git clone https://github.com/magma/magma.git
```

- Locate and navigate to the playbooks directory inside the source repo:

    - `cd ~/magma-dev/magma/experimental/cloudstrapper/playbooks`

- Copy the secrets.yaml file and update the credentials for AWS:

    - `cp roles/vars/secrets.yaml ~/magma-dev/`
- Update your secrets.yaml file with your AWS credentials. The two fields that are required to update are the `AWS Access Key` and the `AWS Secret Key`. **Note:** make sure you have a space between the colon and before your keys.

- Build types: Magma can be deployed via binaries hosted in the community artifactory. These options are enabled in the code by default and Cloudstrapper, when using the ‘community’ option, it will default to these options.
    - Community Builds: Builds are created and labeled by Magma CI teams and available for deployment. The secrets.yaml file does not need any inputs in the docker and github variables.
    - Private Builds: Edit the secrets.yaml file in the `/magma-dev` directory to include the github and docker credentials under the Github Prerequisites and Dockerhub Prerequisites fields.

**Please ensure that:**

- You are deploying orchestrator in a region that supports at least **three Availability Zones**
- The region is clean and has no leftover artifacts from a prior Orc8r deployment. Please use the cleanup script below if needed.
- the value of varFirstInstall is set based on if the account already has an “AWSServiceRoleForAmazonElasticsearchService” role created (to be automated). If it exists, varFirstInstall in the config file would be “false”. If not, varFirstInstall in the config file would be “true”. If the role does not exist already, Orc8r will create the role.
- there is disk space of at least 1+ GB to pull code and create local artifacts
- dirLocalInventory/orc8rClusterName folder does not exist from a previous install (to be automated)

```bash
aws iam list-roles --profile <Your AWS Config profile> | grep -i AWSServiceRoleForAmazonElasticsearchService

"RoleName": "AWSServiceRoleForAmazonElasticsearchService",
"Arn": "arn:aws:iam::<Account Number>:role/aws-service-role/es.amazonaws.com/AWSServiceRoleForAmazonElasticsearchService"     
```

## Deploy Orchestrator

- Prior to completing the following steps you must obtain a domain name for your Orc8r
- Create a parameter file (must end in a .yaml extension) in the ~/magma-dev directory
    - View the example below to examine what a sample parameter file would look like. You must reserve a domain name before completing this step.

```bash
---
dirLocalInventory: ~/magma-dev
orc8rClusterName: Sydney
orc8rDomainName: ens-16-sydney.failedwizard.dev
orc8rLabel: 1.8.0
orc8rVersion: v1.8
awsOrc8rRegion: ap-southeast-2
varBuildType: community
varFirstInstall: "false"
```

- A legend of the variables used:
    - dirLocalInventory: Folder which has secrets.yaml file that includes AWS access and secret keys
    - orc8rClusterName: A local folder created in dirLocalInventory used to store state information
    - orc8rDomainName: Domain name of the Orc8r
    - orc8rLabel: The label to look for in the containers repository
    - orc8rVersion: The version of orc8r tools used to generate artifacts
    - awsOrc8rRegion: The region where this orc8r will be deployed
    - varBuildType: Using either ‘community’ or ‘custom’ binaries
    - varFirstInstall: Indicating if this is the first install of any kind of Magma or not, to skip some of the default, shared roles created

- First change your directory to `~/magma-dev/magma/experimental/cloudstrapper/playbooks`. Next, run the playbook to set up the Orchestrator deployment:

```bash
ansible-playbook orc8r.yaml -e '@<path to parameters file>'
```

- After a successful run of the playbook (30-40 minutes), run terraform to obtain nameserver information to be added to DNS.

Example:

```bash
cloudstrapper:~/magma-dev/Mumbai/terraform #terraform output
nameservers = tolist([
"ns-1006.awsdns-61.net",
"ns-1140.awsdns-14.org",
"ns-2020.awsdns-60.co.uk",
"ns-427.awsdns-53.com",
])
```

## Deploy the AGW

- If you would like to do a customized installation (Ex. generating your own classless interdomain routing) you will first need to create your own .yaml file
    - This is done by navigating to `magma-dev/magma/experimental/cloudstrapper/playbooks/roles/agw-infra/vars/`
    - Next, create your own .yaml file where you can configure your unique parameters. Here is an example:

```bash
**cidrVpc**: 10.7.0.0/16
**cidrSgi**: 10.7.4.0/24
**cidrEnodeb**: 10.7.2.0/24
**cidrBridge**: 10.7.6.0/24
**azHome**: "{{ awsAgwAz }}"
**secGroup**: "{{ secgroupDefault }}"
**sshKey**: "{{ keyHost }}"
**siteName**: MenloPark
```

- If you would like to do a non-customized setup you can use the `MenloPark` idSite for your installation for the steps below.

- Create a parameter file in the `~/magma-dev` directory. A sample parameter file would look like as follows
    - **Note:** Use a base ubuntu image of the region for the awsCloudstrapperAmi variable.
    - Please ensure you have a key with name described in keyHost available in the region. This is the key that would be embedded into the AGW for ssh access
    - Future: awsCloudstrapperAmi will be renamed to awsBastionAmi (task filed)

```bash
    ---
    dirLocalInventory: ~/magma-dev
    awsAgwAmi: ami-00ca08f84d1e324b0
    awsCloudstrapperAmi: ami-02f1c0266c02f885b
    awsAgwRegion: ap-northeast-1
    keyHost: keyMagmaHostBeta
    idSite: MenloPark
    idGw: mpk01


```

- A legend of the variables used:
    - dirLocalInventory: Location of folder with secrets.yaml file that include AWS access and secret keys
    - awsAgwAmi: Id of the AGW AMI available in the region of deployment
    - awsCloudstrapperAmi: Id of the Cloudstrapper AMI available in the region of deployment
    - awsAgwRegion: Region where AGWs will be deployed
    - keyHost: Public key of keypair from region, will be embedded into the launched AGW for SSH access
    - idSite: ID of site and partial name of variable file that has site specific information (such as CIDRs) ([Example](https://github.com/magma/magma/blob/master/experimental/cloudstrapper/playbooks/roles/agw-infra/vars/varSiteMenloPark.yaml))
    - idGw: Id of Gateway to be installed; Used as value of a tag with key as Name; Can be changed for subsequent AGW deployments

- Locate and navigate to the playbooks directory inside the source repo:

```bash
cd ~/magma-dev/magma/experimental/cloudstrapper/playbooks
```

- Run the following command for the first AGW:

```bash
ansible-playbook agw-provision.yaml --tags createNet,createBridge,createGw,inventory -e '@<path to parameters file>'

```

- Run the following for subsequent AGWs :

```bash
ansible-playbook agw-provision.yaml --tags createGw -e '@<path to parameters file>'
```

## Configure AGW and Orchestrator

- Configure AGW [manually](https://magma.github.io/magma/docs/lte/config_agw) or through the playbooks running agw-configure from the Bridge node.
- Start by configuring the Bridge node as a bastion host. Using the Bridge node as Bastion host, configure the newly deployed AGW to communicate with the Orc8r.
- Follow the following steps [here](https://magma.github.io/magma/docs/orc8r/deploy_install) to create an admin user for NMS.
- Generate a challenge key and hardware id and add it to Orc8r from the [Magmacore website documentation](https://magma.github.io/magma/docs/lte/deploy_config_agw).

## Custom AGWs for Snowcone

Snowcone devices require that an image be already embedded in the device before shipping. Hence, the devices require a key to be embedded in the authorized_keys file for the default user (‘ubuntu’) or a customer’s preferred user before the device is ordered. To achieve that run the following commands.

- To generate a custom image from the base Access Gateway AMI

    1. Start an instance of the Access Gateway AMI similar to the “Launching Cloudstrapper” session above.
    2. Once the instance is booted up, add your public key to the ~/.ssh/authorized_keys file
    3. Snapshot the image to create a new AMI. Use this AMI to order your snowcone.

- Alternatively, to build your own AGW AMI from scratch and customize it for your use-case

    1. Follow section 5.1 from the [README](https://github.com/magma/magma/tree/master/experimental/cloudstrapper) file of the Cloudstrapper deployment

## Cleaning up an Orc8r environment

Orc8r cleanup allows the user to target a given region and automatically cleanup all the resources there and ensure it is ready for a new deployment.

Run `terraform destroy` to release all terraform created resources from within the ‘terraform’ directory inside the ‘dirLocalInventory/orc8rClusterName’ folder.

In certain cases, terraform might leave artifacts that weren’t cleaned up that might impact previously deployment installations. Use Cloudstrapper’s cleanup capabilities to cleanup the region of all known artifacts.

**Variables to consider for cleanup:**

- awsOrc8rRegion: Region where Orchestrator runs
    - command: `ansible-playbook cleanup.yaml [--tags various] -e '@<path to parameters file>'`

- For a complete environment cleanup or orc8r run:
    - command: `ansible-playbook cleanup.yaml —skip-tags agw -e ‘@<path to parameters file>’`
