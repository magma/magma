/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {EntityType, FiltersQuery} from './ComparisonViewTypes';

import EquipmentViewQueryRenderer from './EquipmentViewQueryRenderer';
import LinkViewQueryRenderer from './LinkViewQueryRenderer';
import LocationViewQueryRenderer from './LocationViewQueryRenderer';
import PortViewQueryRenderer from './PortViewQueryRenderer';
import React from 'react';

type Props = {
  subject: EntityType,
  filters: FiltersQuery,
  limit?: number,
  onQueryReturn: number => void,
};

const InventoryComparisonViewRouter = (props: Props) => {
  const {limit, subject, filters, onQueryReturn} = props;
  switch (subject) {
    case 'equipment':
      return (
        <EquipmentViewQueryRenderer
          limit={limit}
          onQueryReturn={onQueryReturn}
          filters={filters}
        />
      );
    case 'link':
      return (
        <LinkViewQueryRenderer
          limit={limit}
          onQueryReturn={onQueryReturn}
          filters={filters}
        />
      );
    case 'port':
      return (
        <PortViewQueryRenderer
          limit={limit}
          onQueryReturn={onQueryReturn}
          filters={filters}
        />
      );
    case 'location':
      return (
        <LocationViewQueryRenderer
          limit={limit}
          onQueryReturn={onQueryReturn}
          filters={filters}
        />
      );
    default:
      return null;
  }
};

export default InventoryComparisonViewRouter;
