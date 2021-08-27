#!/bin/bash
###############################################################################
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
###############################################################################

# Script arguments (specifies source DB only):
# -u: DB username
# -w: DB password
# -h: DB host
# -p: DB port
# -b: DB name
# -d: DB SQL dialect
#
# EXAMPLE USAGE
# ./runs-on-nms.sh -h postgres -p 5432 -b nms -d postgres -u root -w password

while getopts ":h:p:b:d:u:w:" opt; do
  case $opt in
    h) arg_host="$OPTARG"
    ;;
    p) arg_port="$OPTARG"
    ;;
    b) arg_db="$OPTARG"
    ;;
    d) arg_dialect="$OPTARG"
    ;;
    u) arg_username="$OPTARG"
    ;;
    w) arg_password="$OPTARG"
    ;;
    \?) echo "Invalid option -$OPTARG" >&2
    ;;
  esac
done

cat << EOF
--------------------------------------------------------------------------------
              Upgrading @fbcnms/sequelize-models to ^0.1.9
--------------------------------------------------------------------------------
EOF
pushd /usr/src/packages/magmalte
yarn upgrade @fbcnms/sequelize-models@^0.1.9
yarn
popd
cat << EOF
--------------------------------------------------------------------------------
              Successfully upgraded @fbcnms/sequelize-models
--------------------------------------------------------------------------------
EOF

cat << EOF
--------------------------------------------------------------------------------
              Running yarn script for NMS DB data migration
--------------------------------------------------------------------------------
EOF
pushd /usr/src/node_modules/\@fbcnms/sequelize-models
yarn dbDataMigrate --host=$arg_host --port=$arg_port --database=$arg_db --dialect=$arg_dialect --username=$arg_username --password=$arg_password --export --confirm
if [ $? -eq 0 ]
then
  echo "SUCCESS"
  echo "You are now ready to upgrade to 1.5"
  exit 0
else
  echo "Failed to migrate NMS DB data: " >&2
  exit 1
fi
