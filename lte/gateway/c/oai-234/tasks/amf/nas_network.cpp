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

  Source      nas_network.cpp

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#include <thread>
#ifdef __cplusplus
extern "C" {
#endif

#include "bstrlib.h"
#include "log.h"
#ifdef __cplusplus
}
#endif
#include "nas5g_network.h"
using namespace std;
namespace magma5g {
//------------------------------------------------------------------------------
void nas_network::free_wrapper(void** ptr) {
  // for debug only
  // AssertFatal(ptr, "Trying to free NULL ptr");
  if (ptr) {
    delete (ptr);
    *ptr = NULL;
  }
}

//------------------------------------------------------------------------------
void nas_network::bdestroy_wrapper(bstring* b) {
  if ((b) && (*b)) {
    bdestroy(*b);
    *b = NULL;
  }
}
}  // namespace magma5g
