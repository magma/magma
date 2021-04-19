---
id: version-1.4.0-deploy_install
title: Installing Federation Gateway
hide_title: true
original_id: deploy_install
---
# Installing Federation Gateway

## Prerequisites

To install the Federation Gateway, there are three required files that are
deployment-specific. These are described below:

* `rootCA.pem` - This file should match the `rootCA.pem` of the Orchestrator
that the Federation Gateway will connect to.

* `control_proxy.yml` - This file is used to configure the `magmad`
and `control_proxy` services to point toward the appropriate Orchestrator.
A sample configuration is provided below. The `bootstrap_address`,
`bootstrap_port`, `controller_address`, and `controller_port` are the
parameters that will likely need to be modified (check
  `/magma/feg/gateway/configs/control_proxy.yml` for the most recent
  format)

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
in the `docker-compose.yml` of the Federation Gateway. A sample configuration
is provided below (please check `magma/feg/gateway/docker/.env` for the most
  recent format):

```
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

COMPOSE_PROJECT_NAME=feg
DOCKER_REGISTRY=<registry>
DOCKER_USERNAME=<username>
DOCKER_PASSWORD=<password>
IMAGE_VERSION=latest
GIT_HASH=master

ROOTCA_PATH=/var/opt/magma/certs/rootCA.pem
CONTROL_PROXY_PATH=/etc/magma/control_proxy.yml
SNOWFLAKE_PATH=/etc/snowflake
CONFIGS_DEFAULT_VOLUME=/etc/magma
CONFIGS_TEMPLATES_PATH=/etc/magma/templates

CERTS_VOLUME=/var/opt/magma/certs
CONFIGS_VOLUME=/var/opt/magma/configs

# This section is unnecessary if using host networking
S6A_LOCAL_PORT=3868
S6A_HOST_PORT=3868
S6A_NETWORK=sctp

SWX_LOCAL_PORT=3869
SWX_HOST_PORT=3869
SWX_NETWORK=sctp

GX_LOCAL_PORT=3870
GX_HOST_PORT=3870
GX_NETWORK=tcp

GY_LOCAL_PORT=3871
GY_HOST_PORT=3871
GY_NETWORK=tcp
```

## Installation

The installation is done using the `install_gateway.sh` script located at
`magma/orc8r/tools/docker`. To install, copy that file and the three files
described above into a directory on the install host. Then

```console
INSTALL_HOST [~/]$ sudo ./install_gateway.sh feg
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
registered with the Orchestrator. At this time, NMS support for FeG
registration is still in-progress.

To register the FeG, go to the Orchestrator's Swagger UI in your browser.
(i.e. https://controller.url.sample:9443/apidocs/v1/).

Now, create a Federation Network. This is found at `/feg` under the
**Federation Networks** section. Then register the gateway under the
**Federation Gateway** section at `/feg/{network_id}/gateways` using the
network ID of the Federation Network and the hardware ID and challenge key
from the previous step.

To verify that the gateway was correctly registered, run:

```console
INSTALL_HOST [~/]$ cd /var/opt/magma/docker
INSTALL_HOST [/var/opt/magma/docker]$ docker-compose exec magmad /usr/local/bin/checkin_cli.py
```

## Upgrades

The Federation Gateway supports NMS initiated upgrades. These can be triggered
from the NMS under the `Configure` section by updating the FeG's tier to the
appropriate `Software Version`. After triggering the upgrade from the NMS,
magmad on the gateway will pull down the specified docker images,
update any static configuration, and update the docker-compose file to the
appropriate version.
