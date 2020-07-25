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

package libgraphql

import (
	"log"
	"os"
)

// get ACCESS_TOKEN from "ixpp <partner short name>"
func ExampleClient_Do() {
	c := NewClient(ClientConfig{
		Token:    os.Getenv("ACCESS_TOKEN"),
		Endpoint: "https://graph.expresswifi.com/graphql",
	})
	op := NewUpsertCustomer(&AppCustomer{
		MobileNumber: "12311728371117",
	})
	if err := c.Do(op); err != nil {
		log.Fatalf("failed executing graphql request: %v", err)
	}
	log.Printf("graphql response: %v", op.Response())
}
