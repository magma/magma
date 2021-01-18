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
#include "magma_logging.h"
namespace magma {
namespace sctpd {

 using namespace std;
SctpdDownlinkImpl::SctpdDownlinkImpl(SctpEventHandler &uplink_handler):
  _uplink_handler(uplink_handler),
  _sctp_4G_connection(nullptr),
  _sctp_5G_connection(nullptr)
{
}

Status SctpdDownlinkImpl::Init(
  ServerContext *context,
  const InitReq *req,
  InitRes *res)
{
  MLOG(MERROR) << "SctpdDownlinkImpl::Init starting ERROR";
  MLOG(MERROR) << "SctpdDownlinkImpl::req->ppid()=" << std::to_string(req->ppid());
  MLOG(MERROR) << "SctpdDownlinkImpl::req->port()=" << std::to_string(req->port());
  
  if (req->ppid() == S1AP) {
	   if (_sctp_4G_connection != nullptr && !req->force_restart()) {
		    MLOG(MERROR) << "SctpdDownlinkImpl::Init reusing existing connection";
		    res->set_result(InitRes::INIT_OK);
		    return Status::OK;
	   }
	   if (_sctp_4G_connection != nullptr) {
		    MLOG(MINFO)<< "SctpdDownlinkImpl::Init cleaning up sctp_desc and listener";
		    auto conn = std::move(_sctp_4G_connection);
		    conn->Close();
	   }
	   MLOG(MERROR) << "SctpdDownlinkImpl::Init creating new socket and listener";
	   try {
	   MLOG(MDEBUG) << "SctpdDownlinkImpl::1";
		    _sctp_4G_connection = std::make_unique<SctpConnection>(*req, _uplink_handler);
		     } catch (...) {
	   MLOG(MERROR) << "SctpdDownlinkImpl::2";
			      res->set_result(InitRes::INIT_FAIL);
			      return Status::OK;
		     }
	   MLOG(MERROR) << "SctpdDownlinkImpl::3";
	    _sctp_4G_connection->Start();
	   MLOG(MERROR) << "SctpdDownlinkImpl::4";
  }else if (req->ppid() == NGAP) {
	  MLOG(MERROR) << "SctpdDownlinkImpl::Init starting for 5G";
	  if (_sctp_5G_connection != nullptr && !req->force_restart()) {
		  MLOG(MERROR) << "SctpdDownlinkImpl::Init reusing existing 5G connection";
		  res->set_result(InitRes::INIT_OK);
		  return Status::OK;
	  }
	  if (_sctp_5G_connection != nullptr) {
		  MLOG(MERROR) << "SctpdDownlinkImpl::Init cleaning up sctp_desc and listener of 5G";
	          auto conn = std::move(_sctp_5G_connection);
		  conn->Close();
		  }
		  MLOG(MERROR) << "SctpdDownlinkImpl::Init creating new socket and listener of 5G";
		  try {
		  	_sctp_5G_connection = std::make_unique<SctpConnection>(*req, _uplink_handler);
			} catch (...) {
			res->set_result(InitRes::INIT_FAIL);
			return Status::OK;
			}
			_sctp_5G_connection->Start();
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
  MLOG(MERROR) << " assoc_id = " << std::to_string(req->assoc_id());
  MLOG(MERROR) << " stream = " << std::to_string(req->stream());

  try {
	  if (req->ppid() == S1AP ){
  		   MLOG(MDEBUG) << "SctpdDownlinkImpl::S1AP";
		  _sctp_4G_connection->Send(req->assoc_id(), req->stream(), req->payload());
		}
	  else{
  		   MLOG(MDEBUG) << "SctpdDownlinkImpl::NGAP";
		  _sctp_5G_connection->Send(req->assoc_id(), req->stream(), req->payload());
	      }
  } catch (...) {
  		   MLOG(MDEBUG) << "SctpdDownlinkImpl::FAIL";
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
