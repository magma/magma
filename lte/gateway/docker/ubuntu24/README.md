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
```
sudo mkdir -p /var/opt/magma/certs
sudo cp ~/rootCA.pem /var/opt/magma/certs/
sudo cp ~/rootCA.key /var/opt/magma/certs/
ls -l /var/opt/magma/certs/
```
## Step 2: Clone Magma Repository
Navigate to your Dockerfile directory:

```
cd $HOME
git clone https://github.com/magma/magma.git
cd magma/lte/gateway/docker
```
## Step 3: Run the Container
Step 3: Build AGW Docker Images (Optional)

If you want to build AGW images locally instead of pulling official images:
```
docker compose build gateway_c gateway_python
```
### 3.1 Tag the images for reference
```
docker tag agw-gateway_c:latest linuxfoundation.jfrog.io/magma-docker/agw_gateway_c:1.9.0
docker tag agw-gateway_python:latest linuxfoundation.jfrog.io/magma-docker/agw_gateway_python:1.9.0
```


## Step 4: Configure AGW
```
cat << EOF | sudo tee /var/opt/magma/configs/control_proxy.yml

fluentd_address: fluentd.orc8r.magmacore.link
fluentd_port: 24224

rootca_cert: /var/opt/magma/certs/rootCA.pem
EOF
```
## Step 5: Start AGW Services on Host
Use Docker Compose v2+ on the host (not inside the container):

```
cd /var/opt/magma/docker
docker compose up -d
```