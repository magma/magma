/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {Equipment, EquipmentPosition} from '../../common/Equipment';
import type {EquipmentType} from '../../common/EquipmentType';

import Card from '@fbcnms/ui/components/design-system/Card/Card';
import DynamicPropertiesGrid from '../DynamicPropertiesGrid';
import EquipmentPositionsGrid from './EquipmentPositionsGrid';
import React from 'react';
import {makeStyles} from '@material-ui/styles';

type Props = $ReadOnly<{|
  equipment: Equipment,
  workOrderId: ?string,
  onAttachingEquipmentToPosition: (
    equipmentType: EquipmentType,
    position: EquipmentPosition,
  ) => void,
  onEquipmentClicked: (equipmentId: string) => void,
  onWorkOrderSelected: (workOrderId: string) => void,
|}>;

const useStyles = makeStyles(() => ({
  card: {
    '&>:not(:first-child)': {
      marginTop: '16px',
    },
  },
}));

function EquipmentDetails(props: Props) {
  const {
    equipment,
    workOrderId,
    onAttachingEquipmentToPosition,
    onEquipmentClicked,
    onWorkOrderSelected,
  } = props;
  const propTypes = equipment.equipmentType.propertyTypes;
  const classes = useStyles();

  return (
    <Card contentClassName={classes.card}>
      <DynamicPropertiesGrid
        properties={equipment.properties}
        propertyTypes={propTypes}
      />
      <EquipmentPositionsGrid
        positionDefinitions={equipment.equipmentType.positionDefinitions}
        positions={equipment.positions}
        workOrderId={workOrderId}
        onAttachingEquipmentToPosition={onAttachingEquipmentToPosition}
        onEquipmentPositionClicked={onEquipmentClicked}
        equipment={equipment}
        onWorkOrderSelected={onWorkOrderSelected}
      />
    </Card>
  );
}

export default EquipmentDetails;
