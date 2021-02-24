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

#include "sctpd_downlink_impl.h"

#include <arpa/inet.h>
#include <assert.h>
#include <netinet/sctp.h>
#include <unistd.h>
#include "sctpd.h"
#include "util.h"

namespace magma {
namespace sctpd {

SctpdDownlinkImpl::SctpdDownlinkImpl(SctpEventHandler &uplink_handler):
  _uplink_handler(uplink_handler),
  _sctp_4G_connection(nullptr),
  _sctp_5G_connection(nullptr)
{
}

Status SctpdDownlinkImpl::create_sctp_connection(
  std::unique_ptr<SctpConnection>& sctp_connection,
  const InitReq *req,
InitRes *res)
{
  if (sctp_connection != nullptr && !req->force_restart()) {
	MLOG(MERROR) << "SctpdDownlinkImpl::Init reusing existing connection";
	res->set_result(InitRes::INIT_OK);
	return Status::OK;
  }
  if (sctp_connection != nullptr) {
	MLOG(MINFO)<< "SctpdDownlinkImpl::Init cleaning up sctp_desc and listener";
	auto conn = std::move(sctp_connection);
	conn->Close();
  }
  MLOG(MINFO) << "SctpdDownlinkImpl::Init creating new socket and listener";
  try {
	sctp_connection = std::make_unique<SctpConnection>(*req, _uplink_handler);
  } catch (...) {
	res->set_result(InitRes::INIT_FAIL);
	return Status::OK;
  }
  sctp_connection->Start();
  return Status::OK;
}

Status SctpdDownlinkImpl::Init(
  ServerContext *context,
  const InitReq *req,
  InitRes *res)
{
  MLOG(MINFO) << "SctpdDownlinkImpl::req->ppid()=" << std::to_string(req->ppid());
  MLOG(MINFO) << "SctpdDownlinkImpl::req->port()=" << std::to_string(req->port());

 Status rc;

  if (req->ppid() == S1AP) {
	rc = create_sctp_connection(_sctp_4G_connection, req, res);

}
  else if (req->ppid() == NGAP) {
	rc = create_sctp_connection(_sctp_5G_connection, req, res);

}
	res->set_result(InitRes::INIT_OK);
	return Status::OK;
}

Status SctpdDownlinkImpl::SendDl(
  ServerContext *context,
  const SendDlReq *req,
  SendDlRes *res)
{
  MLOG(MDEBUG) << "SctpdDownlinkImpl::SendDl starting";

  try {
	  if (req->ppid() == S1AP )
		  _sctp_4G_connection->Send(req->assoc_id(), req->stream(), req->payload());
	  else
		  _sctp_5G_connection->Send(req->assoc_id(), req->stream(), req->payload());

  } catch (...) {
    res->set_result(SendDlRes::SEND_DL_FAIL);
    return Status::OK;
  }

  res->set_result(SendDlRes::SEND_DL_OK);
  return Status::OK;
}

void SctpdDownlinkImpl::stop()
{
	if (_sctp_4G_connection != nullptr) {
		 _sctp_4G_connection->Close();
	}
	if (_sctp_5G_connection != nullptr) {
		 _sctp_5G_connection->Close();
	}
}

} // namespace sctpd
} // namespace magma
