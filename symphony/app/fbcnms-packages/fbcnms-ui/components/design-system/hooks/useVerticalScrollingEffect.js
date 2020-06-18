/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import {useEffect} from 'react';

type VerticalScrollValues = {
  scrollHeight: number,
  scrollTop: number,
  height: number,
  startIsHidden: boolean,
  endIsHidden: boolean,
  hasScrollbar: boolean,
  scrollbarWidth: number,
};

const getVerticalScrollValues: HTMLElement => VerticalScrollValues = (
  htmlElement: HTMLElement,
) => {
  if (!htmlElement) {
    return {};
  }

  return {
    hasScrollbar: htmlElement.offsetHeight < htmlElement.scrollHeight,
    scrollHeight: htmlElement.scrollHeight,
    scrollTop: htmlElement.scrollTop,
    height: htmlElement.offsetHeight,
    startIsHidden: htmlElement.scrollTop > 0,
    endIsHidden:
      htmlElement.offsetHeight + htmlElement.scrollTop <
      htmlElement.scrollHeight,
    scrollbarWidth: htmlElement.offsetWidth - htmlElement.clientWidth,
  };
};

const BOX_SHADOW_COLOR = 'rgba(0, 0, 0, 0.17)';
const topBoxShadow = 'inset 0px 6px 4px -5px ' + BOX_SHADOW_COLOR;
const bottomBoxShadow = 'inset 0px -6px 4px -5px ' + BOX_SHADOW_COLOR;

const calcVerticalScrollingEffect = (
  scrollingContainerElement,
  effect,
  applyScrollingShadows,
) =>
  window.requestAnimationFrame(() => {
    const scrollValues = getVerticalScrollValues(scrollingContainerElement);
    if (applyScrollingShadows) {
      const scrollBoxShadow = [];
      if (scrollValues.startIsHidden) {
        scrollBoxShadow.push(topBoxShadow);
      }
      if (scrollValues.endIsHidden) {
        scrollBoxShadow.push(bottomBoxShadow);
      }
      scrollingContainerElement.style.boxShadow = scrollBoxShadow.join(', ');
    }
    if (effect) {
      effect(scrollValues);
    }
  });

const useVerticalScrollingEffect = (
  element: {current: HTMLElement | null, ...},
  effect?: VerticalScrollValues => void,
  applyScrollingShadows: boolean = true,
) => {
  useEffect(() => {
    if (!element || !element.current || (!effect && !applyScrollingShadows)) {
      return;
    }
    const scrollingContainer = element.current;
    const runEffect = () =>
      calcVerticalScrollingEffect(
        scrollingContainer,
        effect,
        applyScrollingShadows,
      );
    runEffect();
    scrollingContainer.addEventListener('scroll', runEffect);
    return () => scrollingContainer.removeEventListener('scroll', runEffect);
  });
};

export default useVerticalScrollingEffect;
