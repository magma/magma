/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {TokenizerDisplayProps} from '@fbcnms/ui/components/design-system/Token/Tokenizer';

import * as React from 'react';
import StaticNamedNodesTokenizer from './StaticNamedNodesTokenizer';
import {useLocationTypeNodes} from './LocationType';

type LocationTypesTokenizerProps = $ReadOnly<{|
  ...TokenizerDisplayProps,
  selectedLocationTypeIds?: ?$ReadOnlyArray<string>,
  onSelectedLocationTypesIdsChange?: ($ReadOnlyArray<string>) => void,
|}>;

function LocationTypesTokenizer(props: LocationTypesTokenizerProps) {
  const {
    selectedLocationTypeIds,
    onSelectedLocationTypesIdsChange,
    ...rest
  } = props;
  const locationTypes = useLocationTypeNodes();
  return (
    <StaticNamedNodesTokenizer
      allNamedNodes={locationTypes}
      selectedNodeIds={selectedLocationTypeIds}
      onSelectedNodeIdsChange={onSelectedLocationTypesIdsChange}
      {...rest}
    />
  );
}

export default LocationTypesTokenizer;
