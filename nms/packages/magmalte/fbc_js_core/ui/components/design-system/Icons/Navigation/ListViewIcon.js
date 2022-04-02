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
 * @flow strict-local
 * @format
 */

import type {SvgIconStyleProps} from '../SvgIcon';

import React from 'react';
import SvgIcon from '../SvgIcon';

const ListViewIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M19 0a1 1 0 011 1v18a1 1 0 01-1 1H1a1 1 0 01-1-1V1a1 1 0 011-1h18zm-1 2H2v16h16V2zM7 13v2H5v-2h2zm8 0v2H9v-2h6zM7 9v2H5V9h2zm8 0v2H9V9h6zM7 5v2H5V5h2zm8 0v2H9V5h6z" />
  </SvgIcon>
);

export default ListViewIcon;
