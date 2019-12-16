/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
'use strict';

export type TRefCallbackFor<T> = (T | null) => mixed;
export type TRefObjectFor<T> = {current: T, ...};

// NOTE:
// A simple utility type for declaring ref types.
// Please note, remember to use a nullable version of this in functions that are
// wrapped by React.forwardRef.
export type TRefFor<TElement> =
  | TRefObjectFor<TElement | null>
  | TRefCallbackFor<TElement>;
