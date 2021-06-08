"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""


from spyne.model import ComplexModel
from spyne.model.complex import XmlAttribute, XmlData
from spyne.model.primitive import (
    Boolean,
    DateTime,
    Integer,
    String,
    UnsignedInteger,
)
from spyne.util.odict import odict

# Namespaces
XSI_NS = 'http://www.w3.org/2001/XMLSchema-instance'
SOAP_ENV = 'http://schemas.xmlsoap.org/soap/envelope/'
SOAP_ENC = 'http://schemas.xmlsoap.org/soap/encoding/'
CWMP_NS = 'urn:dslforum-org:cwmp-1-0'


class Tr069ComplexModel(ComplexModel):
    """ Base class for TR-069 models, to set common attributes. Does not appear
        in CWMP XSD file. """
    __namespace__ = CWMP_NS

    def as_dict(self):
        """
        Overriding default implementation to fix memory leak. Can remove if
        or after https://github.com/arskom/spyne/pull/579 lands.
        """
        flat_type_info = self.get_flat_type_info(self.__class__)
        return dict((
            (k, getattr(self, k)) for k in flat_type_info
            if getattr(self, k) is not None
        ))


class anySimpleType(Tr069ComplexModel):
    """ Type used to transfer simple data of various types. Data type is
        defined in 'type' XML attribute. Data is handled as a string. """
    _type_info = odict()
    _type_info["type"] = XmlAttribute(String, ns=XSI_NS)
    _type_info["Data"] = XmlData(String)

    def __repr__(self):
        """For types we can't resolve only print the datum"""
        return self.Data


# SOAP Header Elements


class ID(Tr069ComplexModel):
    # Note: for some reason, XmlAttribute/XmlData pairs MUST be ordered, with
    # XmlAttribute coming first. This appears to be a spyne bug (something to do
    # with spyne.interface._base.add_class())
    _type_info = odict()
    _type_info["mustUnderstand"] = XmlAttribute(String, ns=SOAP_ENV)
    _type_info["Data"] = XmlData(String)


class HoldRequests(Tr069ComplexModel):
    _type_info = odict()
    _type_info["mustUnderstand"] = XmlAttribute(String, ns=SOAP_ENV)
    _type_info["Data"] = XmlData(Boolean)


# SOAP Fault Extensions


class SetParameterValuesFault(Tr069ComplexModel):
    _type_info = odict()
    _type_info["ParameterName"] = String
    _type_info["FaultCode"] = UnsignedInteger
    _type_info["FaultString"] = String


class Fault(Tr069ComplexModel):
    _type_info = odict()
    _type_info["FaultCode"] = UnsignedInteger
    _type_info["FaultString"] = String
    _type_info["SetParameterValuesFault"] = SetParameterValuesFault.customize(
        max_occurs='unbounded',
    )


# Type definitions used in messages


class MethodList(Tr069ComplexModel):
    _type_info = odict()
    _type_info["string"] = String(max_length=64, max_occurs='unbounded')
    _type_info["arrayType"] = XmlAttribute(String, ns=SOAP_ENC)


class FaultStruct(Tr069ComplexModel):
    _type_info = odict()
    _type_info["FaultCode"] = Integer
    _type_info["FaultString"] = String(max_length=256)


class DeviceIdStruct(Tr069ComplexModel):
    _type_info = odict()
    _type_info["Manufacturer"] = String(max_length=64)
    _type_info["OUI"] = String(length=6)
    _type_info["ProductClass"] = String(max_length=64)
    _type_info["SerialNumber"] = String(max_length=64)


class EventStruct(Tr069ComplexModel):
    _type_info = odict()
    _type_info["EventCode"] = String(max_length=64)
    _type_info["CommandKey"] = String(max_length=32)


class EventList(Tr069ComplexModel):
    _type_info = odict()
    _type_info["EventStruct"] = EventStruct.customize(max_occurs='unbounded')
    _type_info["arrayType"] = XmlAttribute(String, ns=SOAP_ENC)


class ParameterValueStruct(Tr069ComplexModel):
    _type_info = odict()
    _type_info["Name"] = String
    _type_info["Value"] = anySimpleType


class ParameterValueList(Tr069ComplexModel):
    _type_info = odict()
    _type_info["ParameterValueStruct"] = ParameterValueStruct.customize(
        max_occurs='unbounded',
    )
    _type_info["arrayType"] = XmlAttribute(String, ns=SOAP_ENC)


class ParameterInfoStruct(Tr069ComplexModel):
    _type_info = odict()
    _type_info["Name"] = String(max_length=256)
    _type_info["Writable"] = Boolean


class ParameterInfoList(Tr069ComplexModel):
    _type_info = odict()
    _type_info["ParameterInfoStruct"] = ParameterInfoStruct.customize(max_occurs='unbounded')
    _type_info["arrayType"] = XmlAttribute(String, ns=SOAP_ENC)


