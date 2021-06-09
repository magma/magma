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
  Date        2020/07/28
  Subsystem   Access and Mobility Management Function
  Author      Ashish Prajapati
  Description Defines NG Application Protocol Messages

*****************************************************************************/
#pragma once

#include "bstrlib.h"
#include "ngap_common.h"
#include "intertask_interface_types.h"

int ngap_amf_decode_pdu(Ngap_NGAP_PDU_t* pdu, const_bstring const raw)
    __attribute__((warn_unused_result));
