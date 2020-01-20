/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package diameter

import "github.com/fiorix/go-diameter/v4/diam/datatype"

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
