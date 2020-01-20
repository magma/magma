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

import Button from '@fbcnms/ui/components/design-system/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import EquipmentTypesList from '../EquipmentTypesList';
import LocationEquipments from '../location/LocationEquipments';
import MoveEquipmentToPositionMutation from '../../mutations/MoveEquipmentToPositionMutation';
import React from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import {createFragmentContainer, graphql} from 'react-relay';
import {last} from 'lodash';

import nullthrows from '@fbcnms/util/nullthrows';
import {makeStyles} from '@material-ui/styles';
import {useState} from 'react';

type Props = {
  open: boolean,
  onClose: () => void,
  onEquipmentTypeSelected: (equipmentType: EquipmentType) => void,
  parentEquipment: Equipment,
  position: EquipmentPosition,
};

const useStyles = makeStyles(theme => ({
  tabs: {
    borderBottom: `1px ${theme.palette.grey[200]} solid`,
  },
}));

const AddToEquipmentDialog = (props: Props) => {
  const [tab, setTab] = useState('new');
  const [selectedEquipmentType, setSelectedEquipmentType] = useState(null);
  const [selectedEquipment, setSelectedEquipment] = useState(null);
  const classes = useStyles();

  const locations = props.parentEquipment.locationHierarchy;
  const parentLocation = nullthrows(last(locations));
  return (
    <Dialog maxWidth="sm" open={props.open} onClose={props.onClose}>
      <DialogContent>
        <Tabs
          className={classes.tabs}
          value={tab}
          onChange={(_, newTab) => setTab(newTab)}
          indicatorColor="primary"
          textColor="primary">
          <Tab label="New Equipment" value="new" />
          <Tab label="Existing Equipment" value="existing" />
        </Tabs>
        {tab === 'new' && (
          <EquipmentTypesList
            onSelect={type => setSelectedEquipmentType(type)}
          />
        )}
        {tab === 'existing' && (
          <LocationEquipments
            locationId={parentLocation.id}
            onSelect={equipment => setSelectedEquipment(equipment)}
          />
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} skin="regular">
          Cancel
        </Button>
        <Button
          disabled={
            (tab === 'new' && selectedEquipmentType === null) ||
            (tab === 'existing' && selectedEquipment === null)
          }
          onClick={() => {
            if (tab === 'new') {
              props.onEquipmentTypeSelected(nullthrows(selectedEquipmentType));
            } else {
              MoveEquipmentToPositionMutation(
                {
                  parent_equipment_id: props.parentEquipment.id,
                  position_definition_id: props.position.definition.id,
                  equipment_id: nullthrows(selectedEquipment).id,
                },
                {
                  onCompleted: props.onClose,
                  onError: () => {},
                },
              );
            }
          }}
          color="primary">
          OK
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default createFragmentContainer(AddToEquipmentDialog, {
  parentEquipment: graphql`
    fragment AddToEquipmentDialog_parentEquipment on Equipment {
      id
      locationHierarchy {
        id
      }
    }
  `,
});
