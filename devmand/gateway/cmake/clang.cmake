# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS}" "-std=c++2a")
list(APPEND CMAKE_CXX_WARNINGS
  "-fcolor-diagnostics"
  "-Weverything"
  "-Wno-weak-vtables"
  "-Wno-\#pragma-messages"
  "-Wno-unknown-attributes"
  "-Wno-address-of-packed-member"
  "-Wno-c++98-compat"
  "-Wno-c++98-compat-pedantic"
  "-Wno-padded"
  "-Wno-packed"
  "-Wno-disabled-macro-expansion"
  "-Weffc++"
  "-Werror"
  "-pedantic")
