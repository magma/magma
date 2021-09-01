/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#include <string>

#include "sctp_defs.h"
#include "includes/MConfigLoader.h"
#include "lte/protos/mconfig/mconfigs.pb.h"

namespace magma {
namespace mme {

#define MME_SERVICE "mme"

std::string upstream_sock_from_mconfig() {
  magma::mconfig::MME mconfig;
  magma::MConfigLoader loader;
  if (!loader.load_service_mconfig(MME_SERVICE, &mconfig)) {
    return UPSTREAM_SOCK;
  }
  if (mconfig.upstream_sctp_sock().empty()) {
    return UPSTREAM_SOCK;
  }
  return mconfig.upstream_sctp_sock();
}

std::string downstream_sock_from_mconfig() {
  magma::mconfig::MME mconfig;
  magma::MConfigLoader loader;
  if (!loader.load_service_mconfig(MME_SERVICE, &mconfig)) {
    return DOWNSTREAM_SOCK;
  }
  if (mconfig.downstream_sctp_sock().empty()) {
    return DOWNSTREAM_SOCK;
  }
  return mconfig.downstream_sctp_sock();
}

}  // namespace mme
}  // namespace magma
