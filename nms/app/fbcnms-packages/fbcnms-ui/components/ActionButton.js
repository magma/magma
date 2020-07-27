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

import AddIcon from './design-system/Icons/Actions/AddIcon';
import IconButton from './design-system/IconButton';
import React from 'react';
import RemoveIcon from './design-system/Icons/Actions/RemoveIcon';

type Props = $ReadOnly<{|
  action: 'add' | 'remove',
  onClick: () => void,
|}>;

const ActionButton = (props: Props) => {
  const {action, onClick} = props;

  switch (action) {
    case 'add':
      return <IconButton icon={AddIcon} onClick={onClick} />;
    case 'remove':
      return <IconButton icon={RemoveIcon} skin="gray" onClick={onClick} />;
    default:
      return null;
  }
};

export default ActionButton;
