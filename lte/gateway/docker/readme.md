# Containerization Deploy

* Use the Base Cloudstrapper image with an expanded disk size (64G at least) to create Build Container
    * If you want to use a stand-alone Ubuntu image, [Setup](https://docs.docker.com/engine/install/ubuntu/) Docker to run [locally](https://forums.docker.com/t/couldnt-connect-to-docker-daemon-at-http-docker-localhost-is-it-running/87257). Remember to re-login after local user got added to docker so that the new shell has all the right user provisions
    * fb-magma-aws has ami-03bc7ef7f3b70f77b which is a base Cloudstrapper image with extended storage
    * Run Cloudstrapper AGW deploy with input file as below (SDTI network already exists)
    * devsrv-tokyo:~/magma-master #cat sdti-build1.yaml

        ```
        #Setting AGW AMI and Cloudstrapper AMI to expanded Cloudstrapper image Ubuntu to allow deploy
        dirLocalInventory: ~/magma-master
        awsAgwAmi: ami-03bc7ef7f3b70f77b
        awsCloudstrapperAmi: ami-03bc7ef7f3b70f77b
        awsAgwRegion: ap-northeast-1
        keyHost: keyMagmaHostCharlie
        idSite: SDTI
        idGw: sdti-build1
        awsInstanceType: t3.large
        
        devsrv-tokyo:~/magma-master/magma/experimental/cloudstrapper/playbooks #ansible-playbook —tags createGw agw-provision.yaml -e '@~/magma-master/sdti-build1.yaml'
        ```

* Configure ssh config file to allow SSH access to CTR1
* System information : Runs 5.4 kernel, has IP addresses configured on first and second interfaces

```
sdti-build1:~/magma/lte/gateway/docker #uname -r
5.4.0-1045-aws
sdti-build1:~/magma/lte/gateway/docker #ip addr show
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
2: ens5: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 9001 qdisc mq state UP group default qlen 1000
    link/ether 06:92:f5:05:4a:07 brd ff:ff:ff:ff:ff:ff
    inet 10.23.4.220/24 brd 10.23.4.255 scope global dynamic ens5
       valid_lft 3358sec preferred_lft 3358sec
    inet6 fe80::492:f5ff:fe05:4a07/64 scope link
       valid_lft forever preferred_lft forever
3: ens6: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 9001 qdisc mq state UP group default qlen 1000
    link/ether 06:a4:58:04:db:d3 brd ff:ff:ff:ff:ff:ff
    inet 10.23.2.65/24 brd 10.23.2.255 scope global dynamic ens6
       valid_lft 3356sec preferred_lft 3356sec
    inet6 fe80::4a4:58ff:fe04:dbd3/64 scope link
       valid_lft forever preferred_lft forever
4: docker0: <NO-CARRIER,BROADCAST,MULTICAST,UP> mtu 1500 qdisc noqueue state DOWN group default
    link/ether 02:42:5c:43:5a:64 brd ff:ff:ff:ff:ff:ff
    inet 172.17.0.1/16 brd 172.17.255.255 scope global docker0
       valid_lft forever preferred_lft forever
    inet6 fe80::42:5cff:fe43:5a64/64 scope link
       valid_lft forever preferred_lft forever
```


* Temporary steps
    * sudo apt-get update
    * Install ifupdown
* Run docker-compose build (~45 minutes)

```
Step 56/56 : RUN chmod -R +x /usr/local/bin/generate* /usr/local/bin/set_irq_affinity /usr/local/bin/checkin_cli.py &&   dpkg -i /var/tmp/python3-aioeventlet* &&   pip install jsonpointer>$JSONPOINTER_VERSION &&   mkdir -p /var/opt/magma/
 ---> Running in 46435f4fbad2
Selecting previously unselected package python3-aioeventlet.
(Reading database ... 47574 files and directories currently installed.)
Preparing to unpack .../python3-aioeventlet_0.5.1-2focal_amd64.deb ...
Unpacking python3-aioeventlet (0.5.1-2) ...
Setting up python3-aioeventlet (0.5.1-2) ...
Removing intermediate container 46435f4fbad2
 ---> 08ed545b9000

Successfully built 08ed545b9000
Successfully tagged agw_gateway_python:latest
sdti-build1:~/magma/lte/gateway/docker #
```
* Check images on local host
* ctr-build:~/magma-ctr/magma/lte/gateway/docker #docker images

```
    REPOSITORY           TAG       IMAGE ID       CREATED          SIZE
    agw_gateway_c        latest    755b206a9698   3 minutes ago    1.29GB
    <none>               <none>    5163a21390e8   6 minutes ago    4.41GB
    agw_gateway_python   latest    d61fcb3ed86e   26 minutes ago   894MB
    <none>               <none>    6cb8b28381ce   29 minutes ago   1.71GB
    <none>               <none>    06c85a31b481   2 hours ago      4.41GB
    <none>               <none>    b6e324a3c73a   3 hours ago      1.71GB
    ubuntu               focal     ba6acccedd29   3 weeks ago      72.8MB
```

* Ensure repositories are created on dockerhub
* Tag and Push images to docker hub

```
ctr-build:~/magma-ctr/magma/lte/gateway/docker #docker image tag agw_gateway_python:latest arunuke/agw_gateway_python:9Nov
ctr-build:~/magma-ctr/magma/lte/gateway/docker #docker image tag agw_gateway_c:latest arunuke/agw_gateway_c:9Nov
ctr-build:~/magma-ctr/magma/lte/gateway/docker #docker image push arunuke/agw_gateway_c:9Nov
ctr-build:~/magma-ctr/magma/lte/gateway/docker #docker image push arunuke/agw_python:9Nov
```

* Use the Base Cloudstrapper image with an expanded disk size (64G at least) to create Test Container

```
ansible-playbook agw-provision.yaml --tags createGw -e '@~/magma-master/sdti-ctr1.yaml'

evsrv-tokyo:~/magma-master/magma/experimental/cloudstrapper/playbooks #cat ~/magma-master/sdti-ctr1.yaml
---
#Setting AGW AMI and Cloudstrapper AMI to expanded Cloudstrapper image Ubuntu to allow deploy
dirLocalInventory: ~/magma-master
awsAgwAmi: ami-03bc7ef7f3b70f77b
awsCloudstrapperAmi: ami-03bc7ef7f3b70f77b
awsAgwRegion: ap-northeast-1
keyHost: keyMagmaHostCharlie
idSite: SDTI
idGw: sdti-ctr1
awsInstanceType: t3.large
```

* Prepare the host
    * Install ifupdown
    * Unlink `/etc/resolv.conf` and create a new one with 8.8.8.8 entry
    * Create `/var/opt/magma/certs `and add rootCA.pem to that folder with permissions 400
    * Copy `agw_install_docker.bash` and run script to prepare the host
    * Make changes to `/var/opt/magma/docker/.env` to include the right docker information
    * DOCKER_REGISTRY=[registry.hub.docker.com/arunuke/](http://registry.hub.docker.com/arunuke/)

       ```
        DOCKER_USERNAME=arunuke
        DOCKER_PASSWORD=XXX
       ```

    * Make changes to files in `/var/opt/magma/configs` if needed
        * pipelined, dnsd, enodebd, spgw, mme all will have their eth0/eth1 changed to newer values based on local interface names (eth0 and eth2 references to use the first interface for SGi and eth1 references to use the second interface for S1)
        * pipelined will also set dp_router_enabled to false
    * Make changes to config files and restart services by running `/var/opt/magma/docker/agw_upgrade.sh` or by running the `agw_install_docker.bash` script
* Issues
    * Needs ifup on the host. Install package ifupdown (PR in progress)
    * Need to resolve external IP addresses after bringing up interfaces. unlink /etc/resolv.conf, add a new entry for 8.8.8.8. Add this by creating a new role. (PR needed)
    * Need to setup variables in the [variables file](https://github.com/magma/magma/tree/master/lte/gateway/deploy/roles/agw_docker/vars) under magma_root (/opt/magma) and it has to be documented in the README (PR needed)
    * Need to fix externally pulled images from docker hub and/or aws ECR (works as expected, needs a README note on format)
    * Target config files are set based on localized .env file. Need to fix interface name changes in [config files.](https://github.com/magma/magma/tree/master/lte/gateway/configs) (PR needed)
    * Not required: Need to stop cloning into magma every time which over-writes any existing configuration, or provide a way to start/stop containers alone as a whole (clearly labelled flags)

## AWS Specifics

* CloudFormation
    * Individual stacks for EKS cluster
    * EFS for shared storage (supported across fargate, managed and self-managed nodegroups)
* EKS
    * Fargate compute
        * Supports all Linux workloads
        * Supports EFS for storage
        * Private subnet only
    * Managed node-groups
        * Needed for GPU compute, but AL only
        * Supports ARM
        * Supports Bottle Rocket
        * Custom AMI and CNI support
        * Supports EBS and EFS
        * Supports Daemonsets
    * Self-managed nodes
        * AWS Local Zones and Outpost can support self managed nodes only
        * Can support GPU, but AL only
        * Supports ARM
        * Supports BottleRocket
        * Supports EBS and EFS
        * Supports Daemonsets
* [Nodegroups](https://docs.aws.amazon.com/eks/latest/userguide/eks-compute.html)
* Steps
    * Create launch template for custom Ubuntu AMI that configures host with packages, OVS
