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

#include "proto_msg_to_itti_msg.h"

#include <stdint.h>
#include <string.h>
#include <iostream>
#include <string>

#include "bstrlib.h"
#include "3gpp_23.003.h"
#include "3gpp_24.008.h"
#include "common_ies.h"
#include "feg/protos/csfb.pb.h"
#include "lte/protos/sms_orc8r.pb.h"

extern "C" {
#include "bytes_to_ie.h"
}

namespace magma {
using namespace lte;

// SMS Orc8r Downlink
void convert_proto_msg_to_itti_sgsap_downlink_unitdata(
    const SMODownlinkUnitdata* msg, itti_sgsap_downlink_unitdata_t* itti_msg) {
  std::string imsi = msg->imsi();
  // If north bound is Orc8r itself, IMSI prefix is used;
  // in AGW local tests, IMSI prefix is not used
  // Strip off any IMSI prefix
  if (imsi.compare(0, 4, "IMSI") == 0) {
    imsi = imsi.substr(4, std::string::npos);
  }
  itti_msg->imsi_length = imsi.size();
  strcpy(itti_msg->imsi, imsi.c_str());

  auto nas_msg = msg->nas_message_container();
  if (nas_msg.length() > 0) {
    itti_msg->nas_msg_container =
        bfromcstr_with_str_len(nas_msg.c_str(), nas_msg.length());
  }

  return;
}

}  // namespace magma

