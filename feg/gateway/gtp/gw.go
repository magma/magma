/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gtp

import (
	"github.com/wmnsk/go-gtp/gtpv2"
	"github.com/wmnsk/go-gtp/gtpv2/message"
)

func (c *Client) AddHandlers(handlers map[uint8]gtpv2.HandlerFunc) {
	c.Conn.AddHandlers(handlers)
}

func (c *Client) AddSgwS11Handlers(
	handleCreateSessionRequest,
	handleModifyBearerRequest,
	handleDeleteSessionRequest,
	handleDeleteBearerResponse gtpv2.HandlerFunc) *Client {
	c.Conn.AddHandlers(map[uint8]gtpv2.HandlerFunc{
		message.MsgTypeCreateSessionRequest: handleCreateSessionRequest,
		message.MsgTypeModifyBearerRequest:  handleModifyBearerRequest,
		message.MsgTypeDeleteSessionRequest: handleDeleteSessionRequest,
		message.MsgTypeDeleteBearerResponse: handleDeleteBearerResponse,
	})
	return c
}

func (c *Client) AddSgwS5S8Handlers(
	handleCreateSessionResponse,
	handleModifyBearerRequest,
	handleDeleteSessionResponse,
	handleDeleteBearerRequest gtpv2.HandlerFunc) *Client {
	c.Conn.AddHandlers(map[uint8]gtpv2.HandlerFunc{
		message.MsgTypeCreateSessionResponse: handleCreateSessionResponse,
		message.MsgTypeModifyBearerRequest:   handleModifyBearerRequest,
		message.MsgTypeDeleteSessionResponse: handleDeleteSessionResponse,
		message.MsgTypeDeleteBearerRequest:   handleDeleteBearerRequest,
	})
	return c
}

func (c *Client) AddMmeS11Handlers(
	handleCreateSessionResponse,
	handleModifyBearerResponse,
	handleDeleteSessionResponse gtpv2.HandlerFunc) *Client {

	c.Conn.AddHandlers(map[uint8]gtpv2.HandlerFunc{
		message.MsgTypeCreateSessionResponse: handleCreateSessionResponse,
		message.MsgTypeModifyBearerResponse:  handleModifyBearerResponse,
		message.MsgTypeDeleteSessionResponse: handleDeleteSessionResponse,
	})
	return c
}

func (c *Client) AddPgwS5S8Handlers(
	handleCreateSessionRequest,
	handleDeleteSessionRequest gtpv2.HandlerFunc) *Client {
	c.Conn.AddHandlers(map[uint8]gtpv2.HandlerFunc{
		message.MsgTypeCreateSessionRequest: handleCreateSessionRequest,
		message.MsgTypeDeleteSessionRequest: handleDeleteSessionRequest,
	})
	return c
}
