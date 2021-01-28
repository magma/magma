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

#ifndef FILE_MME_APP_IF_SEEN
#define FILE_MME_APP_IF_SEEN

#define SGsAP_PAGING_REQUEST 0b00000001              // 1
#define SGsAP_PAGING_REJECT 0b00000010               // 2
#define SGsAP_SERVICE_REQUEST 0b00000110             // 6
#define SGsAP_DOWNLINK_UNITDATA 0b00000111           // 7
#define SGsAP_UPLINK_UNITDATA 0b00001000             // 8
#define SGsAP_LOCATION_UPDATE_REQUEST 0b00001001     // 9
#define SGsAP_LOCATION_UPDATE_ACCEPT 0b00001010      // 10
#define SGsAP_LOCATION_UPDATE_REJECT 0b00001011      // 11
#define SGsAP_TMSI_REALLOCATION_COMPLETE 0b00001100  // 12
#define SGsAP_ALERT_REQUEST 0b00001101               // 13
#define SGsAP_ALERT_ACK 0b00001110                   // 14
#define SGsAP_ALERT_REJECT 0b00001111                // 15
#define SGsAP_UE_ACTIVITY_INDICATION 0b00010000      // 16
#define SGsAP_EPS_DETACH_INDICATION 0b00010001       // 17
#define SGsAP_EPS_DETACH_ACK 0b00010010              // 18
#define SGsAP_IMSI_DETACH_INDICATION 0b00010011      // 19
#define SGsAP_IMSI_DETACH_ACK 0b00010100             // 20
#define SGsAP_RESET_INDICATION 0b00010101            // 21
#define SGsAP_RESET_ACK 0b00010110                   // 22
#define SGsAP_SERVICE_ABORT_REQUEST 0b00010111       // 23
#define SGsAP_MO_CSFB_INDICATION 0b00011000          // 24
#define SGsAP_MM_INFORMATION_REQUEST 0b00011010      // 26
#define SGsAP_RELEASE_REQUEST 0b00011011             // 27
#define SGsAP_STATUS 0b00011101                      // 29
#define SGsAP_UE_UNREACHABLE 0b00011111              // 31

#endif /* FILE_MME_APP_IF_SEEN */
