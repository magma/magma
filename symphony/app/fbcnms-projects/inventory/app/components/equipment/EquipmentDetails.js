/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Equipment, EquipmentPosition} from '../../common/Equipment';
import type {EquipmentType} from '../../common/EquipmentType';
import type {WithStyles} from '@material-ui/core';

import DynamicPropertiesGrid from '../DynamicPropertiesGrid';
import EquipmentPositionsGrid from './EquipmentPositionsGrid';
import React from 'react';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  equipment: Equipment,
  workOrderId: ?string,
  onAttachingEquipmentToPosition: (
    equipmentType: EquipmentType,
    position: EquipmentPosition,
  ) => void,
  onEquipmentClicked: (equipmentId: string) => void,
  onWorkOrderSelected: (workOrderId: string) => void,
} & WithStyles<typeof styles>;

const styles = theme => ({
  field: {
    marginLeft: 5,
  },
  positionsGrid: {
    marginTop: theme.spacing(3),
  },
});

class EquipmentDetails extends React.Component<Props> {
  render() {
    const {classes, equipment, workOrderId} = this.props;
    const propTypes = equipment.equipmentType.propertyTypes;
    return (
      <div className={classes.cardDetails}>
        <DynamicPropertiesGrid
          properties={equipment.properties}
          propertyTypes={propTypes}
        />
        <div className={classes.positionsGrid}>
          <EquipmentPositionsGrid
            positionDefinitions={equipment.equipmentType.positionDefinitions}
            positions={equipment.positions}
            workOrderId={workOrderId}
            onAttachingEquipmentToPosition={
              this.props.onAttachingEquipmentToPosition
            }
            onEquipmentPositionClicked={this.props.onEquipmentClicked}
            equipment={equipment}
            onWorkOrderSelected={this.props.onWorkOrderSelected}
          />
        </div>
      </div>
    );
  }
}

export default withStyles(styles)(EquipmentDetails);
