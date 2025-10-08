# -----------------------------
# Stage 1: Builder (Ubuntu 24.04)
# -----------------------------
FROM ubuntu:24.04 AS builder

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
    ifupdown lsb-release gnupg supervisor autoconf automake libtool lksctp-tools libsctp-dev \
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
# Stage 2: Python base image for services
# -----------------------------
FROM python:3.10-slim AS python_base

ENV PYTHONUNBUFFERED=1
WORKDIR /opt/magma

# -----------------------------
# Copy built files from builder
# -----------------------------
COPY --from=builder /opt/venv /opt/venv
COPY --from=builder /magma /opt/magma
COPY --from=builder /usr/local/bin /usr/local/bin

ENV PATH="/opt/venv/bin:$PATH"

# -----------------------------
# Install runtime dependencies
# -----------------------------
RUN apt-get update && apt-get install -y --no-install-recommends \
    redis-server curl netcat iproute2 iputils-ping iptables lksctp-tools \
    && rm -rf /var/lib/apt/lists/*

# -----------------------------
# Dynamic entrypoint for all Magma services
# -----------------------------
COPY <<'EOF' /entrypoint.sh
#!/bin/bash
set -e

echo "Starting Magma AGW Service: ${SERVICE_NAME}"

case "$SERVICE_NAME" in
  magmad)
    exec python3 -m magma.magmad.main
    ;;
  redis)
    exec redis-server --daemonize no
    ;;
  directoryd)
    exec python3 -m magma.directoryd.main
    ;;
  subscriberdb)
    exec python3 -m magma.subscriberdb.main
    ;;
  enodebd)
    exec python3 -m magma.enodebd.main
    ;;
  state)
    exec python3 -m magma.state.main
    ;;
  policydb)
    exec python3 -m magma.policydb.main
    ;;
  health)
    exec python3 -m magma.health.main
    ;;
  monitord)
    exec python3 -m magma.monitord.main
    ;;
  redirectd)
    exec python3 -m magma.redirectd.main
    ;;
  smsd)
    exec python3 -m magma.smsd.main
    ;;
  control_proxy)
    exec nghttpx --conf /var/opt/magma/tmp/nghttpx.conf
    ;;
  ctraced)
    exec python3 -m magma.ctraced.main
    ;;
  sctpd)
    exec /usr/local/bin/sctpd
    ;;
  oai_mme)
    exec /usr/local/bin/oai_mme -c /var/opt/magma/tmp/mme.conf -s /var/opt/magma/tmp/spgw.conf
    ;;
  pipelined)
    exec python3 -m magma.pipelined.main
    ;;
  sessiond)
    exec /usr/local/bin/sessiond
    ;;
  mobilityd)
    exec python3 -m magma.mobilityd.main
    ;;
  td-agent-bit)
    exec /opt/td-agent-bit/bin/td-agent-bit -c /var/opt/magma/tmp/td-agent-bit.conf
    ;;
  eventd)
    exec python3 -m magma.eventd.main
    ;;
  connectiond)
    exec /usr/local/bin/connectiond
    ;;
  liagentd)
    exec /usr/local/bin/liagentd
    ;;
  envoy_controller)
    exec /var/opt/magma/bin/envoy_controller
    ;;
  *)
    echo "Unknown service: ${SERVICE_NAME}"
    exit 1
    ;;
esac
EOF

RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

# -----------------------------
# Wait-for-healthy script
# -----------------------------
COPY <<'EOF' /wait_for_healthy.sh
#!/bin/bash
set -e

SERVICES=(
  "magmad:50052"
  "redis:6379"
  "directoryd:50067"
  "subscriberdb:50051"
  "enodebd:60055"
  "state:50074"
  "policydb:50068"
  "health:50080"
  "monitord:50076"
  "redirectd:50071"
  "smsd:50078"
  "control_proxy:8443"
  "ctraced:50079"
  "sctpd:/tmp/sctpd_downstream.sock"
  "oai_mme:50073"
  "pipelined:50063"
  "sessiond:50065"
  "mobilityd:60051"
  "td-agent-bit:5140"
  "eventd:50075"
  "connectiond:50082"
  "envoy_controller:50081"
)

echo "🕒 Waiting for all Magma AGW services to become healthy..."

for svc in "${SERVICES[@]}"; do
  name="${svc%%:*}"
  port="${svc##*:}"
  echo "Checking $name ($port)..."

  if [[ "$port" == /* ]]; then
    until test -S "$port"; do sleep 3; done
  else
    until nc -zv localhost "$port" >/dev/null 2>&1; do sleep 3; done
  fi

  echo "✅ $name is healthy."
done

echo "🎉 All services are up!"
EOF

RUN chmod +x /wait_for_healthy.sh

# -----------------------------
# Volumes
# -----------------------------
VOLUME ["/var/opt/magma/certs", "/etc/magma/templates", "/etc/magma",
        "/var/opt/magma/configs", "/var/opt/magma/tmp", "/etc/snowflake",
        "/var/opt/magma/fluent-bit", "/var/opt/magma/docker", "/var/run",
        "/tmp", "/var/log", "/etc/openvswitch", "/var/run/docker.sock"]

WORKDIR /var/opt/magma
