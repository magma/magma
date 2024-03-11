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
// WARNING: Do not include this header directly. Use intertask_interface.h
// instead.

// "#pragma once" will not work for this file, because this file is included
// in include/messages_def.h, which is included more than once within enum
// and structure in the file intertask_interface_types.h
// See comment in "lte/gateway/c/core/oai/include/messages_def.h" for details

MESSAGE_DEF(APPLICATION_HEALTHY_MSG, application_healthy_msg_t,
            application_healthy_msg)
MESSAGE_DEF(APPLICATION_UNHEALTHY_MSG, application_unhealthy_msg_t,
            application_unhealthy_msg)
MESSAGE_DEF(APPLICATION_MME_APP_STATS_MSG, application_mme_app_stats_msg_t,
            application_mme_app_stats_msg)
MESSAGE_DEF(APPLICATION_S1AP_STATS_MSG, application_s1ap_stats_msg_t,
            application_s1ap_stats_msg)
MESSAGE_DEF(APPLICATION_AMF_APP_STATS_MSG, application_amf_app_stats_msg_t,
            application_amf_app_stats_msg)
MESSAGE_DEF(APPLICATION_NGAP_STATS_MSG, application_ngap_stats_msg_t,
            application_ngap_stats_msg)
