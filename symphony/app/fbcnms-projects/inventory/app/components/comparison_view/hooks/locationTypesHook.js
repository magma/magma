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
import {buildLocationTypeFilterConfigs, getLocationTypes} from '../FilterUtils';
import {graphql} from 'relay-runtime';
import {useGraphQL} from '@fbcnms/ui/hooks';
import {useMemo} from 'react';

const locationTypesQuery = graphql`
  query locationTypesHookLocationTypesQuery {
    locationTypes(first: 20) {
      edges {
        node {
          id
          name
        }
      }
    }
  }
`;

const useLocationTypes = () => {
  const locationTypesResponse = useGraphQL(
    // $FlowFixMe (T62907961) Relay flow types
    RelayEnvironment,
    locationTypesQuery,
    {},
  );

  return useMemo(() => {
    if (locationTypesResponse.response === null) {
      return null;
    }
    return buildLocationTypeFilterConfigs(
      getLocationTypes(locationTypesResponse.response),
    );
  }, [locationTypesResponse.response]);
};

export default useLocationTypes;
