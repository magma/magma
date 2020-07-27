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

      followingElement.style.top = `${
        followedElementTop + followedElementHeight
      }px`;

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
