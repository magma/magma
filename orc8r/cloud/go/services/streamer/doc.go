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

// Package streamer provides a logical stream for orc8r to push updates
// to gateways.
//
// Orc8r services can implement the StreamProvider servicer interface to
// provide a named stream. Streamer forwards requests for updates under a
// specific stream name to the appropriate remote orc8r service.
//
// E.g., consider a gateway requesting subscriber updates, from the
// "subscriber" stream. This takes the form
//
//		gateway -(a)-> streamer -(b)-> lte
//
// (a) GetUpdates("subscriber")
// (*) streamer: look up service name of provider for "subscriber" stream
// (b) GetUpdates("subscriber")
package streamer

const (
	ServiceName = "STREAMER"
)
