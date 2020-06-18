/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {FilterValue} from '../comparison_view/ComparisonViewTypes';

import InventoryQueryRenderer from '../InventoryQueryRenderer';
import React from 'react';
import SearchIcon from '@material-ui/icons/Search';
import ServicesView from './ServicesView';
import Text from '@fbcnms/ui/components/design-system/Text';
import {graphql} from 'relay-runtime';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  noResultsRoot: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: '100px',
  },
  noResultsLabel: {
    color: theme.palette.grey[600],
  },
  searchIcon: {
    color: theme.palette.grey[600],
    marginBottom: '6px',
    fontSize: '36px',
  },
}));

type Props = {
  onServiceSelected: (serviceId: string) => void,
  limit?: number,
  filters: Array<FilterValue>,
  serviceKey: number,
  onQueryReturn: number => void,
};

const serviceSearchQuery = graphql`
  query ServiceComparisonViewQueryRendererSearchQuery(
    $limit: Int
    $filters: [ServiceFilterInput!]!
  ) {
    serviceSearch(limit: $limit, filters: $filters) {
      services {
        ...ServicesView_service
      }
      count
    }
  }
`;

const ServiceComparisonViewQueryRenderer = (props: Props) => {
  const classes = useStyles();
  const {onServiceSelected, filters, limit, serviceKey, onQueryReturn} = props;

  return (
    <InventoryQueryRenderer
      query={serviceSearchQuery}
      variables={{
        limit: limit,
        filters: filters.map(f => ({
          filterType: f.name.toUpperCase(),
          operator: f.operator.toUpperCase(),
          stringValue: f.stringValue,
          propertyValue: f.propertyValue,
          idSet: f.idSet,
          stringSet: f.stringSet,
        })),
        serviceKey: serviceKey,
      }}
      render={props => {
        const {count, services} = props.serviceSearch;
        onQueryReturn(count);
        if (!services || services.length === 0) {
          return (
            <div className={classes.noResultsRoot}>
              <SearchIcon className={classes.searchIcon} />
              <Text variant="h6" className={classes.noResultsLabel}>
                No results found
              </Text>
            </div>
          );
        }
        return (
          <ServicesView
            service={services}
            onServiceSelected={onServiceSelected}
          />
        );
      }}
    />
  );
};

export default ServiceComparisonViewQueryRenderer;
