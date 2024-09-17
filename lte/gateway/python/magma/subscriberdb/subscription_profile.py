"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

from lte.protos.mconfig.mconfigs_pb2 import SubscriberDB


def get_default_sub_profile(service):
    """
    Returns the default subscription profile to be used when the
    subcribers don't have a profile associated with them.
    """
    if 'default' in service.mconfig.sub_profiles:
        return service.mconfig.sub_profiles['default']
    # No default profile configured for the network. Use the default defined
    # in the code.
    return SubscriberDB.SubscriptionProfile(
        max_ul_bit_rate=service.config['default_max_ul_bit_rate'],
        max_dl_bit_rate=service.config['default_max_dl_bit_rate'],
    )
