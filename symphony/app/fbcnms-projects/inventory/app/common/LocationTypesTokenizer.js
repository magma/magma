/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  Entries,
  TokenizerDisplayProps,
} from '@fbcnms/ui/components/design-system/Token/Tokenizer';
import type {LocationTypeNode} from './LocationType';

import * as React from 'react';
import Tokenizer from '@fbcnms/ui/components/design-system/Token/Tokenizer';
import withSuspense from './withSuspense';
import {useCallback, useEffect, useMemo, useState} from 'react';
import {useLocationTypeNodes} from './LocationType';

type LocationTypesTokenizerProps = $ReadOnly<{|
  ...TokenizerDisplayProps,
  selectedLocationTypeIds?: ?$ReadOnlyArray<string>,
  onSelectedLocationTypesIdsChange?: ($ReadOnlyArray<string>) => void,
|}>;

function LocationTypesTokenizer(props: LocationTypesTokenizerProps) {
  const {
    onSelectedLocationTypesIdsChange,
    selectedLocationTypeIds,
    disabled: disabledProp,
    ...tokenizerDisplayProps
  } = props;

  const wrapAsEntries = useCallback(
    items =>
      (items || []).map(item => ({
        ...item,
        key: item.id,
        label: item.name,
      })),
    [],
  );

  const locationTypes = useLocationTypeNodes();
  const locationTypeEntries = useMemo(() => wrapAsEntries(locationTypes), [
    locationTypes,
    wrapAsEntries,
  ]);

  const [selectedLocationTypes, setSelectedLocationTypes] = useState<
    Entries<LocationTypeNode>,
  >([]);
  useEffect(() => {
    const ids = selectedLocationTypeIds || [];
    const selectionWasChanged =
      ids.length !== selectedLocationTypes.length ||
      selectedLocationTypes.find(lt => !ids.includes(lt.id));
    if (!selectionWasChanged) {
      return;
    }
    const foundLocationTypes = ids
      .map(id => locationTypes.find(lt => lt.id === id))
      .filter(Boolean);
    setSelectedLocationTypes(wrapAsEntries(foundLocationTypes));
  }, [
    locationTypes,
    selectedLocationTypeIds,
    selectedLocationTypes,
    wrapAsEntries,
  ]);

  const updateSelectedLocationTypes = useCallback(
    (newEntries: Entries<LocationTypeNode>) => {
      setSelectedLocationTypes(newEntries);
      if (!onSelectedLocationTypesIdsChange) {
        return;
      }
      onSelectedLocationTypesIdsChange(newEntries.map(lte => lte.id));
    },
    [onSelectedLocationTypesIdsChange],
  );

  const [queryString, setQueryString] = useState('');

  const disabled = disabledProp || onSelectedLocationTypesIdsChange == null;
  return (
    <Tokenizer
      {...tokenizerDisplayProps}
      disabled={disabled}
      tokens={selectedLocationTypes}
      onTokensChange={updateSelectedLocationTypes}
      queryString={queryString}
      onQueryStringChange={setQueryString}
      dataSource={{
        fetchNetwork: () => Promise.resolve(locationTypeEntries),
      }}
    />
  );
}

export default withSuspense(LocationTypesTokenizer);
