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

import type {TRefFor} from '../types/TRefFor.flow';

import * as React from 'react';
import {useEffect, useRef} from 'react';

type Props = {
  isVisible: boolean,
  children: (ref: TRefFor<?HTMLElement>) => React.Node,
  onOutsideClick: () => void,
};

const OnOutsideClickNode = ({isVisible, children, onOutsideClick}: Props) => {
  const elementRef = useRef(null);

  useEffect(() => {
    if (!isVisible) {
      return;
    }

    const listener = (e: MouseEvent) => {
      const node = elementRef.current;
      if (node == null) {
        return;
      }

      const target = e.target;
      if (node instanceof Node && target instanceof Node) {
        if (!node.contains(target)) {
          e.stopPropagation();
          onOutsideClick();
        }
      }
    };

    document.addEventListener('click', listener, true);

    return () => {
      document.removeEventListener('click', listener, true);
    };
  }, [isVisible, onOutsideClick]);

  return <>{children(elementRef)}</>;
};

export default OnOutsideClickNode;
