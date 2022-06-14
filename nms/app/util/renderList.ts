/**
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

export default function renderList(list: Array<string>): string {
  if (!Array.isArray(list)) {
    console.error(
      // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
      `renderList(): expected array, received ${list} (${typeof list})`,
    );
    return '';
  }

  if (list.length < 4) {
    return list.join(', ');
  }

  return `${list[0]}, ${list[1]} & ${list.length - 2} others`;
}
