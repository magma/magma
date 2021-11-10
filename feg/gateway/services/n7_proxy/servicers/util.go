/*
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package servicers

import (
	"fmt"
	"net/url"
)

func GetServerStringFromUrl(urlStr string) (string, error) {
	urlDef, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return "", err
	}
	serverStr := fmt.Sprintf("%s://%s", urlDef.Scheme, urlDef.Host)
	return serverStr, nil
}
