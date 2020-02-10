/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {TRefObjectFor} from '@fbcnms/ui/components/design-system/types/TRefFor.flow';

import {useCallback, useEffect, useState} from 'react';

export type DimensionResizeAttributes = {
  prev: number,
  new: number,
  expanded: boolean,
};
export type ResizeAttributes = {
  width: DimensionResizeAttributes,
  height: DimensionResizeAttributes,
};

const useResize = (
  element: TRefObjectFor<HTMLElement | null>,
  effect: (att: ResizeAttributes) => void,
) => {
  const [lastWidth, setLastWidth] = useState(0);
  const [lastHeight, setLastHeight] = useState(0);
  const handleDimension = useCallback(
    (
      lastValue: number,
      newValue: number,
      updateStateCallback: number => void,
    ) => {
      updateStateCallback(newValue);
      return {
        prev: lastValue,
        new: newValue,
        expanded: newValue > lastValue,
      };
    },
    [],
  );

  const callEffect = useCallback(() => {
    if (!effect || !element?.current) {
      return;
    }
    const trackedElement = element.current;
    const attr = {
      height: handleDimension(
        lastHeight,
        trackedElement.clientHeight,
        setLastHeight,
      ),
      width: handleDimension(
        lastWidth,
        trackedElement.clientWidth,
        setLastWidth,
      ),
    };
    effect(attr);
  }, [effect, element, handleDimension, lastHeight, lastWidth]);

  useEffect(() => {
    window.addEventListener('resize', callEffect);
    return () => {
      window.removeEventListener('resize', callEffect);
    };
  }, [callEffect]);

  useEffect(() => {
    const trackedElement = element?.current;
    if (!trackedElement) {
      return;
    }
    trackedElement.addEventListener('resize', callEffect);
    callEffect();
    return () => {
      trackedElement.removeEventListener('resize', callEffect);
    };
  }, [callEffect, element]);
};

export default useResize;
