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

import * as React from 'react';

type Props = {
  onClick?: () => void,
  error?: boolean,
  children: React.Node,
};

export default function Button(props: Props) {
  const styles = {
    border: '1px solid #bbb',
    borderRadius: 6,
    cursor: 'pointer',
    fontSize: 15,
    padding: '3px 10px',
  };
  if (props.error != null) {
    styles['border'] = '1px solid red';
  }
  return (
    <button style={styles} onClick={props.onClick}>
      {props.children}
    </button>
  );
}
