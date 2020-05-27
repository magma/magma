/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {Equipment} from '../../common/Equipment';

import Button from '@fbcnms/ui/components/design-system/Button';
import Card from '@fbcnms/ui/components/design-system/Card/Card';
import CardHeader from '@fbcnms/ui/components/design-system/Card/CardHeader';
import EquipmentTable from '../equipment/EquipmentTable';
import FormActionWithPermissions from '../../common/FormActionWithPermissions';
import React from 'react';
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_theme => ({
  cardHasNoContent: {
    marginBottom: '0px',
  },
}));

type Props = {
  className?: string,
  equipment: Array<Equipment>,
  selectedWorkOrderId: ?string,
  onEquipmentSelected: Equipment => void,
  onWorkOrderSelected: (workOrderId: string) => void,
  onAddEquipment: () => void,
};

const LocationEquipmentCard = (props: Props) => {
  const {
    equipment,
    className,
    selectedWorkOrderId,
    onEquipmentSelected,
    onWorkOrderSelected,
    onAddEquipment,
  } = props;
  const classes = useStyles();
  return (
    <Card className={className}>
      <CardHeader
        className={classNames({
          [classes.cardHasNoContent]: equipment.filter(Boolean).length === 0,
        })}
        rightContent={
          <FormActionWithPermissions
            permissions={{entity: 'equipment', action: 'create'}}>
            <Button onClick={onAddEquipment}>Add Equipment</Button>
          </FormActionWithPermissions>
        }>
        Equipment
      </CardHeader>
      <EquipmentTable
        equipment={equipment}
        selectedWorkOrderId={selectedWorkOrderId}
        onEquipmentSelected={onEquipmentSelected}
        onWorkOrderSelected={onWorkOrderSelected}
      />
    </Card>
  );
};

export default LocationEquipmentCard;
