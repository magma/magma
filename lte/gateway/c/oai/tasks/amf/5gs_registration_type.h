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
#include <sstream>
#include <thread>
#include "../bstr/bstrlib.h"
using namespace std;
#pragma once
#define AMF_REGISTRATION_TYPE_MINIMUM_LENGTH 1
#define AMF_REGISTRATION_TYPE_MAXIMUM_LENGTH 1

#define AMF_REGISTRATION_TYPE_INITIAL 0b0001
#define AMF_REGISTRATION_TYPE_MOBILITY_UPDATING 0b0010
#define AMF_REGISTRATION_TYPE_PERODIC_UPDATING 0b0011
#define AMF_REGISTRATION_TYPE_EMERGENCY 0b0100
#define AMF_REGISTRATION_TYPE_RESERVED 0b0111

    

namespace magma5g
{


}
/*
0 0 1 initial registration
0 1 0 mobility registration updating
0 1 1 periodic registration updating
1 0 0 emergency registration
1 1 1 reserved
*/