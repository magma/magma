/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {NamedNode} from '../common/EntUtils';

import * as React from 'react';
import Button from '@fbcnms/ui/components/design-system/Button';
import {InventoryAPIUrls} from '../common/InventoryAPI';
import {useHistory} from 'react-router';

type Props = $ReadOnly<{|
  type: string,
  value: ?NamedNode,
|}>;

const NodePropertyValue = (props: Props) => {
  const {type, value} = props;
  const history = useHistory();
  if (value) {
    const onNodeClicked = () => {
      switch (type) {
        case 'equipment':
          history.push(InventoryAPIUrls.equipment(value.id));
          break;
        case 'location':
          history.push(InventoryAPIUrls.location(value.id));
          break;
        case 'service':
          history.push(InventoryAPIUrls.service(value.id));
          break;
        case 'work_order':
          history.push(InventoryAPIUrls.workorder(value.id));
          break;
      }
    };

    switch (type) {
      case 'equipment':
      case 'location':
      case 'service':
      case 'work_order':
        return (
          <Button variant="text" onClick={onNodeClicked}>
            {value.name}
          </Button>
        );
      default:
        return value.name;
    }
  }
  return null;
};

export default NodePropertyValue;
