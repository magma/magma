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
 */

import {OptionsObject} from 'notistack';

export default function useSnackbar(
  message: string,
  config: OptionsObject,
  show: boolean,
  dismissPrevious?: boolean,
): void;

export function useSnackbars(): {
  success: (message: string) => void;
  error: (message: string) => void;
  warning: (message: string) => void;
};

export function useEnqueueSnackbar(): (
  message: string,
  config: OptionsObject,
) => void;
