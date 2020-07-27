/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

#include <devmand/Config.h>

namespace devmand {

DEFINE_string(listen_interface, "eth0", "The interface to listen on.");
DEFINE_string(
    device_configuration_file,
    "/etc/devmand/devices.yml",
    "Accepts .yml or .mconfig files. Inotify watches the file, and applies necessary changes.");
DEFINE_uint64(poll_interval, 55, "The polling interval in seconds.");
DEFINE_uint64(
    debug_print_interval,
    0,
    "The debug print interval in seconds. A value of 0 disables the printing.");
DEFINE_bool(
    devices_readonly,
    false,
    "whether or not devices can be configured");

} // namespace devmand
