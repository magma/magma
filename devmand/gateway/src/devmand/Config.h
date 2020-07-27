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

#pragma once

#include <gflags/gflags.h>

namespace devmand {

DECLARE_string(listen_interface);
DECLARE_string(device_configuration_file);
DECLARE_uint64(poll_interval);
DECLARE_uint64(debug_print_interval);
DECLARE_bool(devices_readonly);

} // namespace devmand
