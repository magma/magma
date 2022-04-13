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

const InfoTinyIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M12 18a6 6 0 100-12 6 6 0 000 12zm-1-7h2v5h-2v-5zm0-3h2v2h-2V8z" />
  </SvgIcon>
);

export default InfoTinyIcon;
