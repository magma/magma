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

  Source      5gs_registration_type.h

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#ifndef GS5_REGISTRATION_TYPE_SEEN
#define GS5_REGISTRATION_TYPE_SEEN

#include <sstream>
#include <thread>
#include "bstrlib.h"
using namespace std;
#define AMF_REGISTRATION_TYPE_MINIMUM_LENGTH 1
#define AMF_REGISTRATION_TYPE_MAXIMUM_LENGTH 1
#if 0
#define AMF_REGISTRATION_TYPE_INITIAL_ 0x01
#define AMF_REGISTRATION_TYPE_MOBILITY_UPDATING_ 0x02
#define AMF_REGISTRATION_TYPE_PERODIC_UPDATING_ 0x03
#define AMF_REGISTRATION_TYPE_EMERGENCY_ 0x04
#define AMF_REGISTRATION_TYPE_RESERVED_ 0x07
#endif

namespace magma5g {}

#endif
/*
0 0 1 initial registration
0 1 0 mobility registration updating
0 1 1 periodic registration updating
1 0 0 emergency registration
1 1 1 reserved
*/
