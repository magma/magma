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

package diameter

import "github.com/fiorix/go-diameter/v4/diam/datatype"

//https://tools.ietf.org/html/rfc4005#page-16
/*  Abort-Session Request (ASR)
< Session-Id >
{ Origin-Host }
{ Origin-Realm }
{ Destination-Realm }
{ Destination-Host }
{ Auth-Application-Id }
[ User-Name ]
[ Origin-State-Id ]
* [ Proxy-Info ]
* [ Route-Record ]
* [ AVP ]
*/
type ASR struct {
	SessionID         string                    `avp:"Session-Id"`
	OriginHost        datatype.DiameterIdentity `avp:"Origin-Host"`
	OriginRealm       datatype.DiameterIdentity `avp:"Origin-Realm"`
	DestinationHost   datatype.DiameterIdentity `avp:"Destination-Host"`
	DestinationRealm  datatype.DiameterIdentity `avp:"Destination-Realm"`
	AuthApplicationId datatype.Unsigned32       `avp:"Auth-Application-Id"`
	UserName          datatype.UTF8String       `avp:"User-Name"`
	OriginStateId     datatype.Unsigned32       `avp:"Origin-State-Id"`
}

/*  Abort-Session Answer (ASA)
< Session-Id >
{ Origin-Host }
{ Origin-Realm }
[ Result-Code ]
[ Experimental-Result ]
[ Origin-State-Id ]
[ Error-Message ]
[ Error-Reporting-Host ]
*[ Failed-AVP ]
*[ Redirected-Host ]
[ Redirected-Host-Usage ]
[ Redirected-Max-Cache-Time ]
*[ Proxy-Info ]
*[ AVP ]
*/
type ASA struct {
	SessionID  string `avp:"Session-Id"`
	ResultCode uint32 `avp:"Result-Code"`
}
