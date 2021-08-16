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

#include "include/amf_metadata_util.h"
#include "include/amf_client_servicer.h"
#include <memory>

struct amf_metadata_s amf_metadata = {};

/* Function to initialize amf metadata information */
void amf_metadata_intialize(amf_metadata_t* metadata_p) {
  metadata_p->amf_client_servicer_init();
}
