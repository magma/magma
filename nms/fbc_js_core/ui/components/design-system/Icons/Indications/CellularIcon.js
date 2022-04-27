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

const CellularIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M7 16v4H5v-4h2zm4-4v8H9v-8h2zm4-4v12h-2V8h2zm4-4v16h-2V4h2z" />
  </SvgIcon>
);

export default CellularIcon;
