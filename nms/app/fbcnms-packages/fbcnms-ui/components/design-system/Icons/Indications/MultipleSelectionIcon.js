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

const MultipleSelectionIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M21 16v2H7v-2h14zM5 16v2H3v-2h2zm16-5v2H7v-2h14zM5 11v2H3v-2h2zm16-5v2H7V6h14zM5 6v2H3V6h2z" />
  </SvgIcon>
);

export default MultipleSelectionIcon;
