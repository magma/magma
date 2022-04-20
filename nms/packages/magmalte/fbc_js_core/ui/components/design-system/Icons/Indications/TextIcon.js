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

const TextIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M18.058 9.243c1.082 0 1.955.275 2.62.824.664.55.996 1.34.996 2.37v4.131c0 .384.025.751.077 1.104.051.352.135.697.249 1.034h-2.112a7.701 7.701 0 01-.185-.687 4.911 4.911 0 01-.099-.67c-.286.436-.663.8-1.13 1.095a2.873 2.873 0 01-1.566.443c-.968 0-1.71-.25-2.229-.747-.518-.498-.777-1.183-.777-2.053 0-.899.355-1.598 1.065-2.1.71-.5 1.709-.75 2.997-.75h1.623v-.817c0-.486-.143-.867-.43-1.142-.286-.275-.695-.412-1.227-.412-.476 0-.853.115-1.134.344a1.09 1.09 0 00-.42.884h-2.01l-.009-.051c-.04-.733.288-1.383.983-1.95.696-.567 1.602-.85 2.718-.85zM8.57 6.203l4.543 12.503h-2.121l-.988-2.92H5.11l-.997 2.92H2L6.594 6.203H8.57zm11.018 8.338H17.93c-.607 0-1.082.15-1.426.451-.343.3-.515.654-.515 1.06 0 .356.116.638.348.847.231.209.562.313.991.313.527 0 1-.132 1.422-.395.42-.263.7-.564.837-.902v-1.374zM7.599 8.702h-.052L5.71 14.043h3.7L7.6 8.703z" />
  </SvgIcon>
);

export default TextIcon;
