#!/usr/bin/env bash

# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


#sudo password
PASSWD="magma"

#Installing Ruby as it is prerequisite for FluentD
install_ruby() {
	echo "***** Installing Ruby as it is prerequisite for FluentD*****"
	{ sleep 0.1; echo '%s'; } | script -q -c 'sudo apt-get update' /dev/null % $1 
	{ sleep 0.1; echo '%s'; } | script -q -c 'sudo apt install -y git curl libssl-dev libreadline-dev zlib1g-dev autoconf bison build-essential libyaml-dev libreadline-dev libncurses5-dev libffi-dev libgdbm-dev python3-prometheus-client python3-flask' /dev/null % $1
	curl -sL https://github.com/rbenv/rbenv-installer/raw/master/bin/rbenv-installer | bash -
	echo 'export PATH="$HOME/.rbenv/bin:$PATH"' >> ~/.bashrc
	echo 'eval "$(rbenv init -)"' >> ~/.bashrc
	echo 'export FLASK_APP="/home/magma/agw_state.py"' >> ~/.bashrc
	{ sleep 0.1; echo '%s'; } | script -q -c 'sudo mkdir -p /var/log/alert_manager' /dev/null % $1 
        echo 'export prometheus_multiproc_dir="/var/log/alert_manager"' >> ~/.bashrc	
	source ~/.bashrc
	rbenv install 2.7.2
	rbenv global 2.7.2
}

#Setting prerequisites like increasing the Max no.of File Descriptors, adding configurations to /etc/sysctl.conf file
set_prerequisites() {
	echo "*****Increase the Maximum Number of File Descriptors*****"
	#Append below lines to /etc/security/limits.conf file to increase Maximum Number of File Descriptors
	cat <<-EOT >> /etc/security/limits.conf
	root soft nofile 65536 
	root hard nofile 65536 
	* soft nofile 65536 
	* hard nofile 65536
	EOT

	#Optimize the Network Kernel Parameters by adding the below configuration to /etc/sysctl.conf file
	cat <<-EOT >> /etc/sysctl.conf
	net.core.somaxconn = 1024 
	net.core.netdev_max_backlog = 5000 
	net.core.rmem_max = 16777216 
	net.core.wmem_max = 16777216 
	net.ipv4.tcp_wmem = 4096 12582912 16777216 
	net.ipv4.tcp_rmem = 4096 12582912 16777216 
	net.ipv4.tcp_max_syn_backlog = 8096 
	net.ipv4.tcp_slow_start_after_idle = 0 
	net.ipv4.tcp_tw_reuse = 1 
	net.ipv4.ip_local_port_range = 10240 65535
	EOT

	echo "*****Running sysctl -p command or reboot your node for the changes to take effect*****"
	sysctl -p
}

#Downloading and Running AGW install on Bare Metal VM
agw_install() {
	echo  "*****Downloading and Running AGW install on Bare Metal VM*****"
	wget https://raw.githubusercontent.com/facebookincubator/magma/v1.3/lte/gateway/deploy/agw_install.sh
	bash agw_install.sh
}

install() {
	install_ruby PASSWD
	set_prerequisites
	agw_install
}

usage() {
	echo "Usage : bash $0 -i"
	echo "      (OR)   "
	echo "Usage : bash $0 --install"
	exit 2
}

#Script starts from here
param="$1"
if [[ ($param == "-i") || ($param == "--install") ]]; then
	install
else
	usage
fi
