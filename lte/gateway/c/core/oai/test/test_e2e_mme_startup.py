#!/usr/bin/env python
#
# Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The OpenAirInterface Software Alliance licenses this file to You under
# the terms found in the LICENSE file in the root of this source tree.
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# -------------------------------------------------------------------------------
# For more information about the OpenAirInterface (OAI) Software Alliance:
#      contact@openairinterface.org
#

from mme_app_driver import MMEAppDriver


def test_non_blocking_mme_init():
    """Test that MME startup does not block during init phase.

    Tests very specifically for s6a init running without blocking MME startup.
    """
    log_conditions = (  # In regex form
        r'Initializing S6a interface',  # S6A init is running
        r'S6a peer connection attempt \d+ / \d+',  # S6A attempting to connect
        r'MME app initialization complete',  # MME proceeded past init steps
    )
    MMEAppDriver().run(log_conditions=log_conditions)


def main():
    """Main method for testing."""
    test_non_blocking_mme_init()


if __name__ == '__main__':
    main()
