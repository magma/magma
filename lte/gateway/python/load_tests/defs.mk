# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# Add the load tests to be run
LOAD_TESTS=loadtest_mobilityd.py:allocate \
loadtest_mobilityd.py:release \
loadtest_pipelined.py:activate_flows \
loadtest_pipelined.py:deactivate_flows \
loadtest_sessiond.py:create \
loadtest_sessiond.py:end \
loadtest_subscriberdb.py:add \
loadtest_subscriberdb.py:list \
loadtest_subscriberdb.py:delete \
loadtest_subscriberdb.py:get \
loadtest_subscriberdb.py:update \
loadtest_policydb.py:enable_static_rules \
loadtest_policydb.py:disable_static_rules
