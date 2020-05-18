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
import emptyFunction from '@fbcnms/util/emptyFunction';
import {createContext, useCallback, useContext, useState} from 'react';
import {Map as immMap} from 'immutable';

export type HierarchyContextValue = $ReadOnly<{|
  childrenValues: immMap<string, ?boolean>,
  parentValue: ?boolean,
  setChildValue: (key: string, value: ?boolean) => void,
|}>;

const DEFUALT_VALUE = {
  childrenValues: new immMap<string, ?boolean>(),
  parentValue: null,
  setChildValue: emptyFunction,
};

const HierarchyContext = createContext<HierarchyContextValue>(DEFUALT_VALUE);

export function useHierarchyContext() {
  return useContext(HierarchyContext);
}

export type HierarchyContextProviderProps = $ReadOnly<{|
  parentValue: ?boolean,
  children: ?React.Node,
|}>;

export function HierarchyContextProvider(props: HierarchyContextProviderProps) {
  const {parentValue, children} = props;
  const [childrenValues, setChildrenValues] = useState(
    new immMap<string, ?boolean>(),
  );

  const setChildValue = useCallback((key: string, value: ?boolean) => {
    setChildrenValues(currentChildren => currentChildren.set(key, value));
  }, []);

  const value = {
    childrenValues,
    parentValue,
    setChildValue,
  };

  return (
    <HierarchyContext.Provider value={value}>
      {children}
    </HierarchyContext.Provider>
  );
}

export default HierarchyContext;
