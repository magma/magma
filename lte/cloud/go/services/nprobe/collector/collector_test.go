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

package collector

import (
	"testing"

	"magma/lte/cloud/go/services/nprobe"

	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
)

func TestMultiStreamsQuery(t *testing.T) {

	expectedQuery := elastic.NewBoolQuery().
		Filter(elastic.NewTermsQuery("stream_name.keyword", stringsToInterfaces(nprobe.GetESStreams())...)).
		Filter(elastic.NewTermsQuery("event_type.keyword", stringsToInterfaces(nprobe.GetESEventTypes())...)).
		Filter(elastic.NewTermsQuery("event_tag.keyword", "001010000000001"))

	params := getMultiStreamsQueryParameters("test", "", []string{"001010000000001"})
	query := params.toElasticBoolQuery()

	assert.Equal(t, expectedQuery, query)
}
