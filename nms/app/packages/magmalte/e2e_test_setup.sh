#!/usr/bin/env bash
#
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e
mkdir -p .cache
mkdir -p /tmp/nms_artifacts

openssl req -nodes -new -x509 -batch -keyout .cache/mock_server.key -out .cache/mock_server.cert
docker-compose --env-file .env.mock -f docker-compose-e2e.yml up -d

i=0
while [ $i -lt 60 ]
do
    val=$(docker-compose -f docker-compose-e2e.yml logs magmalte 2>&1 | grep 'Production server started on port 8081' | wc -l)
    if [ $val -eq 1 ]
    then
        break
    fi
    echo "magmalte server not started yet...sleeping 1 second"
    sleep 1
    i=$[$i+1]
done

docker-compose exec magmalte yarn setAdminPassword magma-test admin@magma.test password1234
docker-compose exec magmalte yarn createOrganization magma-test nms test,test_feg_lte_network

# run the end to end test
cd ../../
set +e
yarn test:e2e
exit_code=$?

cd packages/magmalte
docker-compose -f docker-compose-e2e.yml logs magmalte &> /tmp/nms_artifacts/magmalte.log
docker-compose -f docker-compose-e2e.yml logs mock_server &> /tmp/nms_artifacts/mock_server.log

exit $exit_code
