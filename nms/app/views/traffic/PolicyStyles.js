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
 *
 * @flow strict-local
 * @format
 */
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {colors} from '../../theme/default';

export const policyStyles = {
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
  root: {
    '&$expanded': {
      minHeight: 'auto',
    },
    marginTop: '0px',
    marginBottom: '0px',
  },
  expanded: {marginTop: '-8px', marginBottom: '-8px'},
  block: {
    display: 'block',
  },
  flex: {display: 'flex'},
  panel: {flexGrow: 1},
  removeIcon: {alignSelf: 'baseline'},
  dialog: {height: '640px'},
  title: {textAlign: 'center', margin: 'auto', marginLeft: '0px'},
  description: {
    color: colors.primary.mirage,
  },
  switch: {margin: 'auto 0px'},
};
