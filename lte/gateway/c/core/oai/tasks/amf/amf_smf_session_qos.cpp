/**
 * Copyright 2022 The Magma Authors.
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

#include "lte/gateway/c/core/oai/tasks/amf/include/amf_smf_session_qos.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GQosFlowParam.hpp"

namespace magma5g {

// Fill the qos rule with filter information to be added
void amf_smf_session_api_fill_add_packet_filter(create_new_tft_t* new_tft,
                                                QOSRule* qos_rule) {
  for (int i = 0; i < qos_rule->no_of_pkt_filters; i++) {
    NewQOSRulePktFilter new_qos_rule_pkt_filter = {};
    packet_filter_t* pkt_filter = NULL;
    uint16_t pkt_filter_len = 0;
    pkt_filter = &new_tft[i];

    // Set the spare direction and id
    new_qos_rule_pkt_filter.spare = 0x0;
    new_qos_rule_pkt_filter.pkt_filter_dir = pkt_filter->direction;
    new_qos_rule_pkt_filter.pkt_filter_id = pkt_filter->identifier;
    uint16_t flag = TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG;

    // Check for the max label
    while (flag <= TRAFFIC_FLOW_TEMPLATE_MATCH_ALL_FLAG) {
      switch (pkt_filter->packetfiltercontents.flags & flag) {
        case TRAFFIC_FLOW_TEMPLATE_MATCH_ALL_FLAG: {
          // Match all type
          new_qos_rule_pkt_filter.contents[pkt_filter_len] =
              TRAFFIC_FLOW_TEMPLATE_MATCH_ALL;
          pkt_filter_len++;
        } break;

        case TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG: {
          // IPv4 remote address type
          new_qos_rule_pkt_filter.contents[pkt_filter_len] =
              TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR;
          pkt_filter_len++;

          // Copy IPV4 address and Mask
          for (int j = 0; j < TRAFFIC_FLOW_TEMPLATE_IPV4_ADDR_SIZE; j++) {
            new_qos_rule_pkt_filter.contents[pkt_filter_len] =
                pkt_filter->packetfiltercontents.ipv4remoteaddr[j].addr;
            new_qos_rule_pkt_filter
                .contents[pkt_filter_len +
                          TRAFFIC_FLOW_TEMPLATE_IPV4_ADDR_SIZE] =
                pkt_filter->packetfiltercontents.ipv4remoteaddr[j].mask;
            pkt_filter_len++;
          }

          pkt_filter_len += TRAFFIC_FLOW_TEMPLATE_IPV4_ADDR_SIZE;
        } break;
        case TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT_FLAG: {
          // Remote port type
          new_qos_rule_pkt_filter.contents[pkt_filter_len] =
              TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT;
          pkt_filter_len++;

          new_qos_rule_pkt_filter.contents[pkt_filter_len] =
              (0xFF00 & pkt_filter->packetfiltercontents.singleremoteport) >> 8;
          pkt_filter_len++;
          new_qos_rule_pkt_filter.contents[pkt_filter_len] =
              (0x00FF & pkt_filter->packetfiltercontents.singleremoteport);
          pkt_filter_len++;
        } break;
        default: {
        }
      }
      flag = flag << 1;
    }
    new_qos_rule_pkt_filter.len = pkt_filter_len;
    memcpy(&qos_rule->new_qos_rule_pkt_filter[i], &new_qos_rule_pkt_filter,
           sizeof(NewQOSRulePktFilter));
    qos_rule->len += new_qos_rule_pkt_filter.len + 2;
  }
}

// Fill the qos rule for the filter to be deleted
void amf_smf_session_api_fill_delete_packet_filter(
    delete_packet_filter_t* delete_pkt_filter, QOSRule* qos_rule) {
  NewQOSRulePktFilter new_qos_rule_pkt_filter = {};
  delete_packet_filter_t* pkt_filter = NULL;

  for (int i = 0; i < qos_rule->no_of_pkt_filters; i++) {
    pkt_filter = &delete_pkt_filter[i];
    new_qos_rule_pkt_filter.pkt_filter_id = pkt_filter->identifier;
    memcpy(qos_rule->new_qos_rule_pkt_filter, &new_qos_rule_pkt_filter,
           1 * sizeof(NewQOSRulePktFilter));
  }
}

// Convert the modification request from sessiond to modification message
// for this transaction
int amf_smf_session_api_fill_qos_ie_info(std::shared_ptr<smf_context_t> smf_ctx,
                                         bstring* authorized_qosrules,
                                         bstring* qos_flow_descriptors) {
  // Retrive message from sessoiond
  qos_flow_list_t* pti_flow_list = smf_ctx->get_proc_flow_list();

  uint8_t qos_rules_msg_buffer[QOS_RULES_MSG_BUF_LEN_MAX];
  uint16_t qos_rules_msg_buf_len = 0;

  uint8_t qos_flow_desc_buffer[QOS_FLOW_DESC_BUF_LEN_MAX];
  uint16_t qos_flow_desc_buf_len = 0;

  OAILOG_FUNC_IN(LOG_AMF_APP);

  // Prepare for sending message out to UE/GNB
  for (int i = 0; i < pti_flow_list->maxNumOfQosFlows; i++) {
    qos_flow_setup_request_item* qos_flow_req_item =
        &(pti_flow_list->item[i].qos_flow_req_item);

    // Preparing QoS Rule Msg
    if (qos_flow_req_item->ul_tft.tftoperationcode) {
      QOSRulesMsg qosRuleMsg;
      qosRuleMsg.length = 0;

      QOSRule qos_rule;
      qosRuleMsg.iei = PDU_SESSION_QOS_RULES_IE_TYPE;

      qos_rule.rule_oper_code = qos_flow_req_item->ul_tft.tftoperationcode;

      if (qos_rule.rule_oper_code ==
          TRAFFIC_FLOW_TEMPLATE_OPCODE_CREATE_NEW_TFT) {
        qos_rule.len = QOS_ADD_RULE_MIN_LEN;
      } else if (qos_rule.rule_oper_code ==
                 TRAFFIC_FLOW_TEMPLATE_OPCODE_DELETE_EXISTING_TFT) {
        qos_rule.len = QOS_DEL_RULE_MIN_LEN;
      }

      if (smf_ctx->subscribed_qos.qci ==
          qos_flow_req_item->qos_flow_identifier) {
        qos_rule.dqr_bit = QOS_RULE_DQR_BIT_SET;
        qos_rule.qos_rule_id = qos_flow_req_item->qos_flow_identifier;
      } else {
        qos_rule.dqr_bit = 0;
        qos_rule.qos_rule_id = qos_flow_req_item->qos_flow_identifier;
      }
      qos_rule.no_of_pkt_filters =
          (0x0F & qos_flow_req_item->ul_tft.numberofpacketfilters);
      qos_rule.qos_rule_precedence = 0xff;
      qos_rule.spare = 0x0;
      qos_rule.segregation = 0x0;
      qos_rule.qfi = qos_flow_req_item->qos_flow_identifier;

      // Add or Modify QOS Flow
      if (qos_flow_req_item->qos_flow_action == policy_action_add ||
          qos_flow_req_item->qos_flow_action == policy_action_del) {
        if (qos_flow_req_item->ul_tft.tftoperationcode ==
            TRAFFIC_FLOW_TEMPLATE_OPCODE_CREATE_NEW_TFT) {
          amf_smf_session_api_fill_add_packet_filter(
              qos_flow_req_item->ul_tft.packetfilterlist.createnewtft,
              &qos_rule);
        }
        qosRuleMsg.qos_rule[i] = qos_rule;
      }
      qosRuleMsg.length += qos_rule.len + QOS_RULES_MSG_MIN_LEN;

      // Convert Authorized qos into bstring
      int encoded_result = qosRuleMsg.EncodeQOSRulesMsgData(
          &qosRuleMsg, qos_rules_msg_buffer + qos_rules_msg_buf_len,
          QOS_RULES_MSG_BUF_LEN_MAX);

      if (encoded_result < 0) {
        OAILOG_ERROR(LOG_AMF_APP,
                     "Qos Rule parameters invalid or un-aligned \n");
        OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
      }
      qos_rules_msg_buf_len += encoded_result;
    }  // qos rule creation

    // Preparing qos flow descriptors
    if (qos_flow_req_item->qos_flow_descriptor.qos_flow_identifier) {
      M5GQosFlowDescription flow_des = {};
      flow_des.numOfParams = 0;

      qos_flow_descriptor_t* qos_flow_desc =
          &qos_flow_req_item->qos_flow_descriptor;
      flow_des.operationCode = qos_flow_req_item->ul_tft.tftoperationcode << 5;
      flow_des.qfi = qos_flow_desc->qos_flow_identifier;

      // Set  fiveqi flow descriptor
      if (qos_flow_desc->fiveQi) {
        flow_des.paramList[flow_des.numOfParams].iei =
            magma5g::M5GQosFlowParam::param_id_5qi;
        flow_des.paramList[flow_des.numOfParams].length = sizeof(uint8_t);
        flow_des.paramList[flow_des.numOfParams].element =
            qos_flow_desc->fiveQi;
        flow_des.numOfParams++;
      }

      // Set mbr dl
      if (qos_flow_desc->mbr_dl) {
        flow_des.paramList[flow_des.numOfParams].iei =
            magma5g::M5GQosFlowParam::param_id_mfbr_downlink;
        flow_des.paramList[flow_des.numOfParams].length = 3;
        M5GQosFlowParam* qosParams = &flow_des.paramList[flow_des.numOfParams];
        qosParams->mfbr_gbr_convert(qosParams, qos_flow_desc->mbr_dl);
        flow_des.numOfParams++;
      }

      // Set mbr ul
      if (qos_flow_desc->mbr_ul) {
        flow_des.paramList[flow_des.numOfParams].iei =
            magma5g::M5GQosFlowParam::param_id_mfbr_uplink;
        flow_des.paramList[flow_des.numOfParams].length = 3;
        M5GQosFlowParam* qosParams = &flow_des.paramList[flow_des.numOfParams];
        qosParams->mfbr_gbr_convert(qosParams, qos_flow_desc->mbr_ul);
        flow_des.numOfParams++;
      }

      // Set gbr dl
      if (qos_flow_desc->gbr_dl) {
        flow_des.paramList[flow_des.numOfParams].iei =
            magma5g::M5GQosFlowParam::param_id_gfbr_downlink;
        flow_des.paramList[flow_des.numOfParams].length = 3;
        M5GQosFlowParam* qosParams = &flow_des.paramList[flow_des.numOfParams];
        qosParams->mfbr_gbr_convert(qosParams, qos_flow_desc->gbr_dl);
        flow_des.numOfParams++;
      }

      // Set gbr ul
      if (qos_flow_desc->gbr_ul) {
        flow_des.paramList[flow_des.numOfParams].iei =
            magma5g::M5GQosFlowParam::param_id_gfbr_uplink;
        flow_des.paramList[flow_des.numOfParams].length = 3;
        M5GQosFlowParam* qosParams = &flow_des.paramList[flow_des.numOfParams];
        qosParams->mfbr_gbr_convert(qosParams, qos_flow_desc->gbr_ul);
        flow_des.numOfParams++;
      }

      if (flow_des.numOfParams > 0) {
        flow_des.Ebit = 1;
      } else {
        flow_des.Ebit = 0;
      }

      if (flow_des.numOfParams) {
        // Convert Authorized qos into bstring
        int encoded_result = flow_des.EncodeM5GQosFlowDescription(
            &flow_des, qos_flow_desc_buffer + qos_flow_desc_buf_len,
            QOS_FLOW_DESC_BUF_LEN_MAX);

        if (encoded_result < 0) {
          OAILOG_ERROR(LOG_AMF_APP,
                       "qos flow Description invalid or un-aligned \n");
          OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
        }
        qos_flow_desc_buf_len += encoded_result;
      }
    }  // flow description operation
  }

  // Fill the autorized qos rule buffer
  *authorized_qosrules = blk2bstr(qos_rules_msg_buffer, qos_rules_msg_buf_len);

  // Fill the flow descriptor
  *qos_flow_descriptors = blk2bstr(qos_flow_desc_buffer, qos_flow_desc_buf_len);

  OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNok);
}

// No qos rules received from SMF. Create a default one
void amf_smf_session_set_default_qos_rule(qos_flow_list_t* pti_flow_list) {
  qos_flow_setup_request_item* qos_flow_req_item =
      &(pti_flow_list->item[0].qos_flow_req_item);

  // Default Qos Already present return
  if (qos_flow_req_item->ul_tft.tftoperationcode ==
      TRAFFIC_FLOW_TEMPLATE_OPCODE_CREATE_NEW_TFT) {
    return;
  }

  if (!qos_flow_req_item->qos_flow_identifier) {
    qos_flow_req_item->qos_flow_identifier = PDU_SESSION_DEFAULT_QFI;
  }
  qos_flow_req_item->qos_flow_action = policy_action_add;
  qos_flow_req_item->ul_tft.tftoperationcode =
      TRAFFIC_FLOW_TEMPLATE_OPCODE_CREATE_NEW_TFT;
  qos_flow_req_item->ul_tft.numberofpacketfilters = 1;
  qos_flow_req_item->ul_tft.packetfilterlist.createnewtft[0].direction =
      TRAFFIC_FLOW_TEMPLATE_BIDIRECTIONAL;
  qos_flow_req_item->ul_tft.packetfilterlist.createnewtft[0].identifier = 0x1;
  qos_flow_req_item->ul_tft.packetfilterlist.createnewtft[0]
      .packetfiltercontents.flags = TRAFFIC_FLOW_TEMPLATE_MATCH_ALL_FLAG;
}

// For setting the default qos if nothing frmm sessiond
void amf_smf_session_set_default_qos_info(
    std::shared_ptr<smf_context_t> smf_ctx) {
  qos_flow_list_t* pti_flow_list = smf_ctx->get_proc_flow_list();

  if (!pti_flow_list->maxNumOfQosFlows) {
    pti_flow_list->maxNumOfQosFlows = 1;
  }

  // Default QoS Rules
  amf_smf_session_set_default_qos_rule(pti_flow_list);
}

}  // namespace magma5g
