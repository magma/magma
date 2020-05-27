/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {useEffect, useState} from 'react';

/*
  This hooks helps calling update callback prop.
  It is needed in case we want to call the update callback
  only AFTER context values were recalculated.
  Without it, the callback is called BEFORE the context
  value changes takes effect.
  - SetTimeout is needed for ensuring the used callback is
    surly the one passed AFTER the context values calculations.
  - Using state hook for engaging rendring cycle when triggerred
    (without it, will be using the previous version of given callback).
  - Using effect hook for completing rendering cycle before using callback.
*/

export default function useSideEffectCallback(callback: ?() => void) {
  const [shouldTriggerRunCallback, setShouldTriggerRunCallback] = useState(
    false,
  );
  useEffect(() => {
    if (!shouldTriggerRunCallback) {
      return;
    }
    setShouldTriggerRunCallback(false);
    if (callback == null) {
      return;
    }
    callback();
  }, [callback, shouldTriggerRunCallback]);

  return () => setTimeout(() => setShouldTriggerRunCallback(true));
}