class ParameterNames(Tr069ComplexModel):
    _type_info = odict()
    _type_info["string"] = String.customize(max_occurs='unbounded', max_length=256)
    _type_info["arrayType"] = XmlAttribute(String, ns=SOAP_ENC)


class ParameterKeyType(String.customize(max_length=32)):
    pass


class AccessList(Tr069ComplexModel):
    _type_info = odict()
    _type_info["string"] = String.customize(max_occurs='unbounded', max_length=64)
    _type_info["arrayType"] = XmlAttribute(String, ns=SOAP_ENC)


class SetParameterAttributesStruct(Tr069ComplexModel):
    _type_info = odict()
    _type_info["Name"] = String(max_length=256)
    _type_info["NotificationChange"] = Boolean
    _type_info["Notification"] = Integer
    _type_info["AccessListChange"] = Boolean
    _type_info["AccessList"] = AccessList


class SetParameterAttributesList(Tr069ComplexModel):
    _type_info = odict()
    _type_info["SetParameterAttributesStruct"] = SetParameterAttributesStruct.customize(
        max_occurs='unbounded',
    )
    _type_info["arrayType"] = XmlAttribute(String, ns=SOAP_ENC)


class ParameterAttributeStruct(Tr069ComplexModel):
    _type_info = odict()
    _type_info["Name"] = String(max_length=256)
    _type_info["Notification"] = Integer
    _type_info["AccessList"] = AccessList


class ParameterAttributeList(Tr069ComplexModel):
    _type_info = odict()
    _type_info["ParameterValueStruct"] = ParameterAttributeStruct.customize(
        max_occurs='unbounded',
    )
    _type_info["arrayType"] = XmlAttribute(String, ns=SOAP_ENC)


class CommandKeyType(String.customize(max_length=32)):
    pass


class ObjectNameType(String.customize(max_length=256)):
    pass


# CPE messages


class SetParameterValues(Tr069ComplexModel):
    _type_info = odict()
    _type_info["ParameterList"] = ParameterValueList
    _type_info["ParameterKey"] = ParameterKeyType


class SetParameterValuesResponse(Tr069ComplexModel):
    _type_info = odict()
    _type_info["Status"] = Integer


class GetParameterValues(Tr069ComplexModel):
    _type_info = odict()
    _type_info["ParameterNames"] = ParameterNames


class GetParameterValuesResponse(Tr069ComplexModel):
    _type_info = odict()
    _type_info["ParameterList"] = ParameterValueList


class GetParameterNames(Tr069ComplexModel):
    _type_info = odict()
    _type_info["ParameterPath"] = String.customize(max_length=256)
    _type_info["NextLevel"] = Boolean


class GetParameterNamesResponse(Tr069ComplexModel):
    _type_info = odict()
    _type_info["ParameterList"] = ParameterInfoList


class SetParameterAttributes(Tr069ComplexModel):
    _type_info = odict()
    _type_info["ParameterList"] = SetParameterAttributesList


class SetParameterAttributesResponse(Tr069ComplexModel):
    # Dummy field required because spyne does not allow 'bare' RPC function with
    # no input parameters. This field is never sent by CPE.
    _type_info = odict()
    _type_info["DummyField"] = UnsignedInteger


class GetParameterAttributes(Tr069ComplexModel):
    _type_info = odict()
    _type_info["ParameterNames"] = ParameterNames


class GetParameterAttributesResponse(Tr069ComplexModel):
    _type_info = odict()
    _type_info["ParameterList"] = ParameterAttributeList


class AddObject(Tr069ComplexModel):
    _type_info = odict()
    _type_info["ObjectName"] = ObjectNameType
    _type_info["ParameterKey"] = ParameterKeyType


class AddObjectResponse(Tr069ComplexModel):
    _type_info = odict()
    _type_info["InstanceNumber"] = UnsignedInteger
    _type_info["Status"] = Integer


class DeleteObject(Tr069ComplexModel):
    _type_info = odict()
    _type_info["ObjectName"] = ObjectNameType
    _type_info["ParameterKey"] = ParameterKeyType


class DeleteObjectResponse(Tr069ComplexModel):
    _type_info = odict()
    _type_info["Status"] = Integer


class Download(Tr069ComplexModel):
    _type_info = odict()
    _type_info["CommandKey"] = CommandKeyType
    _type_info["FileType"] = String(max_length=64)
    _type_info["URL"] = String(max_length=256)
    _type_info["Username"] = String(max_length=256)
    _type_info["Password"] = String(max_length=256)
    _type_info["FileSize"] = UnsignedInteger
    _type_info["TargetFileName"] = String(max_length=256)
    _type_info["DelaySeconds"] = UnsignedInteger
    _type_info["SuccessURL"] = String(max_length=256)
    _type_info["FailureURL"] = String(max_length=256)


