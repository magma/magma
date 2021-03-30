#!/bin/bash
################################################################################
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
################################################################################

NMS_SCRIPT_URI="https://gist.githubusercontent.com/andreilee/7aa7d533e2e8f425222b1e6a016a6f5a/raw/5b3c9346dedca6dcfb8d0a13c5f3b4d9100b9279/runs-on-nms.sh"

cat << EOF
================================================================================
                  Magma Pre-1.5-Upgrade NMS DB Data Migration
================================================================================

This script will guide you through the process of copying NMS DB data to the
orc8r DB.

This process is required to be completed prior to the rest of the 1.5 upgrade
process. When this process is completed, you are ready to complete the rest of
your 1.5 upgrade.

PRE-REQUISITES
* You will want to make sure that you have some understanding of your
  kubernetes setup on which you have Magma orc8r and NMS running.
* Ensure that you have access to your kubernetes cluster which is running the
  Magma orc8r and NMS.

--------------------------------------------------------------------------------
EOF

# Get the Magma k8s namespace
read -p "Enter Kubernetes namespace [magma]: " magma_namespace
magma_namespace=${magma_namespace:-magma}
echo ""

# Find NMS pod name
nms_pod_name=$(kubectl -n magma get pods --no-headers -o custom-columns=":metadata.name" | grep nms-magmalte)
echo "Found Magma NMS pod name: $nms_pod_name"
read -p "Enter NMS pod name [$nms_pod_name]: " input
nms_pod_name=${input:-$nms_pod_name}
echo ""

# Find configurator pod name
orc8r_pod_name=$(kubectl -n magma get pods --no-headers -o custom-columns=":metadata.name" | grep orc8r-configurator)
echo "Found Magma configurator pod name: $orc8r_pod_name"
read -p "Enter configurator pod name [$orc8r_pod_name]: " input
orc8r_pod_name=${input:-$orc8r_pod_name}
echo ""

# Get DB connection parameters automatically from an orc8r pod env
# Can't do this through the NMS pod
database_source=$(kubectl -n $magma_namespace exec $orc8r_pod_name -- /bin/bash -c 'echo $DATABASE_SOURCE')
orc8r_db_host=$(echo $database_source | grep -E -o '(?:host=)\S*' | sed 's/host=//')
orc8r_db_port=$(echo $database_source | grep -E -o '(?:port=)\S*' | sed 's/port=//')
orc8r_db_name=$(echo $database_source | grep -E -o '(?:dbname=)\S*' | sed 's/dbname=//')
orc8r_db_username=$(echo $database_source | grep -E -o '(?:user=)\S*' | sed 's/user=//')
orc8r_db_password=$(echo $database_source | grep -E -o '(?:password=)\S*' | sed 's/password=//')
orc8r_db_dialect=$(kubectl -n $magma_namespace exec $orc8r_pod_name -- /bin/bash -c 'echo $SQL_DRIVER')

# Extra whitespacing
echo ""

nms_command="wget $NMS_SCRIPT_URI -O - | bash /dev/stdin -u $orc8r_db_username -w $orc8r_db_password -h $orc8r_db_host -p $orc8r_db_port -b $orc8r_db_name -d $orc8r_db_dialect --confirm"

kubectl -n $magma_namespace exec $nms_pod_name -- /bin/bash -c "$nms_command"
