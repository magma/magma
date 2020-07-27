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

const DeleteIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M14 2a1 1 0 011 1v1h4a1 1 0 011 1v5h-1v11a1 1 0 01-1 1H6a1 1 0 01-1-1V10H4V5a1 1 0 011-1h4V3a1 1 0 011-1h4zm2.999 8h-10L7 20h10l-.001-10zM11 12v6H9v-6h2zm4 0v6h-2v-6h2zm3-4V6H6v2h12z" />
  </SvgIcon>
);

export default DeleteIcon;
