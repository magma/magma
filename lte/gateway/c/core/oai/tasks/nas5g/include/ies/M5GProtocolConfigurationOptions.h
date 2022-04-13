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
#include <sstream>
#include <cstdint>
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.008.h"
#ifdef __cplusplus
};
#endif

namespace magma5g {

// Protocol Configuration Options IE Class
class ProtocolConfigurationOptions {
 public:
  protocol_configuration_options_t pco;

  ProtocolConfigurationOptions();
  ~ProtocolConfigurationOptions();

  int m5g_encode_protocol_configuration_options(
      const protocol_configuration_options_t* const pco, uint8_t* buffer,
      const uint32_t len);

  int EncodeProtocolConfigurationOptions(ProtocolConfigurationOptions* pco,
                                         uint8_t iei, uint8_t* buffer,
                                         uint32_t len);

  int m5g_decode_protocol_configuration_options(
      protocol_configuration_options_t* pco, const uint8_t* const buffer,
      const uint32_t len);

  int DecodeProtocolConfigurationOptions(ProtocolConfigurationOptions* pco,
                                         uint8_t iei, uint8_t* buffer,
                                         uint32_t len);
};
}  // namespace magma5g
