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
/****************************************************************************
  Source      ngap_amf_decoder.h
  Version     0.1
  Date        2020/07/28
  Product     NGAP stack
  Subsystem   Access and Mobility Management Function
  Author      Ashish Prajapati
  Description Defines NG Application Protocol Messages

*****************************************************************************/

#ifndef FILE_NGAP_AMF_DECODER_SEEN
#define FILE_NGAP_AMF_DECODER_SEEN
#include "bstrlib.h"
#include "ngap_common.h"
#include "intertask_interface_types.h"

int ngap_amf_decode_pdu(Ngap_NGAP_PDU_t* pdu, const_bstring const raw)
    __attribute__((warn_unused_result));

// int ngap_amf_decode_pdu(ngap_message *message,  const bstring const raw,
// MessagesIds *messages_id); int ngap_free_amf_decode_pdu(ngap_message
// *message, MessagesIds messages_id);

#endif /* FILE_NGAP_AMF_DECODER_SEEN */
