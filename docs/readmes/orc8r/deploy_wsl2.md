---
id: deploy_wsl2
title: Deploy on WSL2 (Windows)
hide_title: true
---

# Deploying Magma on WSL2 (Windows)

This guide walks through setting up a Magma development environment on
Windows using WSL2 (Windows Subsystem for Linux). This is an alternative
to AWS or bare metal deployment, suitable for local development and testing.

Before following this guide, it may be useful to read through the
[Magma prerequisites](../basics/prerequisites.md) and
[quick start guide](../basics/quick_start_guide.md) sections.

## Prerequisites

- Windows 10 version 2004+ or Windows 11
- At least 8GB RAM (16GB recommended)
- 50GB free disk space
- Administrator access to Windows
- Stable internet connection

## 1. WSL2 Setup

### Install WSL2

Open PowerShell as Administrator:

```powershell
wsl --install -d Ubuntu-24.04
```

Restart your computer when prompted. Ubuntu will launch automatically after
restart. Create a username and password when prompted.

### Verify WSL2 Installation

```powershell
wsl --list --verbose
```

Ensure Ubuntu shows `VERSION 2`. If not, upgrade with:

```powershell
wsl --set-version Ubuntu-24.04 2
```

### Configure WSL2 Resources (Recommended)

Create `C:\Users\<YourUsername>\.wslconfig` with:

```ini
[wsl2]
memory=8GB
processors=4
swap=2GB
```

Restart WSL after making changes:

```powershell
wsl --shutdown
wsl
```

## 2. Docker Installation

### Install Docker Engine

Inside the Ubuntu terminal:

```bash
sudo apt update && sudo apt upgrade -y
sudo apt install docker.io -y
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker $USER
newgrp docker
```

### Verify Docker

```bash
docker --version
docker run hello-world
```

You should see `Hello from Docker!` if installation was successful.

### Install Docker Compose

```bash
sudo apt install docker-compose -y
docker-compose --version
```

:::note
Magma's build scripts use the `--compatibility` flag which was removed in
Docker v28+. If you encounter `unknown flag: --compatibility`, use
`docker-compose` (with hyphen) instead of `docker compose` (with space).
:::

## 3. Clone Magma Repository

```bash
export MAGMA_ROOT=~/magma
git clone https://github.com/magma/magma.git $MAGMA_ROOT
cd $MAGMA_ROOT
```

Add the following to your `~/.bashrc` so the variable persists:

```bash
echo 'export MAGMA_ROOT=~/magma' >> ~/.bashrc
source ~/.bashrc
```

## 4. Orchestrator Deployment

### Build Images

```bash
cd $MAGMA_ROOT/orc8r/cloud/docker

# Create build context
sudo rm -rf /tmp/magma_orc8r_build
python3 build.py --all
```

### Start Services

```bash
docker-compose -f docker-compose.yml \
  -f docker-compose.metrics.yml \
  -f docker-compose.override.yml up -d
```

### Verify Services

```bash
docker ps
```

You should see containers for `controller`, `postgres`, `nginx`, and others.

## 5. Testing with UERANSIM (Alternative)

If full orchestrator deployment is not required, UERANSIM provides a
lightweight 5G UE/gNB simulator that works well on WSL2.

### Install UERANSIM

```bash
# Install dependencies
sudo apt install make g++ libsctp-dev lksctp-tools iproute2 cmake -y

# Clone and build
cd ~
git clone https://github.com/aligungr/UERANSIM
cd UERANSIM
make
```

### Verify Build

```bash
ls build/
# Expected: nr-gnb  nr-ue  nr-cli  nr-binder  libdevbnd.so
```

### Basic Test

Open two terminals:

**Terminal 1 - Start gNB:**

```bash
cd ~/UERANSIM
./build/nr-gnb -c config/open5gs-gnb.yaml
```

**Terminal 2 - Connect UE:**

```bash
cd ~/UERANSIM
./build/nr-ue -c config/open5gs-ue.yaml
```

## 6. Troubleshooting

### Network Connectivity Issues

**Symptom:** `Connection failed`, TLS timeouts, or `apt` errors

**Solution:**

```bash
# Update DNS resolvers
echo "nameserver 8.8.8.8" | sudo tee /etc/resolv.conf
echo "nameserver 8.8.4.4" | sudo tee -a /etc/resolv.conf
```

Then restart WSL from PowerShell:

```powershell
wsl --shutdown
wsl
```

### Permission Denied (Docker)

**Symptom:** `permission denied while trying to connect to Docker daemon`

**Solution:**

```bash
sudo usermod -aG docker $USER
newgrp docker
```

### Build Context Missing

**Symptom:** `build path /tmp/magma_orc8r_build does not exist`

**Solution:**

```bash
cd $MAGMA_ROOT/orc8r/cloud/docker
sudo rm -rf /tmp/magma_orc8r_build
python3 build.py --all
```

### Slow Performance

Store the Magma repository inside the Linux filesystem (`~/magma`), not
the Windows filesystem (`/mnt/c/...`). WSL2 has significantly slower I/O
when accessing Windows filesystem paths.

### Copy-Paste in Ubuntu Terminal

Enable copy-paste by right-clicking the terminal title bar →
**Properties** → check **Use Ctrl+Shift+C/V as Copy/Paste**.

## 7. Additional Resources

- [Magma Documentation](https://magma.github.io/magma/)
- [WSL2 Documentation](https://learn.microsoft.com/en-us/windows/wsl/)
- [UERANSIM](https://github.com/aligungr/UERANSIM)
- [Magma Slack](https://magma-community.slack.com) - `#new-to-magma`---
id: deploy_wsl2
title: Deploy on WSL2 (Windows)
hide_title: true
---
