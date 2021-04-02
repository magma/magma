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

install_td_agent() {
	#gem install fluentd
	echo "*****gem install and fluentd-gem install fluentd-plugin-elasticsearch*****"
	gem install fluentd --no-doc
	fluentd --setup ./fluent
	fluent-gem install fluent-plugin-elasticsearch
	echo "*****Installing Treasure Data(td-agent)*****"
	#td-agent 3
	curl -L https://toolbelt.treasuredata.com/sh/install-debian-stretch-td-agent3.sh | sh
	{ sleep 0.1; echo '%s'; } | script -q -c 'sudo systemctl start td-agent.service' /dev/null % $1
	{ sleep 0.1; echo '%s'; } | script -q -c 'sudo /etc/init.d/td-agent restart' /dev/null % $1
	echo "*****Add Configuration to /etc/td-agent/td-agent.conf file*****"
	cat <<-EOT >> /etc/td-agent/td-agent.conf
	<source>
	  @type tail  
	  format none  
	  secure false  
	  path /var/log/syslog  
	  pos_file /var/log/fluentd/pos/syslog.log.pos  
	  read_from_head true  
	  tag syslog  
	</source> 
	<source>  
 	  @type tail  
	  format none  
	  secure false  
	  path /var/log/mme.log  
	  pos_file /var/log/fluentd/pos/mme.log.pos  
	  read_from_head true  
	  tag magma.log  
	</source> 
	<source>  
	  @type tail  
	  format none  
	  secure false  
	  path /var/log/enodebd.log  
	  pos_file /var/log/fluentd/pos/enodebd.log.pos  
	  read_from_head true  
	  tag magma.log  
	</source> 
	<source>  
	  @type tail  
	  format none  
	  secure false  
	  path /var/log/daemon.log  
	  pos_file /var/log/fluentd/pos/daemon.log.pos  
	  read_from_head true  
	  tag magma.log  
	</source> 
	<source>  
  	  @type tail  
	  format none  
  	  secure false  
  	  path /var/log/user.log  
  	  pos_file /var/log/fluentd/pos/user.log.pos  
 	  read_from_head true  
 	  tag magma.log  
 	</source> 
	<source>  
	  @type tail  
 	  format none  
 	  secure false  
 	  path /var/log/kern.log  
 	  pos_file /var/log/fluentd/pos/kern.log.pos  
 	  read_from_head true  
	  tag magma.log  
	</source> 
	<source>  
 	  @type tail  
 	  format none  
 	  secure false  
 	  path /var/log/dpkg.log  
 	  pos_file /var/log/fluentd/pos/dpkg.log.pos  
 	  read_from_head true  
 	  tag magma.log  
	</source> 
	<source>  
  	  @type tail  
 	  format none  
  	  secure false  
 	  path /var/log/redis/redis-server.log  
 	  pos_file /var/log/fluentd/pos/redis-server.log.pos  
 	  read_from_head true  
  	  tag magma.log  
	</source> 
	<match {syslog,**.log}> 
 	  @type elasticsearch 
 	  host 10.9.9.46 #change me to your host ip-address
 	  port 31714 #change me to your elasticsearch service port number under kubevirt namespace
 	  scheme http 
 	  logstash_format true 
	  secure false 
 	  include_timestamp true 
 	  logstash_prefix fluentd 
 	  <buffer> 
 	    @type file 
	    path /var/log/td-agent/buffer/elasticsearch 
	  </buffer> 
	  <secondary> 
 	    @type secondary_file 
 	    directory /var/log/td-agent/error 
 	    basename syslog 
 	  </secondary> 
	</match>
	EOT
	echo "*****changing /var/log group to td-agent to access logs to EFK*****"
	{ sleep 0.1; echo '%s'; } | script -q -c 'sudo mkdir -p /var/log/fluentd/pos' /dev/null % $1
	{ sleep 0.1; echo '%s'; } | script -q -c 'sudo chmod 777 -R /var/log/fluentd/pos' /dev/null % $1
	cd /var/log/
	umask 001
	chgrp td-agent .
	chmod g+s .
	chgrp td-agent mme.log enodebd.log daemon.log syslog user.log kern.log dpkg.log
	cd redis/
	chgrp td-agent .
	chmod g+s .
	chgrp td-agent redis-server.log
 	/etc/init.d/td-agent restart
	echo "*****Calling update_agw_state as back ground process to update the stats to AlertManger*****"
	cd /home/magma/
	python3 update_agw_state.py &
	echo "*****Run flask with magmadev VM IP and port=5000*****"
	flask run --host=<ip of magmadev VM> --port=5000 &
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
	install_td_agent PASSWD
else
	usage
fi
