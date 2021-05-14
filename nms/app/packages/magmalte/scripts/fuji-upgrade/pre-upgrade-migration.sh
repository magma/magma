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

NMS_SCRIPT_URI="https://raw.githubusercontent.com/magma/magma/master/nms/app/packages/magmalte/scripts/fuji-upgrade/runs-on-nms.sh"

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
* You are running Magma orchestrator v1.3 or v1.4.
* You will want to make sure that you have some understanding of your
  kubernetes setup on which you have Magma orc8r and NMS running.
* Ensure that you have access to your kubernetes cluster which is running the
  Magma orc8r and NMS.

--------------------------------------------------------------------------------
EOF

# Get the Magma k8s namespace
read -r -p "Enter Kubernetes namespace [orc8r]: " magma_namespace
magma_namespace=${magma_namespace:-orc8r}
echo ""

# Find NMS pod name
nms_pod_name="$(kubectl -n $magma_namespace get pod -l app.kubernetes.io/component=magmalte -o jsonpath='{.items[0].metadata.name}')"
echo "Found Magma NMS pod name: $nms_pod_name"
read -r -p "Enter NMS pod name [$nms_pod_name]: " input
nms_pod_name="${input:-$nms_pod_name}"
echo ""

# Find configurator/controller pod name
orc8r_pod_name="$(kubectl -n $magma_namespace get pod -l app.kubernetes.io/component=configurator -o jsonpath='{.items[0].metadata.name}')"
if [ -n "$orc8r_pod_name" ]; then
  echo "Found Magma configurator pod name: $orc8r_pod_name"
  read -r -p "Enter configurator pod name [$orc8r_pod_name]: " input
else
  orc8r_pod_name="$(kubectl -n $magma_namespace get pod -l app.kubernetes.io/component=controller -o jsonpath='{.items[0].metadata.name}')"
  echo "Found Magma controller pod name: $orc8r_pod_name"
  read -r -p "Enter controller pod name [$orc8r_pod_name]: " input
fi
orc8r_pod_name="${input:-$orc8r_pod_name}"
echo ""

# Get DB connection parameters automatically from an orc8r pod env
# Can't do this through the NMS pod
database_source=$(kubectl -n $magma_namespace exec $orc8r_pod_name -- /bin/bash -c 'echo $DATABASE_SOURCE')
orc8r_db_host=$(echo $database_source | awk -F= 'BEGIN { RS=" "; } /host/ { print $2; }')
orc8r_db_port=$(echo $database_source | awk -F= 'BEGIN { RS=" "; } /port/ { print $2; }')
orc8r_db_name=$(echo $database_source | awk -F= 'BEGIN { RS=" "; } /dbname/ { print $2; }')
orc8r_db_username=$(echo $database_source | awk -F= 'BEGIN { RS=" "; } /user/ { print $2; }')
orc8r_db_password=$(echo $database_source | awk -F= 'BEGIN { RS=" "; } /password/ { print $2; }')
orc8r_db_dialect=$(kubectl -n $magma_namespace exec $orc8r_pod_name -- /bin/bash -c 'echo $SQL_DRIVER')

# Extra whitespacing
echo ""

nms_command="wget $NMS_SCRIPT_URI -O - | bash /dev/stdin -u $orc8r_db_username -w $orc8r_db_password -h $orc8r_db_host -p $orc8r_db_port -b $orc8r_db_name -d $orc8r_db_dialect --confirm"

kubectl -n $magma_namespace exec $nms_pod_name -- /bin/bash -c "$nms_command"
