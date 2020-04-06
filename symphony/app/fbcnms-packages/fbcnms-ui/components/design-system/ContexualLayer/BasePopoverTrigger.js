/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {ContextualLayerPosition} from './BaseContexualLayer';
import type {TRefCallbackFor} from '../types/TRefFor.flow';

import * as React from 'react';
import BaseContexualLayer from './BaseContexualLayer';
import HideOnEsc from './HideOnEsc';
import MenuContextProvider from '../Select/MenuContext';
import OnOutsideClickNode from './OnOutsideClickNode';
import {useCallback, useRef, useState} from 'react';

type Props = {
  children: (
    onShow: () => void,
    onHide: () => void,
    contextRef: TRefCallbackFor<?HTMLElement>,
  ) => React.Node,
  popover: React.Node,
  position?: ContextualLayerPosition,
};

const BasePopoverTrigger = ({children, popover, position = 'below'}: Props) => {
  const [isVisible, setIsVisible] = useState(false);
  const contextRef = useRef<?HTMLElement>(null);

  const onHide = useCallback(() => {
    setIsVisible(false);
  }, []);

  const onShow = useCallback(() => {
    if (isVisible) {
      return;
    }

    setIsVisible(true);
  }, [isVisible]);

  const refCallback = useCallback((element: ?HTMLElement) => {
    contextRef.current = element;
  }, []);

  return (
    <>
      {children(onShow, onHide, refCallback)}
      {contextRef.current != null ? (
        <BaseContexualLayer
          context={contextRef.current}
          position={position}
          hidden={!isVisible}>
          <HideOnEsc onEsc={onHide}>
            <MenuContextProvider value={{shown: isVisible, onClose: onHide}}>
              <OnOutsideClickNode isVisible={isVisible} onOutsideClick={onHide}>
                {ref => <div ref={ref}>{popover}</div>}
              </OnOutsideClickNode>
            </MenuContextProvider>
          </HideOnEsc>
        </BaseContexualLayer>
      ) : null}
    </>
  );
};

export default BasePopoverTrigger;
