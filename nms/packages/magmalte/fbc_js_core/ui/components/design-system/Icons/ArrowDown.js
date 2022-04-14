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

import React from 'react';
import SymphonyTheme from '../../../theme/symphony';

type Props = {
  className?: string,
};

const ArrowDown = ({className}: Props) => (
  <svg width="10" height="10" xmlns="http://www.w3.org/2000/svg">
    <path
      className={className}
      d="M4.444 7.916V0h1.112v7.916l3.632-3.472.812.777L5 10 0 5.22l.812-.776 3.632 3.472z"
      fill={SymphonyTheme.palette.primary}
    />
  </svg>
);

export default ArrowDown;
