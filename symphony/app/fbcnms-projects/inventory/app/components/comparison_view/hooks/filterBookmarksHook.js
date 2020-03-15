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
  FilterEntity,
  filterBookmarksHookReportFiltersQueryResponse,
} from './__generated__/filterBookmarksHookReportFiltersQuery.graphql';
import type {SavedSearchConfig} from '../ComparisonViewTypes';

import RelayEnvironment from '../../../common/RelayEnvironment';
import nullthrows from '@fbcnms/util/nullthrows';
import shortid from 'shortid';
import {graphql} from 'relay-runtime';
import {toOperator} from '../FilterUtils';
import {useGraphQL} from '@fbcnms/ui/hooks';
import {useMemo} from 'react';

const reportFilterQuery = graphql`
  query filterBookmarksHookReportFiltersQuery($entity: FilterEntity!) {
    reportFilters(entity: $entity) {
      id
      name
      entity
      filters {
        operator
        key
        filterType
        stringValue
        idSet
        stringSet
        boolValue
        propertyValue {
          index
          name
          type
          stringValue
          intValue
          booleanValue
          floatValue
          latitudeValue
          longitudeValue
          rangeFromValue
          rangeToValue
          isDeleted
        }
      }
    }
  }
`;

const useFilterBookmarks = (entity: FilterEntity): Array<SavedSearchConfig> => {
  const filterBookmarksResponse = useGraphQL(
    // $FlowFixMe (T62907961) Relay flow types
    RelayEnvironment,
    reportFilterQuery,
    {entity: entity},
  );

  return useMemo(() => {
    if (filterBookmarksResponse.response === null) {
      return [];
    }
    return convertSavedSearchToSavedSearchConfig(
      filterBookmarksResponse.response,
    );
  }, [filterBookmarksResponse.response]);
};

const convertSavedSearchToSavedSearchConfig = (
  response: filterBookmarksHookReportFiltersQueryResponse,
): Array<SavedSearchConfig> => {
  return response.reportFilters.map(f => {
    const fullFilters = f.filters.map(x => {
      return {
        ...x,
        id: shortid.generate(),
        idSet: x.idSet?.slice() ?? [],
        stringSet: x.stringSet?.slice() ?? [],
        operator: toOperator(x.operator),
        name: x.filterType,
        propertyValue: x.propertyValue
          ? {
              ...x.propertyValue,
              id: f.id + Math.floor(Math.random() * 1000),
              index: x.propertyValue?.index ?? 0,
              name: nullthrows(x.propertyValue?.name),
              type: nullthrows(x.propertyValue?.type),
            }
          : null,
      };
    });
    return {
      id: f.id,
      key: 'saved_search_' + f.id,
      label: f.name,
      entity: f.entity,
      filters: fullFilters,
    };
  });
};

export default useFilterBookmarks;
