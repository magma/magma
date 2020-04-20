#!/usr/bin/env bash
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# Install and start dynamoDB
sudo mkdir -p /var/tmp/archives
sudo wget -O /var/tmp/archives/dynamo.zip https://s3-us-west-2.amazonaws.com/dynamodb-local/dynamodb_local_latest.zip
sudo unzip /var/tmp/archives/dynamo.zip -d /var/tmp/archives/dynamo
sudo java -Djava.library.path=/var/tmp/archives/dynamo/DynamoDBLocal_lib -jar /var/tmp/archives/dynamo/DynamoDBLocal.jar -dbPath /var/tmp/archives/dynamo -sharedDb &
