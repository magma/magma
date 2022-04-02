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

const ArrowUp = ({className}: Props) => (
  <svg width="10px" height="10px" xmlns="http://www.w3.org/2000/svg">
    <g className={className} fill="none">
      <path d="M17 17H-7V-7h24z" />
      <path
        d="M5.556 2.084V10H4.444V2.084L.812 5.556 0 4.779 5 0l5 4.78-.812.776-3.632-3.472z"
        fill={SymphonyTheme.palette.primary}
      />
    </g>
  </svg>
);

export default ArrowUp;
