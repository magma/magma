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
#pragma once
namespace magma5g {
  #define M5G_SESSION_MANAGEMENT_MESSAGES     0x2e
  #define M5G_MOBILITY_MANAGEMENT_MESSAGES    0x7e
	// 5G Mobility Management Message Types 
  #define REGISTRATION_REQUEST                 0x41
  #define REGISTRATION_ACCEPT                  0x42
  #define REGISTRATION_COMPLETE                0x43
  #define REGISTRATION_REJECT                  0x44
  #define DEREGISTRATION_REQUEST_UE_ORIGIN     0x45
  #define DEREGISTRATION_ACCEPT_UE_ORIGIN      0x46
  #define DEREGISTRATION_REQUEST_UE_TERM       0x47
  #define DEREGISTRATION_ACCEPT_UE_TERM        0x48
  #define SERVICE_REQUEST                      0x4c
  #define SERVICE_REJECT                       0x4d
  #define SERVICE_ACCEPT                       0x4e
  #define CONFIGURATION_UPDATE_COMMAND         0x54
  #define CONFIGURATION_UPDATE_COMPLETE        0x55
  #define AUTHENTICATION_REQUEST               0x56
  #define AUTHENTICATION_RESPONSE              0x57
  #define AUTHENTICATION_REJECT                0x58
  #define AUTHENTICATION_FAILURE               0x59
  #define AUTHENTICATION_RESULT                0x5a
  #define IDENTITY_REQUEST                     0x5b
  #define IDENTITY_RESPONSE                    0x5c
  #define SECURITY_MODE_COMMAND                0x5d
  #define SECURITY_MODE_COMPLETE               0x5e
  #define SECURITY_MODE_REJECT                 0x5f
  #define NOTIFICATION                         0x65
  #define NOTIFICATION_RESPONSE                0x66
  #define ULNASTRANSPORT                       0x67
  #define DLNASTRANSPORT                       0x68
	// IEI for Mobility Management Message
  #define M5GSMOBILEIDENTITY                   0X77
  #define PLMNLIST                             0x4a
  #define TAILIST                              0x54
  #define ALLOWEDNSSAI                         0x15
  #define REJECTEDNSSAI                        0x11
  #define M5GSNETWORKFEATURESUPPORT            0X21
  #define PDUSESSIONSTATUS                     0x50
  #define PDUSESSIONREACTIVATIONRESULT         0x26
  #define PDUSESSIONREACTIVATIONRESULTERROR    0X72
  #define LADNINFORMATION                      0x79
  #define MICOINDICATION                       0xB0
  #define NETWORKSLICINGINDICATION             0x90
  #define SERVICEAREALIST                      0X27
  #define GPRSTIMER3                           0x5E
  #define GPRSTIMER2                           0x5D
  #define EMERGENCYNUMBERLIST                  0x34
  #define EXTENDEDNUMBEREMERGENCYLIST          0x7A
  #define SORTRANSPARANTCONTAINER              0x73
  #define EAPMESSAGE                           0x78
  #define NSSAIINCLUSIONMODE                   0xA0
  #define OPERATORDEFINEDACCESSCATEGORYDEF     0x76
  #define M5GSDRXPARAMETERS                    0x51
  #define NON3GPPNWPROVIDEDPOLICIES            0xD0
  #define EPSBEARERCONTEXTSTATUS               0x60
  #define NASKEYIDENTIFIER                     0xC0
  #define M5GMMCAPABILITY                      0x10
  #define UESECURITYCAPABILITY                 0x2E
  #define S1UENETWORKCAPABILITY                0X17
  #define REQUESTEDNSSAI                       0x2F
  #define UPLINKDATASTATUS                     0x40
  #define UESTATUS                             0x2B
  #define ALLOWEDPDUSESSIONSTATUS              0x25
  #define UEUSSAGESETTING                      0x18
  #define EPSNASMESSAGECONTAINER               0x70
  #define PAYLOADCONTAINERTYPE                 0X80
  #define PAYLOADCONTAINER                     0x7B
  #define M5GSUPDATETYPE                       0x53
  #define NASMESSAGECONTAINER                  0x71
  #define TAI                                  0x52
} // namespace magma5g
