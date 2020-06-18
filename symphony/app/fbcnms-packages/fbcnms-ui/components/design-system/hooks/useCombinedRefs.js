/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import {useEffect, useRef} from 'react';

type RefFunctionType = (HTMLElement | null) => mixed;
type RefType = RefFunctionType | {current: HTMLElement | null};
export type CombinedRefs = Array<RefType>;

const useCombinedRefs = (refs: CombinedRefs) => {
  const targetRef = useRef<?HTMLElement>(null);

  useEffect(() => {
    refs.forEach((ref: ?RefType) => {
      if (ref == null) {
        return;
      }

      if (typeof ref === 'function') {
        const refFunction: RefFunctionType = ref;
        refFunction(targetRef.current || null);
      } else {
        ref.current = targetRef.current || null;
      }
    });
  }, [refs]);

  return targetRef;
};

export default useCombinedRefs;
