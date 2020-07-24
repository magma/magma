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

#ifndef PLMN_LIST_H_
#define PLMN_LIST_H_
#include <stdint.h>

#define PLMN_LIST_MINIMUM_LENGTH 5
#define PLMN_LIST_MAXIMUM_LENGTH 47

typedef struct PlmnList_tag {
  uint8_t mccdigit2 : 4;
  uint8_t mccdigit1 : 4;
  uint8_t mncdigit3 : 4;
  uint8_t mccdigit3 : 4;
  uint8_t mncdigit2 : 4;
  uint8_t mncdigit1 : 4;
} PlmnList;

int encode_plmn_list(
    PlmnList* plmnlist, uint8_t iei, uint8_t* buffer, uint32_t len);

int decode_plmn_list(
    PlmnList* plmnlist, uint8_t iei, uint8_t* buffer, uint32_t len);

void dump_plmn_list_xml(PlmnList* plmnlist, uint8_t iei);

#endif /* PLMN LIST_H_ */
