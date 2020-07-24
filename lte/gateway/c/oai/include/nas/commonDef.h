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

/*

Source      commonDef.h

Version     0.1

Date        2012/02/27

Product     NAS stack

Subsystem   include

Author      Frederic Maurel

Description Contains global common definitions

*****************************************************************************/
#ifndef FILE_COMMONDEF_SEEN
#define FILE_COMMONDEF_SEEN

#include <stdint.h>

/*
 * A list of PLMNs
 */
#define PLMN_LIST_T(SIZE)                                                      \
  struct {                                                                     \
    uint8_t n_plmns;                                                           \
    plmn_t plmn[SIZE];                                                         \
  }

/*
 * A list of TACs
 */
#define TAC_LIST_T(SIZE)                                                       \
  struct {                                                                     \
    uint8_t n_tacs;                                                            \
    TAC_t tac[SIZE];                                                           \
  }

/*
 * A list of TAIs
 */
#define TAI_LIST_T(SIZE)                                                       \
  struct {                                                                     \
    uint8_t list_type;                                                         \
    uint8_t n_tais;                                                            \
    tai_t tai[SIZE];                                                           \
  }

#endif /* FILE_COMMONDEF_SEEN*/
