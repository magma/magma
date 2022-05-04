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

const ProjectsIcon = (props: Props) => (
  <SvgIcon
    color="inherit"
    width="22px"
    height="19px"
    className={props.className}>
    <g transform="translate(-1 -2)" fill="none" fillRule="evenodd">
      <rect
        stroke="currentColor"
        strokeWidth="2"
        x="10"
        y="3"
        width="4"
        height="4"
        rx="1"
      />
      <rect
        stroke="currentColor"
        strokeWidth="2"
        x="10"
        y="16"
        width="4"
        height="4"
        rx="1"
      />
      <rect
        stroke="currentColor"
        strokeWidth="2"
        x="18"
        y="16"
        width="4"
        height="4"
        rx="1"
      />
      <rect
        stroke="currentColor"
        strokeWidth="2"
        x="2"
        y="16"
        width="4"
        height="4"
        rx="1"
      />
      <path
        d="M13 10h8v6h-2v-4h-6v4h-2v-4H5v4H3v-6h8V7h2v3z"
        fill="currentColor"
      />
    </g>
  </SvgIcon>
);

export default ProjectsIcon;
