/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/*****************************************************************************

  Source      nas_network.h

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#ifndef NAS5G_NETWORK_SEEN
#define NAS5G_NETWORK_SEEN

#include <sstream>
#include "amf_config.h"
#ifdef __cplusplus
extern "C" {
#endif
#include "bstrlib.h"
#ifdef __cplusplus
};
#endif
using namespace std;

namespace magma5g {
class nas_network {
 public:
  int nas_network_initialize(amf_config_t);
  void bdestroy_wrapper(bstring* b);
  void free_wrapper(void** ptr);
};
}  // namespace magma5g

#endif
