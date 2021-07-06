"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

This file contains modifications of the core spyne functionality. This is done
using child classes and function override to avoid modifying spyne code itself.
Each function below is a modified version of the parent function. These
modifications are required because:
1) Spyne is not fully python3-compliant
2) Not all parts of the TR-069 spec are possible through spyne APIs (e.g RPC
   calls from server to client in HTTP responses)
3) Minor enhancements for debug-ability
"""

from lxml import etree
from magma.enodebd.logger import EnodebdLogger as logger
from spyne.application import Application
from spyne.interface._base import Interface
from spyne.protocol.soap import Soap11
from spyne.protocol.xml import XmlDocument


class Tr069Interface(Interface):
    """ Modified base interface class. """

    def reset_interface(self):
        super(Tr069Interface, self).reset_interface()
        # Replace default namespace prefix (may not strictly be
        # required, but makes it easier to debug)
        del self.nsmap['tns']
        self.nsmap['cwmp'] = self.get_tns()
        self.prefmap[self.get_tns()] = 'cwmp'
        # To validate against the xsd:<types>, the namespace
        # prefix is expected to be the same
        del self.nsmap['xs']
        self.nsmap['xsd'] = 'http://www.w3.org/2001/XMLSchema'
        self.prefmap['http://www.w3.org/2001/XMLSchema'] = 'xsd'


class Tr069Application(Application):
    """ Modified spyne application. """

    def __init__(
        self, services, tns, name=None, in_protocol=None,
        out_protocol=None, config=None,
    ):
        super(Tr069Application, self).__init__(
            services, tns, name, in_protocol, out_protocol, config,
        )
        # Use modified interface class
        self.interface = Tr069Interface(self)


class Tr069Soap11(Soap11):
    """ Modified SOAP protocol. """

    def __init__(self, *args, **kwargs):
        super(Tr069Soap11, self).__init__(*args, **kwargs)
        # Disabling type resolution as a workaround for
        # https://github.com/arskom/spyne/issues/567
        self.parse_xsi_type = False
        # Bug in spyne is cleaning up the default XSD namespace
        # and causes validation issues on TR-069 clients
        self.cleanup_namespaces = False

    def create_in_document(self, ctx, charset=None):
        """
        In TR-069, the ACS (e.g Magma) is an HTTP server, but acts as a client
        for SOAP messages. This is done by the CPE (e.g ENodeB) sending an
        empty HTTP request, and the ACS responding with a SOAP request in the
        HTTP response. This code replaces an empty HTTP request with a string
        that gets decoded to a call to the 'EmptyHttp' RPC .
        """

        # Try cp437 as default to ensure that we dont get any decoding errors,
        #  since it uses 1-byte encoding and has a 'full' char map
        if not charset:
            charset = 'cp437'

        # Convert from generator to bytes before doing comparison
        # Re-encode to chosen charset to remove invalid characters
        in_string = b''.join(ctx.in_string).decode(charset, 'ignore')
        ctx.in_string = [in_string.encode(charset, 'ignore')]
        if ctx.in_string == [b'']:
            ctx.in_string = [
                b'<soap11env:Envelope xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap11env="http://schemas.xmlsoap.org/soap/envelope/">/n'
                b'   <soap11env:Body>/n'
                b'       <cwmp:EmptyHttp/>/n'
                b'   </soap11env:Body>/n'
                b'</soap11env:Envelope>',
            ]

        super(Tr069Soap11, self).create_in_document(ctx, charset)

    def decompose_incoming_envelope(self, ctx, message=XmlDocument.REQUEST):
        """
        For TR-069, the SOAP fault message (CPE->ACS) contains useful
        information, and should not result in another fault response (ACS->CPE).
        Strip the outer SOAP fault structure, so that the CWMP fault structure
        is treated as a normal RPC call (to the 'Fault' function).
        """
        super(Tr069Soap11, self).decompose_incoming_envelope(ctx, message)

        if ctx.in_body_doc.tag == '{%s}Fault' % self.ns_soap_env:
            faultstring = ctx.in_body_doc.findtext('faultstring')
            if not faultstring or 'CWMP fault' not in faultstring:
                # Not a CWMP fault
                return

            # Strip SOAP fault structure, leaving inner CWMP fault structure
            detail_elem = ctx.in_body_doc.find('detail')
            if detail_elem is not None:
                detail_children = list(detail_elem)
                if len(detail_children):
                    if len(detail_children) > 1:
                        logger.warning(
                            "Multiple detail elements found in SOAP"
                            " fault - using first one",
                        )
                    ctx.in_body_doc = detail_children[0]
                    ctx.method_request_string = ctx.in_body_doc.tag
                    self.validate_body(ctx, message)

    def get_call_handles(self, ctx):
        """
        Modified function to fix bug in receiving SOAP fault. In this case,
        ctx.method_request_string is None, so 'startswith' errors out.
        """
        if ctx.method_request_string is None:
            return []

        return super(Tr069Soap11, self).get_call_handles(ctx)

    def serialize(self, ctx, message):
        # Workaround for issue https://github.com/magma/magma/issues/7869
        # Updates to ctx.descriptor.out_message.Attributes.sub_name are taking
        # effect on the descriptor. But when puled from _attrcache dictionary,
        # it still has a stale value.
        # Force repopulation of dictionary by deleting entry
        # TODO Remove this code once we have a better fix
        if (ctx.descriptor.out_message in self._attrcache):
            del self._attrcache[ctx.descriptor.out_message]  # noqa: WPS529

        super(Tr069Soap11, self).serialize(ctx, message)

        # Keep XSD namespace
        etree.cleanup_namespaces(ctx.out_document, keep_ns_prefixes=['xsd'])
