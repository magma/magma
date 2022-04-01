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
import {useCallback, useEffect} from 'react';

type Props = {
  children: React.Node,
  onEsc: () => void,
};

const HideOnEsc = ({children, onEsc}: Props) => {
  const onKeyUp = useCallback(
    (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        onEsc();
      }
    },
    [onEsc],
  );

  useEffect(() => {
    document.addEventListener('keyup', onKeyUp);
    return () => document.removeEventListener('keyup', onKeyUp);
  }, [onKeyUp]);

  return <>{children}</>;
};

export default HideOnEsc;
