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

type InputChangeFunc<TFormState, TVal> = (
  formUpdate: FormUpdate<TFormState, TVal>,
) => (event: SyntheticInputEvent<HTMLElement>) => void;
type FormUpdate<TFormState, TVal = string> = (
  val: TVal,
  event: SyntheticInputEvent<HTMLElement>,
) => $Shape<TFormState>;

export default function useForm<TFormState: {}>({
  initialState,
  onFormUpdated,
}: {
  initialState: TFormState,
  onFormUpdated?: (state: TFormState) => void,
}): {|
  formState: TFormState,
  updateFormState: (update: $Shape<TFormState>) => TFormState,
  handleInputChange: InputChangeFunc<TFormState, *>,
  updateListItem: (
    listName: $Keys<TFormState>,
    idx: number,
    update: $Shape<TFormState> | TFormState,
  ) => void,
  addListItem: (listName: $Keys<TFormState>, item: {}) => void,
  removeListItem: (listName: $Keys<TFormState>, idx: number) => void,
  setFormState: (f: TFormState) => void,
|} {
  const [formState, setFormState] = React.useState<TFormState>(initialState);
  const formUpdatedRef = React.useRef(onFormUpdated);
  React.useEffect(() => {
    formUpdatedRef.current = onFormUpdated;
  }, [onFormUpdated]);
  const updateFormState = React.useCallback(
    update => {
      const nextState = {
        ...formState,
        ...update,
      };
      setFormState(nextState);
      if (typeof formUpdatedRef.current === 'function') {
        formUpdatedRef.current(nextState);
      }
      return nextState;
    },
    [formState, formUpdatedRef, setFormState],
  );

  /**
   * Immutably updates an item in an array on T.
   * usage:
   * //formState: {list: [{x:1},{x:2}]};
   * updateListItem('list', 0, {x:0})
   * //formState: {{list: [{x:0},{x:2}]}}
   */
  const updateListItem = React.useCallback(
    (
      listName: $Keys<TFormState>,
      idx: number,
      update: $Shape<TFormState> | TFormState,
    ) => {
      updateFormState({
        [listName]: immutablyUpdateArray(
          formState[listName] || [],
          idx,
          update,
        ),
      });
    },
    [formState, updateFormState],
  );

  const removeListItem = React.useCallback(
    (listName: $Keys<TFormState>, idx: number) => {
      if (!formState[listName]) {
        return;
      }
      updateFormState({
        [listName]: formState[listName].filter((_, i) => i !== idx),
      });
    },
    [formState, updateFormState],
  );

  const addListItem = React.useCallback(
    <TItem>(listName: $Keys<TFormState>, item: TItem) => {
      updateFormState({
        [listName]: [...(formState[listName] || []), item],
      });
    },
    [formState, updateFormState],
  );
  /**
   * Passes the event value to an updater function which returns an update
   * object to be merged into the form.
   */
  const handleInputChange = React.useCallback(
    (formUpdate: FormUpdate<TFormState>) => (
      event: SyntheticInputEvent<HTMLElement>,
    ) => {
      const value = event.target.value;
      updateFormState(formUpdate(value, event));
    },
    [updateFormState],
  );

  return {
    formState,
    updateFormState,
    handleInputChange,
    updateListItem,
    addListItem,
    removeListItem,
    setFormState,
  };
}

/**
 * Copies array with the element at idx immutably merged with update
 */
function immutablyUpdateArray<T>(
  array: Array<T>,
  idx: number,
  update: $Shape<T>,
) {
  return array.map((item, i) => {
    if (i !== idx) {
      return item;
    }
    return {...item, ...update};
  });
}
