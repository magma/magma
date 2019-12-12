/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';
import {useCallback, useEffect} from 'react';

type Props = {
  children: React.Node,
  onEsc: () => void,
};

const HideOnEsc = ({children, onEsc}: Props) => {
  const onKeyUp = useCallback(
    (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        onEsc();
      }
    },
    [onEsc],
  );

  useEffect(() => {
    document.addEventListener('keyup', onKeyUp);
    return () => document.removeEventListener('keyup', onKeyUp);
  }, [onKeyUp]);

  return <>{children}</>;
};

export default HideOnEsc;
