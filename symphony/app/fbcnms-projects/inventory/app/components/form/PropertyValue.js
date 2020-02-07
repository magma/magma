/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {Property} from '../../common/Property';

import Button from '@fbcnms/ui/components/design-system/Button';
import React from 'react';
import {InventoryAPIUrls} from '../../common/InventoryAPI';
import {getPropertyValue} from '../../common/Property';
import {useRouter} from '@fbcnms/ui/hooks';

type Props = {
  property: Property,
};

const PropertyValue = ({property}: Props) => {
  const {history} = useRouter();
  const type = property.propertyType
    ? property.propertyType.type
    : property.type;

  switch (type) {
    case 'equipment':
      const equipmentValue = property.equipmentValue;
      if (equipmentValue) {
        return (
          <Button
            variant="text"
            onClick={() =>
              history.push(InventoryAPIUrls.equipment(equipmentValue.id))
            }>
            {equipmentValue.name}
          </Button>
        );
      }
    case 'location':
      const locationValue = property.locationValue;
      if (locationValue) {
        return (
          <Button
            variant="text"
            onClick={() =>
              history.push(InventoryAPIUrls.location(locationValue.id))
            }>
            {locationValue.name}
          </Button>
        );
      }
    case 'service':
      const serviceValue = property.serviceValue;
      if (serviceValue) {
        return (
          <Button
            variant="text"
            onClick={() =>
              history.push(InventoryAPIUrls.service(serviceValue.id))
            }>
            {serviceValue.name}
          </Button>
        );
      }
    default:
      return getPropertyValue(property) ?? '';
  }
};

export default PropertyValue;
