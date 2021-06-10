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
#pragma once

#include <grpc++/grpc++.h>
#include <grpcpp/impl/codegen/status.h>

#include "feg/protos/csfb.grpc.pb.h"

namespace grpc {
class ServerContext;
}  // namespace grpc
namespace magma {
namespace feg {
class AlertRequest;
class DownlinkUnitdata;
class EPSDetachAck;
class IMSIDetachAck;
class LocationUpdateAccept;
class LocationUpdateReject;
class MMInformationRequest;
class PagingRequest;
class ReleaseRequest;
class ResetAck;
class ResetIndication;
class ServiceAbortRequest;
class Status;
}  // namespace feg
namespace orc8r {
class Void;
}  // namespace orc8r
}  // namespace magma

using grpc::ServerContext;

namespace magma {
using namespace feg;
using namespace orc8r;

class CSFBGatewayServiceImpl final : public CSFBGatewayService::Service {
 public:
  CSFBGatewayServiceImpl();

  /*
   * Sent from the VLR to the MME to request an indication
       when the next activity of a UE is detected
   *
   * @param context: the grpc Server context
   * @param request: AlertRequest, contains IMSI
   * @param response (out): Void defined in common.proto
   * @return grpc Status instance
   */
  grpc::Status AlertReq(
      ServerContext* context, const AlertRequest* request,
      Void* response) override;

  /*
   * Sent from the VLR to the MME to transparently relay a NAS message
   *
   * @param context: the grpc Server context
   * @param request: DownlinkUnitdata, contains IMSI, NAS message container
   * @param response (out): Void defined in common.proto
   * @return grpc Status instance
   */
  grpc::Status Downlink(
      ServerContext* context, const DownlinkUnitdata* request,
      Void* response) override;
  /*
   * Sent from the VLR to the MME to acknowledge
       a previous SGsAP-EPS-DETACH-INDICATION message
   *
   * @param context: the grpc Server context
   * @param request: EPSDetachAck, contains IMSI
   * @param response (out): Void defined in common.proto
   * @return grpc Status instance
   */
  grpc::Status EPSDetachAc(
      ServerContext* context, const EPSDetachAck* request,
      Void* response) override;

  /*
   * Sent from the VLR to the MME to acknowledge
       a previous SGsAP-IMSI-DETACH-INDICATION message
   *
   * @param context: the grpc Server context
   * @param request: IMSIDetachAck, contains IMSI
   * @param response (out): Void defined in common.proto
   * @return grpc Status instance
   */
  grpc::Status IMSIDetachAc(
      ServerContext* context, const IMSIDetachAck* request,
      Void* response) override;

  /*
   * Sent from the VLR to the MME to indicate that update or IMSI attach
       in the VLR has been completed
   *
   * @param context: the grpc Server context
   * @param request: LocationUpdateAccept, contains IMSI, LAI,
       new TMSI or IMSI(Optional)
   * @param response (out): Void defined in common.proto
   * @return grpc Status instance
   */
  grpc::Status LocationUpdateAcc(
      ServerContext* context, const LocationUpdateAccept* request,
      Void* response) override;

  /*
   * Sent from the VLR to the MME to indicate that
       location update or IMSI attach has failed
   *
   * @param context: the grpc Server context
   * @param request: LocationUpdateReject, contains IMSI, reject cause,
       LAI (Optional)
   * @param response (out): Void defined in common.proto
   * @return grpc Status instance
   */
  grpc::Status LocationUpdateRej(
      ServerContext* context, const LocationUpdateReject* request,
      Void* response) override;

  /*
   * Sent from the VLR to the MME to provide the UE
       with subscriber specific information
   *
   * @param context: the grpc Server context
   * @param request: MMInformationRequest, contains IMSI, MM information
   * @param response (out): Void defined in common.proto
   * @return grpc Status instance
   */
  grpc::Status MMInformationReq(
      ServerContext* context, const MMInformationRequest* request,
      Void* response) override;

  /*
   * Sent from the VLR to the MME, containing sufficient information
       to allow the paging message to be transmitted
       by the correct cells at the correct time
   *
   * @param context: the grpc Server context
   * @param request: PagingRequest, contains IMSI, VLR name, service indicator,
       and other optional fields
   * @param response (out): Void defined in common.proto
   * @return grpc Status instance
   */
  grpc::Status PagingReq(
      ServerContext* context, const PagingRequest* request,
      Void* response) override;

  /*
   * Sent from the VLR to the MME when the VLR determines that
       there are no more NAS messages to be exchanged
       between the VLR and the UE, or when a further exchange of NAS messages
       for the specified UE is not possible due to an error.
   *
   * @param context: the grpc Server context
   * @param request: ReleaseRequest, contains IMSI, SGs cause (optional)
   * @param response (out): Void defined in common.proto
   * @return grpc Status instance
   */
  grpc::Status ReleaseReq(
      ServerContext* context, const ReleaseRequest* request,
      Void* response) override;

  /*
   * Sent from the VLR to the MME to abort a mobile
       terminating CS fallback call during call establishment
   *
   * @param context: the grpc Server context
   * @param request: ServiceAbortRequest, contains IMSI
   * @param response (out): Void defined in common.proto
   * @return grpc Status instance
   */
  grpc::Status ServiceAbort(
      ServerContext* context, const ServiceAbortRequest* request,
      Void* response) override;

  /*
   * Sent from the VLR to the MME to acknowledge
       a previous SGsAP-RESET-INDICATION message. This message indicates that
       all the SGs associations to the VLR or have been marked as invalid.
   *
   * @param context: the grpc Server context
   * @param request: ResetAck, contains VLR name
   * @param response (out): Void defined in common.proto
   * @return grpc Status instance
   */
  grpc::Status VLRResetAck(
      ServerContext* context, const ResetAck* request, Void* response) override;

  /*
   * Sent from the VLR to the MME to indicate that a failure in the VLR has
   * occurred and all the SGs associations to the VLR are be marked as invalid.
   *
   * @param context: the grpc Server context
   * @param request: ResetIndication, contains VLR name
   * @param response (out): Void defined in common.proto
   * @return grpc Status instance
   */
  grpc::Status VLRResetIndication(
      ServerContext* context, const ResetIndication* request,
      Void* response) override;

  /*
   * Sent from the VLR to the MME to indicate an error
   *
   * @param context: the grpc Server context
   * @param request: Status, contains IMSI (optional), SGs cause, error message
   * @param response (out): Void defined in common.proto
   * @return grpc Status instance
   */
  grpc::Status VLRStatus(
      ServerContext* context, const Status* request, Void* response) override;
};

}  // namespace magma
