# Magma AGW Deployment (Docker) 

This guide explains how to **build and run Magma Access Gateway (AGW)** using Docker, **without modifying the official compose files**.  
It includes certificate setup, Docker build, AGW configuration, and tips for ignoring warnings.

---

##  Prerequisites

- Ubuntu 20.04 / 22.04 host  
- `sudo` privileges  
- Docker & Docker Compose plugin installed (`docker compose`)  
- Orc8r (Orchestrator) running  
- Access to rootCA certificate from Orc8r  

---

## Step 1: Generate Certificates

### Create Root Private Key
```
openssl genrsa -out rootCA.key 4096
```
### Create Self-Signed Root Certificate
```
openssl req -x509 -new -nodes -key rootCA.key -sha256 -days 1024 -out rootCA.pem
```
### Prepare Certificates Directory

sudo mkdir -p /var/opt/magma/certs
sudo cp ~/rootCA.pem /var/opt/magma/certs/
sudo cp ~/rootCA.key /var/opt/magma/certs/
ls -l /var/opt/magma/certs/
# Step 2: Build the Docker Container
Navigate to your Dockerfile directory:

```
cd ../docker/ubuntu24
```

Build the image:

```
docker build -t magma-agw-ubuntu24 .
```
## Step 3: Run the Container
Run interactively with privileged access:

```docker run -it --privileged --name magma-agw magma-agw-ubuntu24 bash ```
 Container name conflict? Remove existing container:

``` docker rm -f magma-agw
docker run -it --privileged --name magma-agw magma-agw-ubuntu24 bash
```
## Step 4: AGW Installation Inside Container
```
wget https://github.com/magma/magma/raw/v1.8/lte/gateway/deploy/agw_install_docker.sh
bash agw_install_docker.sh

reboot
```

## Step 5: Configure AGW
```
cat << EOF | sudo tee /var/opt/magma/configs/control_proxy.yml
fluentd_address: fluentd.orc8r.magmacore.link
fluentd_port: 24224

rootca_cert: /var/opt/magma/certs/rootCA.pem
EOF
```
## Step 6: Start AGW Services on Host
Use Docker Compose v2+ on the host (not inside the container):

```
cd /var/opt/magma/docker
docker compose up -d
```