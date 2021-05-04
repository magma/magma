/*
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

#pragma once

//------------------------------------------------------------------------------
// 6.2 5G Identifiers
//------------------------------------------------------------------------------

typedef uint32_t
    gnb_ue_ngap_id_t; /*!< \brief  An gNB UE NGAP ID shall be allocated so as to
                       uniquely identify the UE over the N2 interface within an
                       gNB. When an AMF receives an gNB UE NGAP ID it shall
                       store it for the duration of the UE-associated logical
                       N2-connection for this UE. Once known to an AMF this IE
                       is included in all UE associated N2-AP signalling. The
                       gNB UE NGAP ID shall be unique within the gNB logical
                       node. */

typedef uint32_t
    amf_ue_ngap_id_t; /*!< \brief  A AMF UE NGAP ID shall be allocated so as to
                       uniquely identify the UE over the N2 interface within the
                       AMF. When an gNB receives AMF UE NGAP ID it shall store
                       it for the duration of the UE-associated logical
                       N2-connection for this UE. Once known to an gNB this IE
                       is included in all UE associated N2-AP signalling. The
                       AMF UE NGAP ID shall be unique within the AMF logical
                       node.*/
typedef uint32_t
    ran_ue_ngap_id_t;  // This IE uniquely identifies the UE association over
                       // the NG interface within the NG-RAN node
