# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set(libnas_src_OBJS
    src/AmfMessage.cpp
    src/M5GRegistrationRequest.cpp
    src/M5GRegistrationAccept.cpp
    src/M5GRegistrationComplete.cpp
    src/M5GRegistrationReject.cpp
    src/M5GDeRegistrationRequestUEInit.cpp
    src/M5GDeRegistrationAcceptUEInit.cpp
    src/M5GAuthenticationRequest.cpp
    src/M5GAuthenticationResponse.cpp
    src/M5GAuthenticationResult.cpp
    src/M5GAuthenticationFailure.cpp
    src/M5GAuthenticationReject.cpp
    src/M5GSecurityModeCommand.cpp
    src/M5GSecurityModeComplete.cpp
    src/M5GSecurityModeReject.cpp
    src/M5GIdentityRequest.cpp
    src/M5GIdentityResponse.cpp
    src/M5GPDUSessionEstablishmentRequest.cpp
    src/M5GPDUSessionEstablishmentAccept.cpp
    src/M5GPDUSessionEstablishmentReject.cpp
    src/M5GPDUSessionModificationRequest.cpp
    src/M5GPDUSessionModificationReject.cpp
    src/M5GPDUSessionReleaseRequest.cpp
    src/M5GPDUSessionReleaseReject.cpp
    src/SmfMessage.cpp
    src/M5GULNASTransport.cpp
    src/M5GDLNASTransport.cpp
    )

set(libnas_ies_OBJS
    src/ies/M5GExtendedProtocolDiscriminator.cpp
    src/ies/M5GNasKeySetIdentifier.cpp
    src/ies/M5GSMobileIdentity.cpp
    src/ies/M5GMessageType.cpp
    src/ies/M5GSecurityHeaderType.cpp
    src/ies/M5GSRegistrationType.cpp
    src/ies/M5GSRegistrationResult.cpp
    src/ies/M5GSDeRegistrationType.cpp
    src/ies/M5GSpareHalfOctet.cpp
    src/ies/M5GABBA.cpp
    src/ies/M5GAuthenticationParameterAUTN.cpp
    src/ies/M5GAuthenticationParameterRAND.cpp
    src/ies/M5GAuthenticationResponseParameter.cpp
    src/ies/M5GEAPMessage.cpp
    src/ies/M5GNASSecurityAlgorithms.cpp
    src/ies/M5GUESecurityCapability.cpp
    src/ies/M5GSIdentityType.cpp
    src/ies/M5GPayloadContainer.cpp
    src/ies/M5GPayloadContainerType.cpp
    src/ies/M5GPTI.cpp
    src/ies/M5GSSCMode.cpp
    src/ies/M5GIntegrityProtMaxDataRate.cpp
    src/ies/M5GPDUSessionType.cpp
    src/ies/M5GPDUSessionIdentity.cpp
    src/ies/M5GPDUAddress.cpp
    src/ies/M5GQOSRules.cpp
    src/ies/M5GSessionAMBR.cpp
    src/ies/M5GSMCause.cpp
    src/ies/M5GMMCause.cpp
    src/ies/M5GIMEISVRequest.cpp
    src/ies/M5GDNN.cpp
    )

add_library(LIB_NAS5G
    ${libnas_src_OBJS}
    ${libnas_ies_OBJS}
    )

target_link_libraries(LIB_NAS5G
    COMMON MAGMA_LOGGING
    LIB_BSTR LIB_SECU LIB_3GPP
    )

target_include_directories(LIB_NAS5G PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CMAKE_CURRENT_SOURCE_DIR}/src
    ${CMAKE_CURRENT_SOURCE_DIR}/src/ies
    ${CMAKE_CURRENT_SOURCE_DIR}/include
    ${CMAKE_CURRENT_SOURCE_DIR}/include/ies
    ${CMAKE_CURRENT_SOURCE_DIR}/test
    )
