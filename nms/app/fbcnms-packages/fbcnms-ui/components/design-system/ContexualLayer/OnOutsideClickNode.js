/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
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
