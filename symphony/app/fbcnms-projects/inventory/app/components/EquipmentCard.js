/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {EquipmentPosition} from '../common/Equipment';
import type {EquipmentType} from '../common/EquipmentType';

import EquipmentAddEditCard from './equipment/EquipmentAddEditCard';
import EquipmentPropertiesCard from './equipment/EquipmentPropertiesCard';
import React from 'react';
import nullthrows from '@fbcnms/util/nullthrows';
import withInventoryErrorBoundary from '../common/withInventoryErrorBoundary';

type Props = {
  mode: 'add' | 'edit' | 'show',
  onSave: () => void,
  onEdit: () => void,
  onCancel: () => void,
  selectedEquipmentId: ?string,
  selectedEquipmentPosition: ?EquipmentPosition,
  selectedLocationId: ?string,
  selectedEquipmentType: ?EquipmentType,
  selectedWorkOrderId: ?string,
  onAttachingEquipmentToPosition: (
    equipmentType: EquipmentType,
    position: EquipmentPosition,
  ) => void,
  onEquipmentClicked: (equipmentId: string) => void,
  onParentLocationClicked: (locationId: string) => void,
  onWorkOrderSelected: (workOrderId: string) => void,
};

class EquipmentCard extends React.Component<Props> {
  render() {
    switch (this.props.mode) {
      case 'add':
        return (
          <EquipmentAddEditCard
            key="new_equipment_type"
            locationId={this.props.selectedLocationId}
            equipmentPosition={this.props.selectedEquipmentPosition}
            workOrderId={this.props.selectedWorkOrderId}
            type={this.props.selectedEquipmentType}
            onCancel={this.props.onCancel}
            onSave={this.props.onSave}
          />
        );
      case 'edit':
        return (
          <EquipmentAddEditCard
            key={'new_equipment_type_' + String(this.props.selectedEquipmentId)}
            locationId={this.props.selectedLocationId}
            editingEquipmentId={this.props.selectedEquipmentId}
            equipmentPosition={this.props.selectedEquipmentPosition}
            workOrderId={this.props.selectedWorkOrderId}
            type={this.props.selectedEquipmentType}
            onCancel={this.props.onCancel}
            onSave={this.props.onSave}
          />
        );
      case 'show':
        return (
          <EquipmentPropertiesCard
            equipmentId={nullthrows(this.props.selectedEquipmentId)}
            onAttachingEquipmentToPosition={
              this.props.onAttachingEquipmentToPosition
            }
            onEquipmentClicked={this.props.onEquipmentClicked}
            onParentLocationClicked={this.props.onParentLocationClicked}
            workOrderId={this.props.selectedWorkOrderId}
            onEdit={this.props.onEdit}
            onWorkOrderSelected={this.props.onWorkOrderSelected}
          />
        );
      default:
        return null;
    }
  }
}

export default withInventoryErrorBoundary(EquipmentCard);
