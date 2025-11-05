#!/bin/bash
#
# Copyright 2025 The Magma Authors.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# Author: Nitin Rajput (coRAN LABS)
#
# eBPF GTP Installation Script for Magma AGW

set -e

MAGMA_USER="ubuntu"
MAGMA_VERSION="${MAGMA_VERSION:-ebpf-dev}"
GIT_URL="${GIT_URL:-https://github.com/magma/magma.git}"
DEPLOY_PATH="/opt/magma/lte/gateway/deploy"
EBPF_GTP_ENABLED="${EBPF_GTP_ENABLED:-true}"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' 

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

if [[ $EUID -ne 0 ]]; then
    log_error "This script must be run as root"
    exit 1
fi

if ! grep -q 'Ubuntu' /etc/issue; then
    log_error "Ubuntu is not installed"
    exit 1
fi

ROOTCA="/var/opt/magma/certs/rootCA.pem"
if [[ ! -f "$ROOTCA" ]]; then
    log_error "Upload rootCA to $ROOTCA before running this script"
    exit 1
fi

echo "=============================================="
echo "  Magma AGW Installation with eBPF GTP       "
echo "=============================================="
echo ""
log_info "Repository: $GIT_URL"
log_info "Version/Branch: $MAGMA_VERSION"
log_info "eBPF GTP Enabled: $EBPF_GTP_ENABLED"
echo ""

log_info "Verifying repository access..."
if ! git ls-remote "$GIT_URL" >/dev/null 2>&1; then
    log_error "Cannot access repository: $GIT_URL"
    log_error "Please ensure:"
    log_error "1. Repository URL is correct"
    log_error "2. Repository is public or you have access"
    log_error "3. Network connectivity is working"
    exit 1
fi
log_success "Repository access verified"

export MAGMA_VERSION="$MAGMA_VERSION"
export GIT_URL="$GIT_URL"

BACKUP_DIR="/tmp/magma_install_backup_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$BACKUP_DIR"

if [[ -d /opt/magma ]]; then
    log_info "Backing up existing Magma installation to $BACKUP_DIR"
    cp -r /opt/magma "$BACKUP_DIR/" 2>/dev/null || true
fi

log_info "Starting Magma AGW installation..."
log_info "This will take several minutes..."

if [[ ! -f "./agw_install_docker.sh" ]]; then
    log_info "Downloading installation script..."
    wget -O agw_install_docker.sh \
        "https://raw.githubusercontent.com/magma/magma/master/lte/gateway/deploy/agw_install_docker.sh"
    chmod +x agw_install_docker.sh
fi

if ./agw_install_docker.sh; then
    log_success "Magma AGW installation completed successfully"
else
    log_error "Magma AGW installation failed"
    exit 1
fi

echo ""
echo "=============================================="
echo "           Installation Summary               "
echo "=============================================="
echo ""
log_success "Magma AGW with eBPF GTP installation completed!"
echo ""
echo "Repository Used: $GIT_URL"
echo "Version/Branch: $MAGMA_VERSION"
echo "eBPF GTP Status: Enabled"
echo ""
echo "Next Steps:"
echo "1. Reboot the system: sudo reboot"
echo "Configuration Files:"
echo "- Pipelined Config: /etc/magma/pipelined.yml" 
echo "- Docker Compose: /var/opt/magma/docker/docker-compose.yaml"
echo ""
echo "Backup Location: $BACKUP_DIR"
echo ""
log_info "Installation completed successfully!"
echo "=============================================="
