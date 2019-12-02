#!/bin/bash

set -e

if [[ -z "${MAGMA_ROOT}" ]]; then
  echo "Must set env var MAGMA_ROOT"
  exit 1
fi

export CERT_DIR=$MAGMA_ROOT/.cache/test_certs
if [[ ! -d "${CERT_DIR}" ]]; then
  echo "Certs directory is missing"
  exit 1
fi

export FILES_DIR=$MAGMA_ROOT/.circleci/devmand/test_files
export NETWORK="test_network"
export TIER="default"
export AGENT="test_agent"
export BASE_URL="https://localhost:9443/magma/v1"

# This test assumes orchestrator is running locally and that
# a symphony agent is running locally
# For curl commands, use the -k flag because we're using a self-signed cert

echo "Checking base symphony url"
curl -k -X GET "${BASE_URL}/symphony" \
      -H "accept: application/json" -H "content-type: application/json" \
      --key "${CERT_DIR}"/admin_operator.key.pem \
      --cert "${CERT_DIR}"/admin_operator.pem -v

# Create network
echo "Creating network ${NETWORK}"
curl -k -X POST "${BASE_URL}/symphony" \
      -H "accept: application/json" -H "content-type: application/json" \
      --key "${CERT_DIR}"/admin_operator.key.pem \
      --cert "${CERT_DIR}"/admin_operator.pem \
      -d @"${FILES_DIR}"/test_network_payload.json

# Create upgrade tier
echo "Creating tier ${TIER}"
curl -k -X POST "${BASE_URL}/networks/${NETWORK}/tiers" \
      -H "accept: application/json" -H "content-type: application/json" \
      --key "${CERT_DIR}"/admin_operator.key.pem \
      --cert "${CERT_DIR}"/admin_operator.pem \
      -d @"${FILES_DIR}"/default_tier_payload.json

# Create agent
echo "Creating agent ${AGENT}"
curl -k -X POST "${BASE_URL}/symphony/${NETWORK}/agents" \
      -H "accept: application/json" -H "content-type: application/json" \
      --key "${CERT_DIR}"/admin_operator.key.pem \
      --cert "${CERT_DIR}"/admin_operator.pem \
      -d @"${FILES_DIR}"/test_agent_payload.json

echo "Checking things"
curl -k -X GET "${BASE_URL}/symphony" \
      -H "accept: application/json" -H "content-type: application/json" \
      --key "${CERT_DIR}"/admin_operator.key.pem \
      --cert "${CERT_DIR}"/admin_operator.pem
curl -k -X GET "${BASE_URL}/symphony/${NETWORK}" \
      -H "accept: application/json" -H "content-type: application/json" \
      --key "${CERT_DIR}"/admin_operator.key.pem \
      --cert "${CERT_DIR}"/admin_operator.pem


# It takes up to 60 seconds for the configs to land on the box
# so instead let's just restart magmad
docker exec "$(docker ps -qf name=symphony-agent)" ls -la /var/opt/magma/configs
echo "============ WE ARE HERE 2 ============"
docker exec "$(docker ps -qf name=symphony-agent)" journalctl -n 100
docker exec "$(docker ps -qf name=symphony-agent)" systemctl restart magmad
sleep 5
echo "============ WE ARE HERE 3 ============"
docker exec "$(docker ps -qf name=symphony-agent)" journalctl -n 100
sleep 10
echo "============ WE ARE HERE 4 ============"
docker exec "$(docker ps -qf name=symphony-agent)" journalctl -n 100
sleep 10
echo "============ WE ARE HERE 5 ============"
docker exec "$(docker ps -qf name=symphony-agent)" journalctl -n 100
sleep 10
echo "============ WE ARE HERE 6 ============"
docker exec "$(docker ps -qf name=symphony-agent)" journalctl -n 100
sleep 10
echo "============ WE ARE HERE 7 ============"
docker exec "$(docker ps -qf name=symphony-agent)" journalctl -n 100
docker exec "$(docker ps -qf name=symphony-agent)" ls -la /var/opt/magma/configs
if ! docker exec "$(docker ps -qf name=symphony-agent)" test -f /var/opt/magma/configs/gateway.mconfig; then
  echo "Couldn't bootstrap successfully!"
  exit 1
fi

echo "Deleting network ${NETWORK}"
curl -k -X DELETE "${BASE_URL}/symphony/${NETWORK}" \
      -H "accept: application/json" \
      --key "${CERT_DIR}"/admin_operator.key.pem \
      --cert "${CERT_DIR}"/admin_operator.pem
