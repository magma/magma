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
import PowerSearchPortsResultsTable from './PowerSearchPortsResultsTable';
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

const portSearchQuery = graphql`
  query PortViewQueryRendererSearchQuery(
    $limit: Int
    $filters: [PortFilterInput!]!
  ) {
    equipmentPorts(first: $limit, filterBy: $filters) {
      edges {
        node {
          ...PowerSearchPortsResultsTable_ports
        }
      }
      totalCount
    }
  }
`;

const PortViewQueryRenderer = (props: Props) => {
  const classes = useStyles();
  const {limit, filters, onQueryReturn} = props;

  return (
    <InventoryQueryRenderer
      query={portSearchQuery}
      variables={{
        limit: limit,
        filters: filters.map(f => ({
          filterType: f.name.toUpperCase(),
          operator: f.operator.toUpperCase(),
          stringValue: f.stringValue,
          propertyValue: f.propertyValue,
          idSet: f.idSet,
          stringSet: f.stringSet,
          boolValue: f.boolValue,
        })),
      }}
      render={props => {
        const {totalCount, edges} = props.equipmentPorts;
        onQueryReturn(totalCount);
        if (totalCount === 0) {
          return <ComparisonViewNoResults />;
        }
        const ports = edges.map(edge => edge.node);
        return (
          <div className={classes.searchResults}>
            <PowerSearchPortsResultsTable ports={ports} />
          </div>
        );
      }}
    />
  );
};

export default PortViewQueryRenderer;
