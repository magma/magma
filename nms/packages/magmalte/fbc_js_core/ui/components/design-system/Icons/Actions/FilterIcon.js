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

const FilterIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M6 6v1l3.791 2.844L11.608 18h.792l1.809-8.156L18 7V6H6zM4 4h16v3.5a1 1 0 01-.4.8L16 11l-1.995 9h-4L8 11 4.4 8.3a1 1 0 01-.4-.8V4z" />
  </SvgIcon>
);

export default FilterIcon;
