/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Contains utilities for working with SubscriberIDs
package protos

import (
	"fmt"
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

// ParseIMSIfromSessionIdNoPrefix extracts IMSI from a sessionId and returns only the IMSI without prefix
// SessionId format is is considered to be IMMSIxxxxxx-1234, where xxxxx is the imsi to be extracted
// ie:  IMSI123456789012345-54321   ->  123456789012345
func ParseIMSIfromSessionIdNoPrefix(sessionId string) (string, error) {
	sessionId = strings.TrimPrefix(sessionId, "IMSI")
	return ParseIMSIfromSessionIdWithPrefix(sessionId)
}

// ParseIMSIfromSessionIdWithPrefix extracts IMSI from a sessionId and returns the IMSI with prefix
// SessionId format is is considered to be IMMSIxxxxxx-1234, where xxxxx is the imsi to be extracted
// ie:  IMSI123456789012345-54321   ->  IMSI123456789012345
func ParseIMSIfromSessionIdWithPrefix(sessionId string) (string, error) {
	data := strings.Split(sessionId, "-")
	if len(data) != 2 {
		return "", fmt.Errorf("Session ID %s does not match format 'IMSI-RandNum'", sessionId)
	}
	return data[0], nil

}
