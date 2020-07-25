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

#define LOG_WITH_GLOG
#include <devmand/channels/cli/Spd2Glog.h>
#include <magma_logging.h>
#include <spdlog/details/log_msg.h>

void devmand::channels::cli::Spd2Glog::_sink_it(
    const spdlog::details::log_msg& msg) {
  toGlog(msg);
}

void devmand::channels::cli::Spd2Glog::toGlog(
    const spdlog::details::log_msg& msg) {
  if (msg.level == trace || msg.level == debug) {
    MLOG(MDEBUG) << msg.formatted.str();
  } else if (msg.level == info || msg.level == warn) {
    MLOG(MINFO) << msg.formatted.str();
  } else {
    MLOG(MERROR) << msg.formatted.str();
  }
}

void devmand::channels::cli::Spd2Glog::flush() {}
