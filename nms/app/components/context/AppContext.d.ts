/**
 * Copyright 2022 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
import React from 'react';

type Props = {
  children: React.ReactNode;
  isOrganizations?: boolean;
  networkIDs?: Array<string>;
};

export function AppContextProvider(props: Props): JSX.Element;

declare const Context: React.Context<any>;

export default Context;
