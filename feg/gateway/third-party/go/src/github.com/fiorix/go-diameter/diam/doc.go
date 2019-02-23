// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

/*
Package diam provides support for the Diameter Base Protocol for Go.
See RFC 6733 for details.

go-diameter is minimalist implementation of the Diameter Base Protocol,
organized in sub-packages with specific functionality:

 * diam: the main package, provides the capability of encoding and
         decoding messages, and a client and server API similar to net/http.

 * diam/diamtest: Server test API analogous to net/http/httptest.

 * diam/avp: Diameter attribute-value-pairs codes and flags.

 * diam/datatype: AVP data types (e.g. Unsigned32, OctetString).

 * diam/dict: a dictionary parser that supports collections of dictionaries.

If you're looking to go right into code, see the examples subdirectory for
applications like clients and servers.


Diameter Applications

All diameter applications require at least the following:

 * A dictionary with the application id, its commands and message formats
 * A program that implements the application, driven by the dictionary

The diam/dict sub-package supports the base application (id 0, RFC 6733)
and the credit control application (id 4, RFC 4006). Each application
has its own commands and messages, and their pre-defined AVPs.

AVP data have specific data types, like UTF8String, Unsigned32 and so on.
Fortunately, those data types map well with Go types, which makes things
easier for us. However, the AVP data types have specific properties like
padding for certain strings, which have to be taken care of. The sub-package
diam/datatype handles it all.

At last, the diam package is used to build clients and servers using
an API very similar to the one of net/http. To initiate the client or
server, you'll have to pass a dictionary. Messages sent and received
are encoded and decoded using the dictionary automatically.

The API of clients and servers require that you assign handlers for
certain messages, similar to how you route HTTP endpoints. In the
handlers, you'll receive messages already decoded.
*/
package diam
