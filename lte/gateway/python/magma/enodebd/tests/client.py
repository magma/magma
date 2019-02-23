"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
from spyne import RemoteProcedureBase, RemoteService
from spyne.client.http import HttpClient
# pylint: disable=no-name-in-module, import-error
from spyne.util.six.moves.urllib.error import HTTPError
from spyne.util.six.moves.urllib.request import Request, urlopen

from magma.common.misc_utils import get_ip_from_if
from magma.configuration.service_configs import load_service_config
from magma.enodebd.tr069.models import CWMP_NS, DummyInput, ID
from magma.enodebd.tr069.rpc_methods import AutoConfigServer
from magma.enodebd.tr069.spyne_mods import Tr069Application, Tr069Soap11


class _Tr069RemoteProcedure(RemoteProcedureBase):
    """ Modified remote procedure class. Search for 'Tr069' to see mods. """
    def __call__(self, *args, **kwargs):
        # there's no point in having a client making the same request more than
        # once, so if there's more than just one context, it is a bug.
        # the comma-in-assignment trick is a general way of getting the first
        # and the only variable from an iterable. so if there's more than one
        # element in the iterable, it'll fail miserably.
        # pylint: disable=attribute-defined-outside-init
        self.ctx, = self.contexts

        # sets ctx.out_object
        self.get_out_object(self.ctx, args, kwargs)

        # sets ctx.out_string
        self.get_out_string(self.ctx)

        # Tr069 modifications - fix bug in handling of binary string.
        # May be a python3 issue.
        out_string = b''.join(self.ctx.out_string)
        request = Request(self.url, out_string, {"Content-Type": "text/xml"})
        code = 200
        try:
            response = urlopen(request)
            self.ctx.in_string = [response.read()]

        except HTTPError as e:
            code = e.code
            self.ctx.in_string = [e.read()]

        # this sets ctx.in_error if there's an error, and ctx.in_object if
        # there's none.
        self.get_in_object(self.ctx)

        if not (self.ctx.in_error is None):
            raise self.ctx.in_error
        elif code >= 400:
            raise self.ctx.in_error
        else:
            return self.ctx.in_object


class Tr069HttpClient(HttpClient):
    """ Modified HTTP client class. """
    def __init__(self, url, app):
        super(Tr069HttpClient, self).__init__(url, app)

        # Tr069 modifications - use modified remote procedure class
        self.service = RemoteService(_Tr069RemoteProcedure, url, app)


def main():
    """ This module is used for manual testing of the TR-069 server """
    config = load_service_config("enodebd")

    app = Tr069Application([AutoConfigServer], CWMP_NS,
                           in_protocol=Tr069Soap11(validator="soft"),
                           out_protocol=Tr069Soap11())

    ip_address = get_ip_from_if(config['tr069']['interface'])
    client = Tr069HttpClient(
        "http://%s:%s" % (ip_address, config["tr069"]["port"]),
        app)

    client.set_options(out_header=ID("123", mustUnderstand="1"))
    rpc_methods = client.service.get_rpc_methods()
    for rpc_method in rpc_methods:
        print("Method: %s" % rpc_method)

    inform_req = client.factory.create("Inform")
    inform_req.DeviceId = client.factory.create("DeviceIdStruct")
    inform_req.DeviceId.Manufacturer = "Magma"
    inform_req.DeviceId.OUI = "ABCDEF"
    inform_req.DeviceId.ProductClass = "TopClass"
    inform_req.DeviceId.SerialNumber = "123456789"
    inform_req.Event = None
    inform_req.MaxEnvelopes = 1
    inform_req.CurrentTime = None
    inform_req.RetryCount = 4
    inform_req.ParameterList = None
    client.set_options(out_header=ID("456", mustUnderstand="1"))
    client.service.Inform(inform_req)

    dummy = DummyInput()
    dummy.Field1 = 5
    rsp = client.service.EmptyHttp(dummy)
    print("EmptyHttp response = ", rsp)

    paramNames = client.factory.create("GetParameterNamesResponse")
    paramNames.ParameterList = client.factory.create("ParameterInfoList")
    paramNames.ParameterList.ParameterInfoStruct =\
        [client.factory.create("ParameterInfoStruct")]
    paramNames.ParameterList.ParameterInfoStruct[0].Name = "Parameter1"
    paramNames.ParameterList.ParameterInfoStruct[0].Writable = True
    rsp = client.service.GetParameterNamesResponse(paramNames)
    print("GetParameterNamesResponse response = ", rsp)


if __name__ == "__main__":
    main()
