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
'use strict';

export type TRefCallbackFor<T> = (T | null) => mixed;
export type TRefObjectFor<T> = {current: T, ...};

// NOTE:
// A simple utility type for declaring ref types.
// Please note, remember to use a nullable version of this in functions that are
// wrapped by React.forwardRef.
export type TRefFor<TElement> =
  | TRefObjectFor<TElement | null>
  | TRefCallbackFor<TElement>;
