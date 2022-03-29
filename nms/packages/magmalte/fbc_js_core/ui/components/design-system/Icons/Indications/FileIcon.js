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

const FileIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path
      d="M15 3l5 6v11a1 1 0 01-1 1H5a1 1 0 01-1-1V4a1 1 0 011-1h10zm-2 2H6v14h12v-9h-5V5zm3 10v2H8v-2h8zm0-3v2H8v-2h8zm-1-5.875V8h1.563L15 6.125z"
      fillRule="nonzero"
    />
  </SvgIcon>
);

export default FileIcon;
