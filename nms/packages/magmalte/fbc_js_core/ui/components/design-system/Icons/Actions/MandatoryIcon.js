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

const MandatoryIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M13 3v7.267l6.294-3.633 1 1.732L14 12l6.294 3.634-1 1.732L13 13.732V21h-2v-7.268l-6.294 3.634-1-1.732 6.293-3.635-6.293-3.633 1-1.732L11 10.267V3h2z" />
  </SvgIcon>
);

export default MandatoryIcon;
