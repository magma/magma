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
import type {Location} from '../../common/Location.js';

import Card from '@fbcnms/ui/components/design-system/Card/Card';
import CardHeader from '@fbcnms/ui/components/design-system/Card/CardHeader';
import DynamicPropertiesGrid from '../DynamicPropertiesGrid';
import LocationDetailsCard from './LocationDetailsCard';
import LocationEquipmentCard from './LocationEquipmentCard';
import React from 'react';
import {makeStyles} from '@material-ui/styles';

type Props = {
  location: Location,
  selectedWorkOrderId: ?string,
  onEquipmentSelected: Equipment => void,
  onWorkOrderSelected: (workOrderId: string) => void,
  onAddEquipment: () => void,
};

const useStyles = makeStyles(_theme => ({
  card: {
    marginBottom: '16px',
  },
}));

const LocationDetailsTab = (props: Props) => {
  const classes = useStyles();
  const {
    location,
    selectedWorkOrderId,
    onEquipmentSelected,
    onWorkOrderSelected,
    onAddEquipment,
  } = props;

  const propTypes = location.locationType.propertyTypes;

  return (
    <div>
      <LocationDetailsCard className={classes.card} location={location} />
      <Card className={classes.card}>
        <CardHeader>Properties</CardHeader>
        <DynamicPropertiesGrid
          hideTitle={true}
          properties={location.properties}
          propertyTypes={propTypes}
        />
      </Card>
      <LocationEquipmentCard
        className={classes.card}
        equipment={location.equipments}
        selectedWorkOrderId={selectedWorkOrderId}
        onEquipmentSelected={onEquipmentSelected}
        onWorkOrderSelected={onWorkOrderSelected}
        onAddEquipment={onAddEquipment}
      />
    </div>
  );
};

export default LocationDetailsTab;
