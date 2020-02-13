/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {EntityType, FiltersQuery} from './ComparisonViewTypes';

import EquipmentPowerSearchBar from './EquipmentPowerSearchBar';
import LinksPowerSearchBar from './LinksPowerSearchBar';
import LocationsPowerSearchBar from './LocationsPowerSearchBar';
import PortsPowerSearchBar from './PortsPowerSearchBar';
import React from 'react';

type Props = {
  subject: EntityType,
  filters: FiltersQuery,
  onFiltersChanged: FiltersQuery => void,
  footer?: ?string,
};

const PwerSearchBarRouter = (props: Props) => {
  const {subject, filters, onFiltersChanged, footer} = props;
  switch (subject) {
    case 'equipment':
      return (
        <EquipmentPowerSearchBar
          onFiltersChanged={onFiltersChanged}
          filters={filters}
          footer={footer}
        />
      );
    case 'link':
      return (
        <LinksPowerSearchBar
          onFiltersChanged={onFiltersChanged}
          filters={filters}
          footer={footer}
        />
      );
    case 'port':
      return (
        <PortsPowerSearchBar
          onFiltersChanged={onFiltersChanged}
          filters={filters}
          footer={footer}
        />
      );
    case 'location':
      return (
        <LocationsPowerSearchBar
          onFiltersChanged={onFiltersChanged}
          filters={filters}
          footer={footer}
        />
      );
    default:
      return null;
  }
};

export default PwerSearchBarRouter;
