# Magma AGW + Orc8r (NMS) Deployment Guide

## Overview

This guide walks through deploying the Magma LTE core network using virtual machines. It covers:

* Orc8r (Orchestrator) + NMS setup
* Access Gateway (AGW) setup
* Networking configuration
* Troubleshooting and verification

---

## Architecture

Magma consists of three main components:

* **Orc8r (Orchestrator)**
  Central control plane managing gateways, configs, and APIs

* **NMS (Network Management System)**
  Web UI for managing networks and gateways

* **AGW (Access Gateway)**
  Connects LTE base stations to the core network

---

## Deployment Environment

* VirtualBox VMs
* Ubuntu 20.04

### VM Setup

| Component   | OS     | Purpose            |
| ----------- | ------ | ------------------ |
| Orc8r + NMS | Ubuntu | Control plane + UI |
| AGW         | Ubuntu | Gateway            |

### Networking

* Adapter 1 → NAT (internet)
* Adapter 2 → Host-Only (VM communication)

---

## Step 1: Clone Magma

```bash
git clone https://github.com/magma/magma.git
cd magma
```

Verify:

```bash
ls
```

---

## Step 2: Setup Orc8r (Controller)

```bash
cd orc8r/cloud/docker
```

### Install Docker

(Follow official Docker install guide)

---

### Build Services

```bash
./build.py --all
```

### Fix Permission Error

```bash
sudo usermod -aG docker $USER
reboot
```

---

### Start Services

```bash
docker compose up -d
```

Verify:

```bash
docker compose ps
```

---

## Step 3: Setup NMS

### Set Admin Credentials

```bash
docker compose exec magmalte yarn setAdminPassword host admin@magma.test password1234
```

---

### Access UI

```
http://host.localhost:8081/host
```

---

### If NMS Fails

Reset:

```bash
docker compose down --volumes --rmi all --remove-orphans
docker compose up -d
```

Re-set password:

```bash
docker compose exec magmalte yarn setAdminPassword host admin@magma.test password1234
```

---

## Step 4: Setup AGW

### Install AGW

```bash
wget https://github.com/magma/magma/raw/v1.8/lte/gateway/deploy/agw_install_docker.sh
bash agw_install_docker.sh
```

---

### Add Root Certificate

On Orc8r:

```bash
cat ~/magma/.cache/test_certs/rootCA.pem
```

On AGW:

```bash
sudo mkdir -p /var/opt/magma/certs
sudo nano /var/opt/magma/certs/rootCA.pem
```

---

### Create Config File

```bash
sudo mkdir -p /var/opt/magma/configs
```

```bash
cat << EOF | sudo tee /var/opt/magma/configs/control_proxy.yml
cloud_address: <ORC8R_IP>
cloud_port: 443
bootstrap_address: <ORC8R_IP>
bootstrap_port: 443
fluentd_address: <ORC8R_IP>
fluentd_port: 24224
rootca_cert: /var/opt/magma/certs/rootCA.pem
EOF
```

---

## Step 5: Fix Networking (CRITICAL)

### Problem

NAT alone does NOT allow reliable VM communication.

---

### Solution: Host-Only Network

* Add Adapter 2 → Host-Only
* Restart both VMs

---

### Get Orc8r IP

```bash
ip a
```

Example:

```
192.168.56.102
```

---

### Update AGW Hosts File

```bash
sudo nano /etc/hosts
```

Add:

```
192.168.56.102 controller.magma.test
192.168.56.102 bootstrapper-controller.magma.test
```

---

## Step 6: Start AGW

```bash
cd /var/opt/magma/docker
docker compose up -d
```

---

## Step 7: Register Gateway

```bash
docker exec magmad show_gateway_info.py
```

Use output in NMS:

* Hardware ID
* UUID

---

## Step 8: Verify Connection

```bash
docker exec magmad checkin_cli.py
```

Expected:

* TCP connection ✅
* SSL success ✅
* Cloud check-in ✅

---

## Troubleshooting

### Fluentd Issue

Disable in `docker-compose.yml`:

```yaml
# fluentd:
#   ...
```

---

### No Controller Connectivity

Check:

```bash
ping controller.magma.test
```

---

### No IP on VM

Fix Netplan:

```bash
sudo nano /etc/netplan/*.yaml
```

Enable:

```yaml
dhcp4: true
```

Apply:

```bash
sudo netplan apply
```

---

## Final Result

* Orc8r running ✅
* NMS accessible ✅
* AGW connected ✅
* Gateway registered ✅

---

## Next Steps

* Integrate LTE RAN (e.g., srsRAN)
* Test subscriber connectivity
