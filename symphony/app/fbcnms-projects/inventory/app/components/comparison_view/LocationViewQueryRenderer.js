/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {FiltersQuery} from './ComparisonViewTypes';

import ComparisonViewNoResults from './ComparisonViewNoResults';
import InventoryQueryRenderer from '../InventoryQueryRenderer';
import PowerSearchLocationsResultsTable from './PowerSearchLocationsResultsTable';
import React from 'react';

import {graphql} from 'relay-runtime';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_theme => ({
  searchResults: {
    flexGrow: 1,
  },
}));

type Props = {
  filters: FiltersQuery,
  limit?: number,
  onQueryReturn: number => void,
};

const locationSearchQuery = graphql`
  query LocationViewQueryRendererSearchQuery(
    $limit: Int
    $filters: [LocationFilterInput!]!
  ) {
    locationSearch(limit: $limit, filters: $filters) {
      locations {
        ...PowerSearchLocationsResultsTable_locations
      }
      count
    }
  }
`;

const LocationViewQueryRenderer = (props: Props) => {
  const classes = useStyles();
  const {limit, filters, onQueryReturn} = props;

  return (
    <InventoryQueryRenderer
      query={locationSearchQuery}
      variables={{
        limit: limit,
        filters: filters.map(f => ({
          filterType: f.name.toUpperCase(),
          operator: f.operator.toUpperCase(),
          stringValue: f.stringValue,
          propertyValue: f.propertyValue,
          boolValue: f.boolValue,
          idSet: f.idSet,
          stringSet: f.stringSet,
        })),
      }}
      render={props => {
        const {count, locations} = props.locationSearch;
        onQueryReturn(count);
        if (count === 0) {
          return <ComparisonViewNoResults />;
        }
        return (
          <div className={classes.searchResults}>
            <PowerSearchLocationsResultsTable locations={locations} />
          </div>
        );
      }}
    />
  );
};

export default LocationViewQueryRenderer;