class DownloadResponse(Tr069ComplexModel):
    _type_info = odict()
    _type_info["Status"] = Integer
    _type_info["StartTime"] = DateTime
    _type_info["CompleteTime"] = DateTime


class Reboot(Tr069ComplexModel):
    _type_info = odict()
    _type_info["CommandKey"] = CommandKeyType


class RebootResponse(Tr069ComplexModel):
    # Dummy field required because spyne does not allow 'bare' RPC function with
    # no input parameters. This field is never sent by CPE.
    _type_info = odict()
    _type_info["DummyField"] = UnsignedInteger


# ACS messages


class Inform(Tr069ComplexModel):
    _type_info = odict()
    _type_info["DeviceId"] = DeviceIdStruct
    _type_info["Event"] = EventList
    _type_info["MaxEnvelopes"] = UnsignedInteger
    _type_info["CurrentTime"] = DateTime
    _type_info["RetryCount"] = UnsignedInteger
    _type_info["ParameterList"] = ParameterValueList


class InformResponse(Tr069ComplexModel):
    _type_info = odict()
    _type_info["MaxEnvelopes"] = UnsignedInteger


class TransferComplete(Tr069ComplexModel):
    _type_info = odict()
    _type_info["CommandKey"] = CommandKeyType
    _type_info["FaultStruct"] = FaultStruct
    _type_info["StartTime"] = DateTime
    _type_info["CompleteTime"] = DateTime


class TransferCompleteResponse(Tr069ComplexModel):
    # Dummy field required because spyne does not allow 'bare' RPC function with
    # no input parameters. This field is never sent by ACS.
    _type_info = odict()
    _type_info["DummyField"] = UnsignedInteger


class GetRPCMethods(Tr069ComplexModel):
    _type_info = odict()
    _type_info["DummyField"] = UnsignedInteger


class GetRPCMethodsResponse(Tr069ComplexModel):
    _type_info = odict()
    _type_info["MethodList"] = MethodList


#
# Miscellaneous
#

class ParameterListUnion(Tr069ComplexModel):
    """ Union of structures that get instantiated as 'ParameterList' in ACS->CPE
        messages. This is required because AcsToCpeRequests can only have one
        parameter named 'ParameterList', so that must also be a union """
    _type_info = odict()

    # Fields from ParameterValueList
    _type_info["ParameterValueStruct"] = ParameterValueStruct.customize(
        max_occurs='unbounded',
    )
    _type_info["arrayType"] = XmlAttribute(String, ns=SOAP_ENC)

    # Fields from SetParameterAttributesList
    _type_info["SetParameterAttributesStruct"] = \
        SetParameterAttributesStruct.customize(max_occurs='unbounded')
    # arrayType = XmlAttribute(String, ns=SOAP_ENC) - Already covered above


class AcsToCpeRequests(Tr069ComplexModel):
    """ Union of all ACS->CPE requests. Only fields for one request is populated
        per message instance """
    _type_info = odict()

    # Fields for SetParameterValues
    _type_info["ParameterList"] = ParameterListUnion  # See ParameterListUnion for explanation
    _type_info["ParameterKey"] = ParameterKeyType

    # Fields for GetParameterValues
    # _type_info["ParameterList"] = ParameterValueList - Already covered above

    # Fields for GetParameterNames
    _type_info["ParameterPath"] = String.customize(max_length=256)
    _type_info["NextLevel"] = Boolean

    # Fields for SetParameterAttributes
    # _type_info["ParameterList"] = SetParameterAttributesList - Already covered above

    # Fields for GetParameterAttributes
    _type_info["ParameterNames"] = ParameterNames

    # Fields for AddObject
    _type_info["ObjectName"] = ObjectNameType
    _type_info["ParameterKey"] = ParameterKeyType

    # Fields for DeleteObject
    # _type_info["ObjectName"] = ObjectNameType - Already covered above
    # _type_info["ParameterKey"] = ParameterKeyType - Already covered above

    # Fields for Download
    _type_info["CommandKey"] = CommandKeyType
    _type_info["FileType"] = String(max_length=64)
    _type_info["URL"] = String(max_length=256)
    _type_info["Username"] = String(max_length=256)
    _type_info["Password"] = String(max_length=256)
    _type_info["FileSize"] = UnsignedInteger
    _type_info["TargetFileName"] = String(max_length=256)
    _type_info["DelaySeconds"] = UnsignedInteger
    _type_info["SuccessURL"] = String(max_length=256)
    _type_info["FailureURL"] = String(max_length=256)

    # Fields for Reboot
    # _type_info["CommandKey"] = CommandKeyType - Already covered above


class DummyInput(Tr069ComplexModel):
    """ Dummy complex model. Used for 'EmptyHttp' function, because spyne Does
        not handle 'bare' function with no inputs """
    _type_info = odict()
    _type_info["DummyField"] = UnsignedInteger
