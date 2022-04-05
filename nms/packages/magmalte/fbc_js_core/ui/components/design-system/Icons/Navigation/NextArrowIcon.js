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

const NextArrowIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M13 6l-1.41 1.41L15.173 11H5v2h10.173l-3.583 3.59L13 18l5.293-5.293a1 1 0 000-1.414L13 6z" />
  </SvgIcon>
);

export default NextArrowIcon;
