/*
   Copyright 2021 The Magma Authors.
   This source code is licensed under the BSD-style license found in the
   LICENSE file in the root directory of this source tree.
   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
 */

#pragma once
namespace magma5g {
// RequestType Class
class RequestType {
 public:
  uint8_t iei : 4;
  uint32_t type_val : 3;

  RequestType();
  ~RequestType();
  int EncodeRequestType(RequestType* reqest_type, uint8_t iei, uint8_t* buffer,
                        uint32_t len);
  int DecodeRequestType(RequestType* reqest_type, uint8_t iei, uint8_t* buffer,
                        uint32_t len);
  void copy(const RequestType& p) {
    iei = p.iei;
    type_val = p.type_val;
  }
  bool isEqual(const RequestType& p) {
    return ((iei == p.iei) && (type_val == p.type_val));
  }
};
}  // namespace magma5g
