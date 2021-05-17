---
id: deploy_install
title: Installing Carrier Wifi Gateway
hide_title: true
---
# Installing Carrier Wifi Gateway

## Prerequisites

To install the Carrier Wifi Gateway, there are three required files that are
deployment-specific. These are described below:

* `rootCA.pem` - This file should match the `rootCA.pem` of the Orchestrator
that the Carrier Wifi Gateway will connect to.

* `control_proxy.yml` - This file is used to configure the `magmad`
and `control_proxy` services to point toward the appropriate Orchestrator.
A sample configuration is provided below. The `bootstrap_address`,
`bootstrap_port`, `controller_address`, and `controller_port` are the
parameters that will likely need to be modified.
 
```
#
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# nghttpx config will be generated here and used
nghttpx_config_location: /var/tmp/nghttpx.conf

# Location for certs
rootca_cert: /var/opt/magma/certs/rootCA.pem
gateway_cert: /var/opt/magma/certs/gateway.crt
gateway_key: /var/opt/magma/certs/gateway.key

# Listening port of the proxy for local services. The port would be closed
# for the rest of the world.
local_port: 8443

# Cloud address for reaching out to the cloud.
cloud_address: controller.magma.test
cloud_port: 443

bootstrap_address: bootstrapper-controller.magma.test
bootstrap_port: 443

# Option to use nghttpx for proxying. If disabled, the individual
# services would establish the TLS connections themselves.
proxy_cloud_connections: True

# Allows http_proxy usage if the environment variable is present
allow_http_proxy: True
```

* `.env` - This file provides any deployment specific environment variables used
in the `docker-compose.yml` of the Carrier Wifi Gateway. A sample configuration
is provided below:

```
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

COMPOSE_PROJECT_NAME=cwf
DOCKER_REGISTRY=<registry>
DOCKER_USERNAME=<username>
DOCKER_PASSWORD=<password>
IMAGE_VERSION=latest
GIT_HASH=master

BUILD_CONTEXT=https://github.com/facebookincubator/magma.git#master

ROOTCA_PATH=/var/opt/magma/certs/rootCA.pem
CONTROL_PROXY_PATH=/etc/magma/control_proxy.yml
CONFIGS_TEMPLATES_PATH=/etc/magma/templates

CERTS_VOLUME=/var/opt/magma/certs
CONFIGS_OVERRIDE_VOLUME=/var/opt/magma/configs
CONFIGS_DEFAULT_VOLUME=/etc/magma
```

## Installation

The installation is done using the `install_gateway.sh` script located at
`magma/orc8r/tools/docker`. To install, copy that file and the three files
described above into a directory on the install host. Then

```console
INSTALL_HOST [~/]$ sudo ./install_gateway.sh cwag
```

After this completes, you should see: `Installed successfully!!`

## Registration

After installation, the next step is to register the gateway with the Orchestrator.
To do so:

```console
INSTALL_HOST [~/]$ cd /var/opt/magma/docker
INSTALL_HOST [/var/opt/magma/docker]$ docker-compose exec magmad /usr/local/bin/show_gateway_info.py
```

This will output a hardware ID and a challenge key. This information must be
registered with the Orchestrator.

To register the Carrier Wifi Gateway, go to the Orchestrator's APIdocs in your browser. 
**Note: It is highly encouraged to use V1 of the apidocs**
(i.e. https://controller.url.sample:9443/apidocs/v1/).

Now, create a Carrier Wifi Network. This is found at `/cwf` under the
**Carrier Wifi Networks** section. Then register the gateway under the
**Carrier Wifi Gateway** section at `/cwf/{network_id}/gateways` using the 
network ID of the Carrier Wifi Network and the hardware ID and challenge key 
from the previous step.

To verify that the gateway was correctly registered, run:

```console
INSTALL_HOST [~/]$ cd /var/opt/magma/docker
INSTALL_HOST [/var/opt/magma/docker]$ docker-compose exec magmad /usr/local/bin/checkin_cli.py
```

## Upgrades

The Carrier Wifi Gateway supports NMS initiated upgrades. These can be triggered
from the NMS under the `Configure` section by updating the CWF's tier to the
appropriate `Software Version`. After triggering the upgrade from the NMS,
magmad on the gateway will pull down the specified docker images,
update any static configuration, and update the docker-compose file to the
appropriate version.

### Prepare Gateway Node for Upgrade:

1. Configuring Docker With a Proxy

 In order to the set the proxy for Docker, you will need to create a configuration file for the Docker service. No configuration files exist by default, so one will have to be created.

```console 
a. Create a new directory for Docker service configurations
  
  sudo mkdir -p /etc/systemd/system/docker.service.d

b. Create a file called http-proxy.conf in configuration directory. 

  sudo vim /etc/systemd/system/docker.service.d/http-proxy.conf

c. Add the following contents, changing the values to match your environment

  [Service]
  Environment=HTTP_PROXY=http://bootstrapper-controller.magma.com:443
  Environment=HTTPS_PROXY=https://bootstrapper-controller.magma.com:443

d. Save your changes and Reload the daemon configuration.

  sudo systemctl daemon-reload

e. Install rootCA cert on ubuntu machine.

  sudo cp /var/opt/magma/certs/rootCA.pem /usr/local/share/ca-certificates/rootCA.crt
  sudo update-ca-certificates

f. Restart Docker to apply changes.

  sudo service docker restart
```

2. update orc8r to support proxy

```console
a. update orc8r-proxy values by editing  vals.yml
  
proxy:
  spec:
    http_proxy_docker_hostname: "docker.io"
    http_proxy_github_hostname: "github.com"
        
b. upgrade helm deployment 

  cd magma/orc8r/cloud/helm/orc8r
  helm upgrade orc8r . --values=PATH_TO_VALS/vals.yml
  kubectl -n magma get pods
```

3. create or update upgrade tier with latest tag/commit id

```console
a. open Orchestrator's APIdocs in your browser

 https://<orc8r_ip>:9443/apidocs/v1/#/Upgrades/post_networks__network_id__tiers

 {
   "gateways": [
     "cwf01"
   ],
   "id": "stable",
   "images": [
     {
       "name": "string",
       "order": 0
     }
   ],
   "name": "Stable Tier",
   "version": "1.0.0-123456789-<commit_id/tag_id>" 
 }
```

4. tail magmad logs on the gateway to see the upgrade status. 
```console
  [gateway]$ cd /var/opt/magma/docker
	[gateway]$ docker-compose logs -f magmad
```
