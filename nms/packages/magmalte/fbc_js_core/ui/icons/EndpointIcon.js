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
import SvgIcon from '@material-ui/core/SvgIcon';

type Props = {
  className?: string,
};

const EndpointIcon = (props: Props) => (
  <SvgIcon color="inherit" viewBox="0 0 16 18" className={props.className}>
    <g
      fill="none"
      fillRule="evenodd"
      transform="translate(0 1)"
      stroke="#B8C2D3"
      strokeWidth="2">
      <path
        d="m9.5 0.86603 3.9282 2.2679c0.9282 0.5359 1.5 1.5263 1.5 2.5981v4.5359c0 1.0718-0.5718 2.0622-1.5 2.5981l-3.9282 2.2679c-0.9282 0.5359-2.0718 0.5359-3 0l-3.9282-2.2679c-0.9282-0.5359-1.5-1.5263-1.5-2.5981v-4.5359c0-1.0718 0.5718-2.0622 1.5-2.5981l3.9282-2.2679c0.9282-0.5359 2.0718-0.5359 3 0z"
        fill="#EDF0F9"
      />
      <polyline points="14.857 4.1905 8 8 1.1429 4.1905" />
      <path d="m8 8v7" strokeLinecap="square" />
    </g>
  </SvgIcon>
);

export default EndpointIcon;
