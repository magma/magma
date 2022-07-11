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
 */

import * as React from 'react';

type InputChangeFunc<TFormState, TVal> = (
  formUpdate: FormUpdate<TFormState, TVal>,
) => (event: React.ChangeEvent<HTMLInputElement>) => void;
type FormUpdate<TFormState, TVal = string> = (
  val: TVal,
  event: React.ChangeEvent<HTMLInputElement>,
) => Partial<TFormState>;

export default function useForm<TFormState extends Record<string, any>>({
  initialState,
  onFormUpdated,
}: {
  initialState: Partial<TFormState>;
  onFormUpdated?: (state: TFormState) => void;
}): {
  formState: TFormState;
  updateFormState: (update: Partial<TFormState>) => TFormState;
  handleInputChange: InputChangeFunc<TFormState, any>;
  updateListItem: (
    listName: keyof TFormState,
    idx: number,
    update: Partial<TFormState[keyof TFormState][number]>,
  ) => void;
  addListItem: (listName: keyof TFormState, item: object) => void;
  removeListItem: (listName: keyof TFormState, idx: number) => void;
  setFormState: (f: TFormState) => void;
} {
  // TODO[TS-migration] is formState a Partial TFormState?
  const [formState, setFormState] = React.useState<TFormState>(
    initialState as TFormState,
  );
  const formUpdatedRef = React.useRef(onFormUpdated);
  React.useEffect(() => {
    formUpdatedRef.current = onFormUpdated;
  }, [onFormUpdated]);
  const updateFormState = React.useCallback(
    (update: Partial<TFormState>) => {
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
    <K extends keyof TFormState, VALUE extends TFormState[K] & Array<any>>(
      listName: K,
      idx: number,
      update: Partial<VALUE[number]>,
    ) => {
      updateFormState({
        [listName]: immutablyUpdateArray(
          (formState[listName] || []) as VALUE,
          idx,
          update,
        ),
      } as Partial<TFormState>);
    },
    [formState, updateFormState],
  );

  const removeListItem = React.useCallback(
    (listName: keyof TFormState, idx: number) => {
      if (!formState[listName]) {
        return;
      }
      updateFormState({
        [listName]: (formState[listName] as Array<any>).filter(
          (_, i) => i !== idx,
        ),
      } as Partial<TFormState>);
    },
    [formState, updateFormState],
  );

  const addListItem = React.useCallback(
    <TItem>(listName: keyof TFormState, item: TItem) => {
      updateFormState({
        [listName]: [...(formState[listName] || []), item],
      } as Partial<TFormState>);
    },
    [formState, updateFormState],
  );
  /**
   * Passes the event value to an updater function which returns an update
   * object to be merged into the form.
   */
  const handleInputChange = React.useCallback(
    (formUpdate: FormUpdate<TFormState>) => (
      event: React.ChangeEvent<HTMLInputElement>,
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
  update: Partial<T>,
): Array<T> {
  return array.map((item, i) => {
    if (i !== idx) {
      return item;
    }
    return {...item, ...update};
  });
}
