/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import React from 'react';
import WorkOrderDetailsPaneItem from './WorkOrderDetailsPaneItem';
import {createFragmentContainer, graphql} from 'react-relay';
import type {FutureState} from '../../common/WorkOrder';
import type {WorkOrderDetailsPaneEquipmentItem_equipment} from './__generated__/WorkOrderDetailsPaneEquipmentItem_equipment.graphql.js';

type Props = {
  equipment: WorkOrderDetailsPaneEquipmentItem_equipment,
  futureState: FutureState,
};

class WorkOrderDetailsPaneEquipmentItem extends React.Component<Props> {
  render() {
    const {equipment} = this.props;
    return (
      <WorkOrderDetailsPaneItem
        text={this._getEquipmentDescription(equipment, 'INSTALL')}
      />
    );
  }

  _getEquipmentDescription(
    equipment: WorkOrderDetailsPaneEquipmentItem_equipment,
    futureState: FutureState,
  ) {
    const action =
      futureState === 'INSTALL'
        ? `Add ${equipment.name} to`
        : `Remove ${equipment.name} from`;
    if (equipment.parentLocation) {
      return `${action} ${equipment.parentLocation.name}`;
    }
    if (equipment.parentPosition) {
      return `${action} ${equipment.parentPosition.parentEquipment.name} - ${equipment.parentPosition.definition.name}`;
    }
  }
}

export default createFragmentContainer(WorkOrderDetailsPaneEquipmentItem, {
  equipment: graphql`
    fragment WorkOrderDetailsPaneEquipmentItem_equipment on Equipment {
      id
      name
      equipmentType {
        id
        name
      }
      parentLocation {
        id
        name
        locationType {
          id
          name
        }
      }
      parentPosition {
        id
        definition {
          name
          visibleLabel
        }
        parentEquipment {
          id
          name
        }
      }
    }
  `,
});
