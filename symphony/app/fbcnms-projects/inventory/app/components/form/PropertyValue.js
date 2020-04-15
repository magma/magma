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
import {useHistory} from 'react-router';

type Props = {
  property: Property,
};

const PropertyValue = ({property}: Props) => {
  const history = useHistory();
  const propType = property.propertyType ? property.propertyType : property;

  switch (propType.type) {
    case 'node':
      const nodeValue = property.nodeValue;
      if (nodeValue) {
        switch (propType.nodeType) {
          case 'equipment':
            return (
              <Button
                variant="text"
                onClick={() =>
                  history.push(InventoryAPIUrls.equipment(nodeValue.id))
                }>
                {nodeValue.name}
              </Button>
            );
          case 'location':
            return (
              <Button
                variant="text"
                onClick={() =>
                  history.push(InventoryAPIUrls.location(nodeValue.id))
                }>
                {nodeValue.name}
              </Button>
            );
          case 'service':
            return (
              <Button
                variant="text"
                onClick={() =>
                  history.push(InventoryAPIUrls.service(nodeValue.id))
                }>
                {nodeValue.name}
              </Button>
            );
        }
      }

    default:
      return getPropertyValue(property) ?? '';
  }
};

export default PropertyValue;
