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
import axios from 'axios';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

import type {GenericRule, RuleInterfaceMap} from './RuleInterface';

/**
 * Loads alert rules for each rule type. Rules are loaded in parallel and if one
 * rule loader fails or exceeds the load timeout, it will be cancelled and a
 * snackbar will be enqueued.
 */
export function useLoadRules<TRuleUnion>({
  ruleMap,
  lastRefreshTime,
}: {
  ruleMap: RuleInterfaceMap<TRuleUnion>,
  lastRefreshTime: string,
}): {rules: Array<GenericRule<TRuleUnion>>, isLoading: boolean} {
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [isLoading, setIsLoading] = React.useState(true);
  const [rules, setRules] = React.useState<Array<GenericRule<TRuleUnion>>>([]);

  React.useEffect(() => {
    const promises = Object.keys(ruleMap || {}).map((ruleType: string) => {
      const cancelSource = axios.CancelToken.source();
      const request = {
        // for magma api
        networkId: match.params.networkId,
        cancelToken: cancelSource.token,
      };

      const ruleInterface = ruleMap[ruleType];

      return new Promise(resolve => {
        ruleInterface
          .getRules(request)
          .then(response => resolve(response))
          .catch(error => {
            console.error(error);
            enqueueSnackbar(
              `An error occurred while loading ${ruleInterface.friendlyName} rules`,
              {
                variant: 'error',
              },
            );
            resolve([]);
          });
      });
    });

    Promise.all(promises).then(results => {
      const allResults = [].concat.apply([], results);
      setRules(allResults);
      setIsLoading(false);
    });
  }, [enqueueSnackbar, lastRefreshTime, match.params.networkId, ruleMap]);

  return {
    rules,
    isLoading,
  };
}

type InputChangeFunc<TFormState> = (
  formUpdate: FormUpdate<TFormState>,
) => (event: SyntheticInputEvent<HTMLElement>) => void;
type FormUpdate<TFormState> = (val: string) => $Shape<TFormState>;

export function useForm<TFormState: {}>({
  initialState,
  onFormUpdated,
}: {
  initialState: TFormState,
  onFormUpdated?: (state: TFormState) => void,
}): {|
  formState: TFormState,
  updateFormState: (update: $Shape<TFormState>) => TFormState,
  handleInputChange: InputChangeFunc<TFormState>,
  updateListItem: (
    listName: $Keys<TFormState>,
    idx: number,
    update: $Shape<TFormState> | TFormState,
  ) => void,
  addListItem: (listName: $Keys<TFormState>, item: {}) => void,
  removeListItem: (listName: $Keys<TFormState>, idx: number) => void,
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
      return nextState;
    },
    [formState, setFormState],
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
      const updated = updateFormState(formUpdate(value));
      if (typeof onFormUpdated === 'function') {
        onFormUpdated(updated);
      }
    },
    [onFormUpdated, updateFormState],
  );

  return {
    formState,
    updateFormState,
    handleInputChange,
    updateListItem,
    addListItem,
    removeListItem,
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
