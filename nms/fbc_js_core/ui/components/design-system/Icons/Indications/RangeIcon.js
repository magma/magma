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

import type {SvgIconStyleProps} from '../SvgIcon';

import React from 'react';
import SvgIcon from '../SvgIcon';

const RangeIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M7 6.343l1.414 1.414L4.172 12l4.242 4.243L7 17.657l-4.243-4.243L1.343 12 7 6.343zm10 0l4.243 4.243L22.657 12 17 17.657l-1.414-1.414L19.828 12l-4.242-4.243L17 6.343zM16 11v2H8v-2h8z" />
  </SvgIcon>
);

export default RangeIcon;
