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
import {useState} from 'react';

type AddNewCategoryMethod = () => void;
type MethodsType = {
  addNewCategory: AddNewCategoryMethod,
};
type ExtractReturnObjectType = <V>(V) => V => void;
type MethodsContextType<M> = {
  call: M,
  override: $ObjMap<M, ExtractReturnObjectType>,
};

const CheckListCategoryContext = React.createContext<
  MethodsContextType<MethodsType>,
>({
  call: {
    addNewCategory: emptyFunction,
  },
  override: {
    addNewCategory: emptyFunction,
  },
});

type Props = {
  children: React.Node,
};

export function CheckListCategoryContextProvider(props: Props) {
  const [addNewCategory, setAddNewCategory] = useState(() => emptyFunction);
  const providerValue = {
    call: {
      addNewCategory,
    },
    override: {
      addNewCategory: newAdd => {
        setAddNewCategory(() => newAdd);
      },
    },
  };
  return (
    <CheckListCategoryContext.Provider value={providerValue}>
      {props.children}
    </CheckListCategoryContext.Provider>
  );
}

export default CheckListCategoryContext;
