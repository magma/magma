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

//go:generate bash -c "protoc -I /usr/include -I $MAGMA_ROOT --go_out=plugins=grpc,Mgoogle/protobuf/field_mask.proto=google.golang.org/genproto/protobuf/field_mask:$MAGMA_ROOT/.. $MAGMA_ROOT/lte/protos/oai/*.proto"
package oai
