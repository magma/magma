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

const INTERVAL = 500;

const stickfollowingElement = (
  followingElement: HTMLElement | null,
  followedElement: HTMLElement | null,
): Promise<void> => {
  const p = new Promise<void>(resolve =>
    window.requestAnimationFrame(() => {
      if (!followingElement || !followedElement) {
        resolve();
        return;
      }

      const followedElementDimensions = followedElement.getClientRects()[0];
      const followedElementTop = followedElementDimensions.top;
      const followedElementHeight = followedElementDimensions.height;

      followingElement.style.top = `${followedElementTop +
        followedElementHeight}px`;

      resolve();
    }),
  );

  return p;
};

const useFollowElement = (
  followingElement: {current: HTMLElement | null, ...},
  followedElement: {current: HTMLElement | null, ...},
) => {
  useEffect(() => {
    if (!followingElement?.current || !followedElement?.current) {
      return;
    }

    const followingElementStyle = window.getComputedStyle(
      followingElement.current,
    );
    if (
      followingElementStyle.display === 'none' ||
      followingElementStyle.visibility !== 'visible'
    ) {
      return;
    }

    let timeoutId = null;
    const runEffect = () => {
      timeoutId = setTimeout(() => {
        stickfollowingElement(
          followingElement.current,
          followedElement.current,
        ).finally(() => runEffect());
      }, INTERVAL);
    };
    runEffect();
    return () => clearTimeout(timeoutId);
  });
};

export default useFollowElement;
