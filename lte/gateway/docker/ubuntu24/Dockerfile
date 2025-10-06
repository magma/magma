# -----------------------------
# Magma Core AGW (LTE + 5G) Ubuntu 24.04 
# -----------------------------
FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive
ENV PATH="/opt/venv/bin:$PATH"

# -----------------------------
# Step 0: Install base dependencies
# -----------------------------
RUN apt-get update && apt-get install -y --no-install-recommends \
        sudo net-tools iproute2 bridge-utils openvswitch-switch iputils-ping tcpdump iptables \
        python3 python3-venv python3-dev python3-pip \
        curl wget git unzip make build-essential cmake pkg-config software-properties-common \
        libsystemd-dev libffi-dev libssl-dev libxml2-dev libxslt1-dev libgmp-dev zlib1g-dev rsync zip \
        ifupdown lsb-release gnupg supervisor autoconf automake libtool lksctp-tools \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# -----------------------------
# Step 1: Add Magma user
# -----------------------------
RUN useradd -ms /bin/bash magma && echo "magma ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers
WORKDIR /home/magma
VOLUME /home/magma

# -----------------------------
# Step 2: Setup Python virtual environment
# -----------------------------
RUN python3 -m venv /opt/venv && \
    /opt/venv/bin/pip install --upgrade pip setuptools wheel cython

# -----------------------------
# Step 3: Clone Magma repo
# -----------------------------
WORKDIR /magma
RUN git clone https://github.com/magma/magma.git . || true

# -----------------------------
# Step 4: Install Python dependencies
# -----------------------------
WORKDIR /magma/lte/gateway/python
RUN /opt/venv/bin/pip install --no-cache-dir grpcio requests six protobuf

# -----------------------------
# Step 5: Enable IP forwarding
# -----------------------------
RUN echo "net.ipv4.ip_forward=1" >> /etc/sysctl.conf

# -----------------------------
# Step 6: Create Magma config directories
# -----------------------------
RUN mkdir -p /etc/magma /var/opt/magma

# -----------------------------
# Step 7: Install SCTP daemon
# -----------------------------
RUN apt-get update && apt-get install -y --no-install-recommends \
        lksctp-tools libsctp-dev \
    && apt-get clean && rm -rf /var/lib/apt/lists/*
