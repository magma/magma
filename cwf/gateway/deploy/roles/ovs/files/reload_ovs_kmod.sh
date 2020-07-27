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

SRC_VERSION_FN=/sys/module/openvswitch/srcversion
if test -f "$SRC_VERSION_FN"; then
    echo "Checking if upgrade is necessary"
    LOADED_VER=$(cat /sys/module/openvswitch/srcversion)
    INSTALLED_VER=$(modinfo -F srcversion openvswitch)
    echo "Loaded Version: $LOADED_VER, Installed Version:$INSTALLED_VER"
    if [ "$LOADED_VER" != "$INSTALLED_VER" ]; then
        echo "Version mismatch, Reloading openvswitch kernel module"
        /usr/share/openvswitch/scripts/ovs-ctl force-reload-kmod
    else
        echo "Skipping reload, version match"
    fi
else
    echo "No src version file found, Reloading openvswitch kernel module"
    /usr/share/openvswitch/scripts/ovs-ctl load-kmod
fi
