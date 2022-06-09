#!/bin/bash

# Copyright 2022 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

MME_state=$(state_cli.py keys IMSI | grep -c MME)
S1AP_state=$(state_cli.py keys IMSI | grep -c S1AP)
SPGW_state=$(state_cli.py keys IMSI | grep -c SPGW)
Sessiond_state=$(state_cli.py parse  sessiond:sessions|grep -c '\"create_session_response\":')
table_0_dl_bearers=$(ovs-ofctl dump-flows gtp_br0 table=0|grep -c "nw_dst")
(( table_0_dl_bearers = table_0_dl_bearers/2 ))
table_0_paging_rules=$(ovs-ofctl dump-flows gtp_br0 table=0|grep -c CONTROLLER)
table_12_dl_flows=$(ovs-ofctl dump-flows gtp_br0 table=12|grep -c "nw_dst")
table_13_dl_flows=$(ovs-ofctl dump-flows gtp_br0 table=13|grep -c "nw_dst")
table_12_ul_flows=$(ovs-ofctl dump-flows gtp_br0 table=12|grep -c "nw_src")
table_13_ul_flows=$(ovs-ofctl dump-flows gtp_br0 table=13|grep -c "nw_src")
mobility_state=$(mobility_cli.py get_subscriber_table | grep -c -v SID)
(( idle_ues = table_0_paging_rules - table_0_dl_bearers))
echo "MME_state = $MME_state"
echo "S1AP_state = $S1AP_state"
echo "SPGW_state = $SPGW_state"
echo "Sessiond_state = $Sessiond_state"
if (( table_0_dl_bearers > 0 )); then
  echo "table0 DL bearers = $table_0_dl_bearers"
fi
if (( idle_ues > 0 )); then
  echo "table0 Idle UEs = $idle_ues"
fi
echo "table12 DL flows = $table_12_dl_flows ; table12 UL flows = $table_12_ul_flows"
echo "table13 DL flows = $table_13_dl_flows ; table13 UL flows = $table_13_ul_flows"
echo "Mobilityd state = $mobility_state"
