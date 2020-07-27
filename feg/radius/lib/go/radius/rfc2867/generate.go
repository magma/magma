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

//go:generate go run ../cmd/radius-dict-gen/main.go -package rfc2867 -output generated.go -ref Acct-Status-Type:fbc/lib/go/radius/rfc2866 /usr/share/freeradius/dictionary.rfc2867

package rfc2867
