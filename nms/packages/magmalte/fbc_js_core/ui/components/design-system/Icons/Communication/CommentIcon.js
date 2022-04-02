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

const CommentIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M4 19.263L7.394 17H20V7H4v12.263zM3 5h18a1 1 0 011 1v12a1 1 0 01-1 1H8l-6 4V6a1 1 0 011-1z" />
  </SvgIcon>
);

export default CommentIcon;
