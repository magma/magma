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
import EquipmentTypeahead from './typeahead/EquipmentTypeahead';
import LocationTypeahead from './typeahead/LocationTypeahead';
import ServiceTypeahead from './typeahead/ServiceTypeahead';
import WorkOrderTypeahead from './typeahead/WorkOrderTypeahead';

type Props = $ReadOnly<{|
  type: string,
  value: ?NamedNode,
  onChange: (?NamedNode) => void,
  label?: ?string,
|}>;

const NodePropertyInput = (props: Props) => {
  const {type, value, onChange, label} = props;
  switch (type) {
    case 'equipment':
      return (
        <EquipmentTypeahead
          margin="dense"
          selectedEquipment={value}
          onEquipmentSelection={onChange}
          headline={label}
        />
      );
    case 'location':
      return (
        <LocationTypeahead
          margin="dense"
          selectedLocation={value}
          onLocationSelection={onChange}
          headline={label}
        />
      );
    case 'service':
      return (
        <ServiceTypeahead
          margin="dense"
          selectedService={value}
          onServiceSelection={onChange}
          headline={label}
        />
      );
    case 'work_order':
      return (
        <WorkOrderTypeahead
          margin="dense"
          selectedWorkOrder={value}
          onWorkOrderSelected={onChange}
          headline={label}
        />
      );
  }
  return null;
};

export default NodePropertyInput;
