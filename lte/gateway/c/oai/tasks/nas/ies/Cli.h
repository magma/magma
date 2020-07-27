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

#ifndef CLI_SEEN
#define CLI_SEEN

#include <stdint.h>

#include "bstrlib.h"

#define CLI_MINIMUM_LENGTH 3
#define CLI_MAXIMUM_LENGTH 14

typedef bstring Cli;

int encode_cli(Cli cli, uint8_t iei, uint8_t* buffer, uint32_t len);

int decode_cli(Cli* cli, uint8_t iei, uint8_t* buffer, uint32_t len);

#endif /* CLI_SEEN */