namespace magma {
using namespace feg;
void convert_proto_msg_to_itti_sgsap_eps_detach_ack(
    const EPSDetachAck* msg, itti_sgsap_eps_detach_ack_t* itti_msg) {
  auto imsi             = msg->imsi();
  itti_msg->imsi_length = imsi.length();
  strcpy(itti_msg->imsi, imsi.c_str());

  return;
}

void convert_proto_msg_to_itti_sgsap_imsi_detach_ack(
    const IMSIDetachAck* msg, itti_sgsap_imsi_detach_ack_t* itti_msg) {
  auto imsi             = msg->imsi();
  itti_msg->imsi_length = imsi.length();
  strcpy(itti_msg->imsi, imsi.c_str());

  return;
}

void convert_proto_msg_to_itti_sgsap_location_update_accept(
    const LocationUpdateAccept* msg,
    itti_sgsap_location_update_acc_t* itti_msg) {
  auto imsi             = msg->imsi();
  itti_msg->imsi_length = imsi.length();
  strcpy(itti_msg->imsi, imsi.c_str());

  auto lai = msg->location_area_identifier();
  lai_t itti_lai;
  bytes_to_lai(lai.c_str(), &itti_lai);
  itti_msg->laicsfb = itti_lai;

  if (msg->newIMSITMSI_case() == LocationUpdateAccept::kNewImsi) {
    itti_msg->presencemask            = SGSAP_MOBILE_IDENTITY;
    itti_msg->mobileid.typeofidentity = MOBILE_IDENTITY_IMSI;
    auto new_imsi                     = msg->new_imsi();
    itti_msg->mobileid.length         = new_imsi.length();
    strcpy(itti_msg->mobileid.u.imsi, new_imsi.c_str());
  } else if (msg->newIMSITMSI_case() == LocationUpdateAccept::kNewTmsi) {
    auto new_tmsi = msg->new_tmsi();
    if (new_tmsi.length() != TMSI_SIZE) {
      std::cout << "[MWARNING] "
                << "Expected length of new TMSI in Location Update Accept: "
                << new_tmsi.length() << ", got " << TMSI_SIZE
                << " instead. Ignoring the TMSI" << std::endl;
      itti_msg->presencemask = 0;
      return;
    }
    itti_msg->presencemask            = SGSAP_MOBILE_IDENTITY;
    itti_msg->mobileid.typeofidentity = MOBILE_IDENTITY_TMSI;
    itti_msg->mobileid.length         = new_tmsi.length();
    for (int i = 0; i < TMSI_SIZE; ++i) {
      itti_msg->mobileid.u.tmsi[i] = new_tmsi[i];
    }
  } else {
    itti_msg->presencemask = 0;
  }

  return;
}

void convert_proto_msg_to_itti_sgsap_location_update_reject(
    const LocationUpdateReject* msg,
    itti_sgsap_location_update_rej_t* itti_msg) {
  auto imsi             = msg->imsi();
  itti_msg->imsi_length = imsi.length();
  strcpy(itti_msg->imsi, imsi.c_str());

  auto cause      = msg->reject_cause();
  itti_msg->cause = static_cast<SgsRejectCause_t>(cause[0]);

  auto lai = msg->location_area_identifier();
  if (lai.length() != 0) {
    itti_msg->presencemask = SGSAP_LAI;
    lai_t itti_lai;
    bytes_to_lai(lai.c_str(), &itti_lai);
    itti_msg->laicsfb = itti_lai;
  } else {
    itti_msg->presencemask = 0;
  }

  return;
}

void convert_proto_msg_to_itti_sgsap_paging_request(
    const PagingRequest* msg, itti_sgsap_paging_request_t* itti_msg) {
  auto imsi             = msg->imsi();
  itti_msg->imsi_length = imsi.length();
  strcpy(itti_msg->imsi, imsi.c_str());

  auto service_indicator      = msg->service_indicator();
  itti_msg->service_indicator = static_cast<uint8_t>(service_indicator[0]);

  uint8_t presencemask = 0;
  // optional fields
  auto tmsi = msg->tmsi();
  if (tmsi.length() != 0) {
    presencemask = presencemask | PAGING_REQUEST_TMSI_PARAMETER_PRESENT;
    tmsi_t itti_tmsi;
    bytes_to_tmsi(tmsi.c_str(), &itti_tmsi);
    itti_msg->opt_tmsi = itti_tmsi;
  }

  auto cli = msg->cli();
  if (cli.length() != 0) {
    presencemask = presencemask | PAGING_REQUEST_CLI_PARAMETER_PRESENT;

    unsigned char* cli_char = new unsigned char[cli.length()];
    for (int i = 0; i < cli.length(); ++i) {
      cli_char[i] = cli[i];
    }

    tagbstring* cli_tbstr = new tagbstring();
    cli_tbstr->mlen       = cli.length();
    cli_tbstr->slen       = cli.length();
    cli_tbstr->data       = cli_char;

    itti_msg->opt_cli = cli_tbstr;
  }

  auto lai = msg->location_area_identifier();
  if (lai.length() != 0) {
    presencemask = presencemask | PAGING_REQUEST_LAI_PARAMETER_PRESENT;
    lai_t itti_lai;
    bytes_to_lai(lai.c_str(), &itti_lai);
    itti_msg->opt_lai = itti_lai;
  }

  itti_msg->presencemask = presencemask;

  return;
}

void convert_proto_msg_to_itti_sgsap_status_t(
    const Status* msg, itti_sgsap_status_t* itti_msg) {
  uint8_t presencemask = 0;
  auto imsi            = msg->imsi();
  if (imsi.length() != 0) {
    presencemask          = presencemask | SGSAP_IMSI;
    itti_msg->imsi_length = imsi.length();
    strcpy(itti_msg->imsi, imsi.c_str());
  }

  auto sgs_cause  = msg->sgs_cause();
  itti_msg->cause = static_cast<SgsCause_t>(sgs_cause[0]);

  itti_msg->presencemask = presencemask;

  return;
}

void convert_proto_msg_to_itti_sgsap_vlr_reset_indication(
    const ResetIndication* msg, itti_sgsap_vlr_reset_indication_t* itti_msg) {
  auto vlr_name    = msg->vlr_name();
  itti_msg->length = vlr_name.length();
  strcpy(itti_msg->vlr_name, vlr_name.c_str());
  return;
}

void convert_proto_msg_to_itti_sgsap_downlink_unitdata(
    const DownlinkUnitdata* msg, itti_sgsap_downlink_unitdata_t* itti_msg) {
  auto imsi             = msg->imsi();
  itti_msg->imsi_length = imsi.length();
  strcpy(itti_msg->imsi, imsi.c_str());

  auto nas_msg = msg->nas_message_container();
  if (nas_msg.length() > 0) {
    itti_msg->nas_msg_container =
        bfromcstr_with_str_len(nas_msg.c_str(), nas_msg.length());
  }

  return;
}

void convert_proto_msg_to_itti_sgsap_release_req(
    const ReleaseRequest* msg, itti_sgsap_release_req_t* itti_msg) {
  uint8_t presencemask  = 0;
  auto imsi             = msg->imsi();
  itti_msg->imsi_length = imsi.length();
  strcpy(itti_msg->imsi, imsi.c_str());

  auto sgs_cause = msg->sgs_cause();
  if (sgs_cause.length() != 0) {
    presencemask        = presencemask | RELEASE_REQ_CAUSE_PARAMETER_PRESENT;
    itti_msg->opt_cause = static_cast<SgsCause_t>(sgs_cause[0]);
  }
  itti_msg->presencemask = presencemask;
  return;
}

void convert_proto_msg_to_itti_sgsap_alert_request(
    const AlertRequest* msg, itti_sgsap_alert_request_t* itti_msg) {
  auto imsi             = msg->imsi();
  itti_msg->imsi_length = imsi.length();
  strcpy(itti_msg->imsi, imsi.c_str());

  return;
}

void convert_proto_msg_to_itti_sgsap_service_abort_req(
    const ServiceAbortRequest* msg, itti_sgsap_service_abort_req_t* itti_msg) {
  auto imsi             = msg->imsi();
  itti_msg->imsi_length = imsi.length();
  strcpy(itti_msg->imsi, imsi.c_str());
  return;
}

void convert_proto_msg_to_itti_sgsap_mm_information_req(
    const MMInformationRequest* msg,
    itti_sgsap_mm_information_req_t* itti_msg) {
  auto imsi             = msg->imsi();
  itti_msg->imsi_length = imsi.length();
  strcpy(itti_msg->imsi, imsi.c_str());

  return;
}

}  // namespace magma
