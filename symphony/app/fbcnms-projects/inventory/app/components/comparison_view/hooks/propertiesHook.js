/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import RelayEnvironment from '../../../common/RelayEnvironment';
import {graphql} from 'relay-runtime';
import {useGraphQL} from '@fbcnms/ui/hooks';
import type {EntityType} from '../ComparisonViewTypes';

const propertiesQuery = graphql`
  query propertiesHookPossiblePropertiesQuery($entityType: PropertyEntity!) {
    possibleProperties(entityType: $entityType) {
      name
      type
      stringValue
    }
  }
`;

const usePropertyFilters = (entityType: EntityType) => {
  return useGraphQL(RelayEnvironment, propertiesQuery, {
    entityType: entityType.toString().toUpperCase(),
  });
};

export default usePropertyFilters;
