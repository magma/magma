/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {EquipmentViewQueryRendererSearchQueryResponse} from './__generated__/EquipmentViewQueryRendererSearchQuery.graphql.js';
import type {FiltersQuery} from './ComparisonViewTypes';

import ComparisonViewNoResults from './ComparisonViewNoResults';
import InventoryQueryRenderer from '../InventoryQueryRenderer';
import PowerSearchEquipmentResultsTable from './PowerSearchEquipmentResultsTable';
import React from 'react';
import useRouter from '@fbcnms/ui/hooks/useRouter';
import {InventoryAPIUrls} from '../../common/InventoryAPI';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';

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

const equipmentSearchQuery = graphql`
  query EquipmentViewQueryRendererSearchQuery(
    $limit: Int
    $filters: [EquipmentFilterInput!]!
  ) {
    equipmentSearch(limit: $limit, filters: $filters) {
      equipment {
        ...PowerSearchEquipmentResultsTable_equipment
      }
      count
    }
  }
`;

const EquipmentViewQueryRenderer = (props: Props) => {
  const classes = useStyles();
  const {limit, filters, onQueryReturn} = props;
  const {history} = useRouter();

  return (
    <InventoryQueryRenderer
      query={equipmentSearchQuery}
      variables={{
        limit: limit,
        filters: filters.map(f => ({
          filterType: f.name.toUpperCase(),
          operator: f.operator.toUpperCase(),
          stringValue: f.stringValue,
          propertyValue: f.propertyValue,
          idSet: f.idSet,
        })),
      }}
      render={(props: EquipmentViewQueryRendererSearchQueryResponse) => {
        const {count, equipment} = props.equipmentSearch;
        onQueryReturn(count);
        if (count === 0) {
          return <ComparisonViewNoResults />;
        }
        return (
          <div className={classes.searchResults}>
            <PowerSearchEquipmentResultsTable
              equipment={equipment}
              onEquipmentSelected={equipment => {
                ServerLogger.info(
                  LogEvents.EQUIPMENT_COMPARISON_VIEW_EQUIPMENT_CLICKED,
                );
                history.replace(InventoryAPIUrls.equipment(equipment.id));
              }}
              onWorkOrderSelected={workOrderId =>
                history.replace(InventoryAPIUrls.workorder(workOrderId))
              }
            />
          </div>
        );
      }}
    />
  );
};

export default EquipmentViewQueryRenderer;
