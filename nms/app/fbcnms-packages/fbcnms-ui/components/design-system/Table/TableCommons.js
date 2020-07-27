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
 *
 * @flow
 * @format
 */
import {makeStyles} from '@material-ui/styles';

const styles = {
  cell: {
    padding: '4px 8px',
    minHeight: '48px',
    height: '48px',
    width: '100%',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    boxSizing: 'border-box',
    '&:first-child': {
      paddingLeft: '12px',
    },
  },
};

export const useTableCommonStyles = makeStyles<{}, typeof styles>(() => styles);
