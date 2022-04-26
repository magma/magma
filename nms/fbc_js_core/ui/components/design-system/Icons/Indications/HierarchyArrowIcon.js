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

const HierarchyArrowIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M9 7v5l3.999-.001L13 9l5 4-5 4-.001-3.001L7 14V7h2z" />
  </SvgIcon>
);

export default HierarchyArrowIcon;
