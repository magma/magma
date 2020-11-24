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

package serdes

import (
	"magma/fbinternal/cloud/go/services/testcontroller"
	"magma/fbinternal/cloud/go/services/testcontroller/obsidian/models"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/serdes"
)

const (
	testcontrollerSerdeDomain = "testcontroller"
)

var (
	// Network contains the full set of configurator network config serdes
	// used in the fbinternal module
	Network        = serdes.Network.MustMerge(models.NetworkSerdes)
	TestController = serde.NewRegistry(
		serde.NewBinarySerde(testcontrollerSerdeDomain, testcontroller.EnodedTestCaseType, &models.EnodebdTestConfig{}),
		serde.NewBinarySerde(testcontrollerSerdeDomain, testcontroller.EnodedTestExcludeTraffic, &models.EnodebdTestConfig{}),
	)
)
