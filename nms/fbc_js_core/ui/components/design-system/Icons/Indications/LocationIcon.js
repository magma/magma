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

const LocationIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M12 1c4.866 0 8.847 3.868 8.996 8.735l.004.268c-.012 4.469-2.731 8.65-8.047 12.542l-.939.678-.58-.399C5.774 18.934 2.918 14.65 3 9.984l.004-.264A9 9 0 0112 1zm0 2C8.215 3 5.119 6.009 5.004 9.766L5 10.017c-.065 3.75 2.225 7.33 6.99 10.753l.366-.276c4.454-3.41 6.634-6.908 6.644-10.478l-.003-.235A7 7 0 0012 3zm0 3a4 4 0 110 8 4 4 0 010-8zm0 2a2 2 0 100 4 2 2 0 000-4z" />
  </SvgIcon>
);

export default LocationIcon;
