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

const CheckIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M10.292 16.293a.996.996 0 01-1.41 0l-3.59-3.59a.996.996 0 111.41-1.41l2.88 2.88 6.88-6.88a.996.996 0 111.41 1.41l-7.58 7.59z" />
  </SvgIcon>
);

export default CheckIcon;
