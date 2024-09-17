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

// Contains utilities for working with SubscriberIDs
package protos

import (
	"fmt"
	"strconv"
	"strings"
)

// Function to convert the Sid string to the proto struct
func SidProto(sid string) (*SubscriberID, error) {
	sidUpper := strings.ToUpper(sid)
	for typ, name := range SubscriberID_IDType_name {
		if strings.HasPrefix(sidUpper, name) {
			return &SubscriberID{
				Type: SubscriberID_IDType(typ),
				Id:   sid[len(name):]}, nil
		}
	}
	return nil, fmt.Errorf("Invalid sid string: %s", sid)
}

// Function to convert the Sid proto struct to its string representation
func SidString(pb *SubscriberID) string {
	if pb == nil {
		return ""
	}
	return SubscriberID_IDType_name[int32(pb.Type)] + pb.Id
}

// SidFromString converts the string representation of Sid to Sid proto struct
func SidFromString(sid string) *SubscriberID {
	for nm, typ := range SubscriberID_IDType_value {
		if strings.HasPrefix(sid, nm) {
			return &SubscriberID{
				Id:   strings.TrimPrefix(sid, nm),
				Type: SubscriberID_IDType(typ),
			}
		}
	}
	return nil
}

// GetIMSIFromSessionId extracts IMSI from a sessionId and returns only the IMSI without prefix
// SessionId format is is considered to be IMMSIxxxxxx-1234, where xxxxx is the imsi to be extracted
// ie:  IMSI123456789012345-54321   ->  123456789012345
func GetIMSIFromSessionId(sessionId string) (string, error) {
	sessionId = strings.TrimPrefix(sessionId, "IMSI")
	return GetIMSIwithPrefixFromSessionId(sessionId)
}

// GetIMSIwithPrefixFromSessionId extracts IMSI from a sessionId and returns the IMSI with prefix
// SessionId format is is considered to be IMMSIxxxxxx-1234, where xxxxx is the imsi to be extracted
// ie:  IMSI123456789012345-54321   ->  IMSI123456789012345
func GetIMSIwithPrefixFromSessionId(sessionId string) (string, error) {
	data := strings.Split(sessionId, "-")
	if len(data) != 2 {
		return "", fmt.Errorf("Session ID %s does not match format 'IMSI-RandNum'", sessionId)
	}
	return data[0], nil
}

// StripPrefixFromIMSIandFormat extracts IMSI from an IMSI with prefix. It also checks that the IMSI is only numeric
// ie:  IMSI123456789012345   ->  123456789012345
// It returns IMSI in two formats: string and uint64
func StripPrefixFromIMSIandFormat(imsiWithPrefix string) (string, uint64, error) {
	imsiNoPrefix := strings.TrimPrefix(imsiWithPrefix, "IMSI")
	imsiUint, err := strconv.ParseUint(imsiNoPrefix, 10, 64)
	if err != nil {
		return "", 0, fmt.Errorf("IMSI is not numeric: %s", imsiWithPrefix)
	}
	return imsiNoPrefix, imsiUint, nil
}
