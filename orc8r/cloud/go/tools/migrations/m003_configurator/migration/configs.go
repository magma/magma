/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package migration

import (
	"bytes"

	"github.com/golang/glog"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

var migratorRegistry = map[string]ConfigMigrator{}

func migrateConfig(configType string, configVal []byte) ([]byte, error) {
	migrator, found := migratorRegistry[configType]
	if !found {
		glog.Infof("no migrator found for config type %s, skipping", configType)
		return nil, nil
	}
	return migrator.ToNewConfig(configVal)
}

// THIS CODE HAS BEEN DUPLICATED FROM ANOTHER LOCATION
// DO NOT --EVER-- CHANGE THIS CODE, EXCEPT TO DELETE THE ENTIRE MIGRATION

func Unmarshal(bt []byte, msg proto.Message) error {
	return (&jsonpb.Unmarshaler{AllowUnknownFields: true}).Unmarshal(
		bytes.NewBuffer(bt),
		msg)
}
