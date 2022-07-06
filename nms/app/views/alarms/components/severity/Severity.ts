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

import grey from '@material-ui/core/colors/grey';
import orange from '@material-ui/core/colors/orange';
import red from '@material-ui/core/colors/red';
import yellow from '@material-ui/core/colors/yellow';

export const SEVERITY = {
  NOTICE: {
    name: 'NOTICE',
    order: 0,
    color: grey[500],
  },
  INFO: {
    name: 'INFO',
    order: 1,
    color: grey[500],
  },
  WARNING: {
    name: 'WARNING',
    order: 2,
    color: yellow.A400,
  },
  MINOR: {
    name: 'MINOR',
    order: 3,
    color: yellow.A400,
  },
  MAJOR: {
    name: 'MAJOR',
    order: 4,
    color: orange.A400,
  },
  CRITICAL: {
    name: 'CRITICAL',
    order: 5,
    color: red.A400,
  },
} as const;
