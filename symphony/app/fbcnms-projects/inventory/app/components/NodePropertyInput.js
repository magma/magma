/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {NamedNode} from '../common/EntUtils';

import * as React from 'react';
import EquipmentTypeahead from './typeahead/EquipmentTypeahead';
import LocationTypeahead from './typeahead/LocationTypeahead';
import ServiceTypeahead from './typeahead/ServiceTypeahead';
import UserTypeahead from './typeahead/UserTypeahead';
import WorkOrderTypeahead from './typeahead/WorkOrderTypeahead';

type Props = $ReadOnly<{|
  type: string,
  value: ?NamedNode,
  onChange: (?NamedNode) => void,
  label?: ?string,
|}>;

const NodePropertyInput = (props: Props) => {
  const {type, value, onChange, label} = props;
  const basicValue = value != null ? {id: value.id, name: value.name} : null;
  switch (type) {
    case 'equipment':
      return (
        <EquipmentTypeahead
          margin="dense"
          selectedEquipment={basicValue}
          onEquipmentSelection={onChange}
          headline={label}
        />
      );
    case 'location':
      return (
        <LocationTypeahead
          margin="dense"
          selectedLocation={basicValue}
          onLocationSelection={onChange}
          headline={label}
        />
      );
    case 'service':
      return (
        <ServiceTypeahead
          margin="dense"
          selectedService={basicValue}
          onServiceSelection={onChange}
          headline={label}
        />
      );
    case 'work_order':
      return (
        <WorkOrderTypeahead
          margin="dense"
          selectedWorkOrder={basicValue}
          onWorkOrderSelected={onChange}
          headline={label}
        />
      );
    case 'user':
      return (
        <UserTypeahead
          margin="dense"
          selectedUser={value ? {id: value.id, email: value.name ?? ''} : null}
          onUserSelection={newUser =>
            onChange(newUser ? {id: newUser.id, name: newUser.email} : null)
          }
          headline={label}
        />
      );
  }
  return null;
};

export default NodePropertyInput;
