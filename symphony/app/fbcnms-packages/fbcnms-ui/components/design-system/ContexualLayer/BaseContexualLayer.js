/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import * as React from 'react';
import Portal from '../Core/Portal';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useLayoutEffect, useRef, useState} from 'react';

const useStyles = makeStyles(() => ({
  root: {
    position: 'absolute',
    zIndex: 1400,
  },
}));

export type ContextualLayerPosition = 'below' | 'above';

export type ContextualLayerOptions = $ReadOnly<{|
  align?: 'middle' | 'stretch',
  position?: ContextualLayerPosition,
|}>;

type PositionRect = {
  bottom: number,
  left: number,
  right: number,
  top: number,
};

function getElementPosition(element: Element): PositionRect {
  const rect = element.getBoundingClientRect();
  return {
    bottom: rect.bottom,
    left: rect.left,
    right: rect.right,
    top: rect.top,
  };
}

type Props = {
  ...ContextualLayerOptions,
  children: React.Node,
  context: Element,
  hidden?: boolean,
};

const BaseContexualLayer = ({
  children,
  position,
  context,
  hidden = false,
  align = 'middle',
}: Props) => {
  const classes = useStyles();

  const [_hasCalculated, setHasCalculated] = useState(false);
  const contextualLayerRef = useRef<HTMLDivElement | null>(null);

  const recalculateStyles = useCallback(() => {
    const contextRect = getElementPosition(context);
    const getPositioningStyles = () => {
      const style = {};
      switch (position) {
        case 'below':
          style.top = contextRect.bottom;
          style.left = contextRect.left;
          break;
        case 'above':
          style.left = contextRect.left;
          style.top = contextRect.top;
          style.transform = 'translate(0, -100%)';
          break;
      }
      if (align === 'stretch') {
        style.width = contextRect.right - contextRect.left;
      }

      return style;
    };

    const contextualLayerElement = contextualLayerRef.current;
    if (contextualLayerElement != null) {
      contextualLayerElement.removeAttribute('style'); // Unset all styles
      const style = getPositioningStyles();
      Object.keys(style).forEach(styleKey => {
        const value = style[styleKey];
        contextualLayerElement.style.setProperty(
          styleKey,
          typeof value === 'number' ? String(value) + 'px' : value,
        );
      });
    }
    setHasCalculated(true);
  }, [context, position, align]);

  useLayoutEffect(() => {
    if (!hidden) {
      recalculateStyles();
    }
  }, [recalculateStyles, hidden]);

  const contextualLayerFunctionRef = useCallback(
    contextualLayerElement => {
      contextualLayerRef.current = contextualLayerElement;
      recalculateStyles();
    },
    [recalculateStyles],
  );

  return (
    <Portal target={document.body}>
      <div
        className={classes.root}
        ref={contextualLayerFunctionRef}
        hidden={hidden}>
        {children}
      </div>
    </Portal>
  );
};

export default BaseContexualLayer;
