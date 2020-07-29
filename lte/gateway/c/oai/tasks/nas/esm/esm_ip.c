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
#include <stdio.h>

#include "bstrlib.h"

char ip_addr_str[100];

char* esm_data_get_ipv4_addr(const_bstring ip_addr) {
  if ((ip_addr) && (ip_addr->slen == 4)) {
    snprintf(
        ip_addr_str, sizeof(ip_addr_str), "%u.%u.%u.%u", ip_addr->data[0],
        ip_addr->data[1], ip_addr->data[2], ip_addr->data[3]);
    return ip_addr_str;
  }

  return (NULL);
}

char* esm_data_get_ipv6_addr(const_bstring ip_addr) {
  if ((ip_addr) && (ip_addr->slen == 8)) {
    snprintf(
        ip_addr_str, sizeof(ip_addr_str), "%x%.2x:%x%.2x:%x%.2x:%x%.2x",
        ip_addr->data[0], ip_addr->data[1], ip_addr->data[2], ip_addr->data[3],
        ip_addr->data[4], ip_addr->data[5], ip_addr->data[6], ip_addr->data[7]);
    return ip_addr_str;
  }

  return (NULL);
}

char* esm_data_get_ipv4v6_addr(const_bstring ip_addr) {
  if ((ip_addr) && (ip_addr->slen == 12)) {
    snprintf(
        ip_addr_str, sizeof(ip_addr_str),
        "%u.%u.%u.%u / %x%.2x:%x%.2x:%x%.2x:%x%.2x", ip_addr->data[0],
        ip_addr->data[1], ip_addr->data[2], ip_addr->data[3], ip_addr->data[4],
        ip_addr->data[5], ip_addr->data[6], ip_addr->data[7], ip_addr->data[8],
        ip_addr->data[9], ip_addr->data[10], ip_addr->data[11]);
    return ip_addr_str;
  }

  return (NULL);
}
