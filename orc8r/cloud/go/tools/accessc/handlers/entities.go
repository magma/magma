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

// Package handlers implements individual accessc commands as well as common
// across multiple commands functionality
package handlers

import (
	"fmt"
	"regexp"
	"strings"

	accessprotos "magma/orc8r/cloud/go/services/accessd/protos"
)

// Entity is a CLI tool representation of ACLed entity (network, operator, etc.)
type Entity struct {
	id   string
	perm int32
}

// Entities - list of entities compiled from command line flags
type Entities []Entity

// String - stringer for entities
func (ents *Entities) String() string {
	if ents == nil {
		return "<nil>"
	}
	res := ""
	for _, ent := range *ents {
		fmtspec := "%s:%s"
		if len(res) > 0 {
			fmtspec = "; %s:%s"
		}
		res += fmt.Sprintf(
			fmtspec,
			ent.id,
			accessprotos.AccessControl_Permission(ent.perm).ToString())
	}
	return res
}

// Set adds a new Entity into Entities Map from provided flag value string
func (ents *Entities) Set(value string) error {
	sepIdx := strings.LastIndex(value, ":")
	if sepIdx <= 0 {
		return fmt.Errorf(
			"Invalid Entity Specification for '%s', expected <id>:R|W|RW",
			value)
	}
	id := strings.TrimSpace(value[:sepIdx])
	if len(id) == 0 {
		return fmt.Errorf(
			"Invalid Entity Specification for '%s', Id cannot be empty",
			value)
	}
	permStr := strings.ToUpper(strings.TrimSpace(value[sepIdx+1:]))
	// Check if it's in the form: R, W, RW or R&W, R+W, R|W or any combination
	// of them
	match, err := regexp.MatchString("^[RW+&|]+$", permStr)
	if err != nil || (!match) {
		return fmt.Errorf(
			"Invalid Entity Permissions for '%s': %s. Expected <id>:R|W|RW",
			id, permStr)
	}
	perm := accessprotos.AccessControl_NONE
	if strings.Contains(permStr, "R") {
		perm |= accessprotos.AccessControl_READ
	}
	if strings.Contains(permStr, "W") {
		perm |= accessprotos.AccessControl_WRITE
	}
	if perm == accessprotos.AccessControl_NONE {
		return fmt.Errorf(
			"Invalid Entity Specification for '%s', at least one R/W permission must be specified",
			value)
	}
	*ents = append(*ents, Entity{id, int32(perm)})
	return nil
}
