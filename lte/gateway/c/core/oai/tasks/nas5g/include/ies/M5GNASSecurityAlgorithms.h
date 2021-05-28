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

using namespace std;
namespace magma5g {
class NASSecurityAlgorithmsMsg {
 public:
  uint8_t tca;
  uint8_t tia;
#define NAS_SECURITY_ALGORITHMS_MINIMUM_LENGTH 1

  typedef struct M5GNasSecurityAlgorithms_tag {
#define M5G_NAS_SECURITY_ALGORITHMS_5G_EA0 0b000
#define M5G_NAS_SECURITY_ALGORITHMS_128_5G_EA1 0b001
#define M5G_NAS_SECURITY_ALGORITHMS_128_5G_EA2 0b010
#define M5G_NAS_SECURITY_ALGORITHMS_128_5G_EA3 0b011
#define M5G_NAS_SECURITY_ALGORITHMS_5G_EA4 0b100
#define M5G_NAS_SECURITY_ALGORITHMS_5G_EA5 0b101
#define M5G_NAS_SECURITY_ALGORITHMS_5G_EA6 0b110
#define M5G_NAS_SECURITY_ALGORITHMS_5G_EA7 0b111
    uint8_t m5gtypeofcipheringalgorithm : 3;
#define M5G_NAS_SECURITY_ALGORITHMS_5G_IA0 0b000
#define M5G_NAS_SECURITY_ALGORITHMS_128_5G_IA1 0b001
#define M5G_NAS_SECURITY_ALGORITHMS_128_5G_IA2 0b010
#define M5G_NAS_SECURITY_ALGORITHMS_128_5G_IA3 0b011
#define M5G_NAS_SECURITY_ALGORITHMS_5G_EIA4 0b100
#define M5G_NAS_SECURITY_ALGORITHMS_5G_EIA5 0b101
#define M5G_NAS_SECURITY_ALGORITHMS_5G_EIA6 0b110
#define M5G_NAS_SECURITY_ALGORITHMS_5G_EIA7 0b111
    uint8_t m5gtypeofintegrityalgorithm : 3;
  } M5GNasSecurityAlgorithms;

  M5GNasSecurityAlgorithms M5GNasSecurityAlgorithms_;
  NASSecurityAlgorithmsMsg();
  ~NASSecurityAlgorithmsMsg();
  int EncodeNASSecurityAlgorithmsMsg(
      NASSecurityAlgorithmsMsg* nas_sec_algorithms, uint8_t iei,
      uint8_t* buffer, uint32_t len);
  int DecodeNASSecurityAlgorithmsMsg(
      NASSecurityAlgorithmsMsg* nas_sec_algorithms, uint8_t iei,
      uint8_t* buffer, uint32_t len);
};
}  // namespace magma5g
