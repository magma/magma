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

package encoding

func getBearerActivationSpecificParams(eventData map[string]interface{}) EPSSpecificParameters {
	apn := eventData["apn"].(string)
	sessionID := eventData["session_id"].(string)
	userLocation := eventData["user_location"].(string)
	activationType := DefaultBearer // Not forwarded yet
	apnAmbr := "test"               // Not forwarded yet
	bearerQos := "test"             // Not forwarded yet

	addressAllocation := []byte("")
	if activationType == DefaultBearer {
		addressAllocation = getPdnAddressAllocation(eventData)
	}

	return EPSSpecificParameters{
		PDNAddressAllocation: []byte(addressAllocation),
		APN:                  []byte(apn),
		EPSBearerIdentity:    []byte(sessionID),
		RATType:              []byte{byte(RatTypeEutran)},
		EPSBearerQoS:         []byte(bearerQos),
		BearerActivationType: activationType,
		ApnAmbr:              []byte(apnAmbr),
		EPSLocationOfTheTarget: EPSLocation{
			UserLocationInfo: []byte(userLocation),
		},
	}
}

func getStartInterceptWithActiveBearerSpecificParams(eventData map[string]interface{}) EPSSpecificParameters {
	return EPSSpecificParameters{}
}
