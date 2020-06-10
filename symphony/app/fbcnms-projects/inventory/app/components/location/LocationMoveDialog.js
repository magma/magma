/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import Button from '@fbcnms/ui/components/design-system/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import IconButton from '@fbcnms/ui/components/design-system/IconButton';
import LocationTypeahead from '../typeahead/LocationTypeahead';
import RadioGroup from '@fbcnms/ui/components/design-system/RadioGroup/RadioGroup';
import React, {useState} from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import {CloseIcon} from '@fbcnms/ui/components/design-system/Icons';

import fbt from 'fbt';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    position: 'relative',
  },
  closeButton: {
    position: 'absolute',
    top: '24px',
    right: '24px',
  },
  dialogTitle: {
    padding: '24px',
    paddingBottom: '16px',
  },
  dialogContent: {
    overflowY: 'hidden',
  },
  dialogActions: {
    padding: '24px',
    bottom: 0,
    display: 'flex',
    justifyContent: 'flex-end',
    width: '100%',
    backgroundColor: 'rgba(255, 255, 255, 0.9)',
    zIndex: 2,
  },
}));

type Props = $ReadOnly<{|
  locationId: string,
  locationParentId: ?string,
  open: boolean,
  onClose: () => void,
  onLocationSelected: (id: ?string) => void,
|}>;

const LocationMoveDialog = (props: Props) => {
  const {
    locationId,
    locationParentId,
    open,
    onClose,
    onLocationSelected,
  } = props;
  const [selectedLocationId, setSelectedLocationId] = useState(null);
  const [selectedValue, setSelectedValue] = useState('LOCATION');
  const classes = useStyles();
  return (
    <Dialog
      maxWidth="sm"
      open={open}
      onClose={onClose}
      fullWidth={true}
      className={classes.root}>
      <DialogTitle className={classes.dialogTitle}>
        <Text variant="h6" weight="medium">
          <fbt desc="Select location to move to dialog title">
            Choose Location
          </fbt>
        </Text>
        <IconButton
          skin="gray"
          className={classes.closeButton}
          icon={CloseIcon}
          onClick={onClose}
        />
      </DialogTitle>
      <DialogContent className={classes.dialogContent}>
        <RadioGroup
          options={[
            {
              value: 'LOCATION',
              label: (
                <LocationTypeahead
                  headline={null}
                  margin="dense"
                  selectedLocation={null}
                  onLocationSelection={location =>
                    setSelectedLocationId(location?.id ?? null)
                  }
                />
              ),
              details: '',
            },
            {
              value: 'TOP_LEVEL',
              label: fbt(
                'Make this a top-level location',
                'Caption for menu option for moving a location to be a top-level location',
              ),
              details: '',
            },
          ]}
          value={selectedValue}
          onChange={value => setSelectedValue(value)}
        />
      </DialogContent>
      <DialogActions className={classes.dialogActions}>
        <Button onClick={onClose} skin="secondaryGray">
          <fbt desc="">Cancel</fbt>
        </Button>
        <Button
          disabled={
            (selectedValue === 'LOCATION' && selectedLocationId === null) ||
            selectedLocationId === locationId ||
            selectedLocationId === locationParentId
          }
          onClick={() => {
            onLocationSelected(selectedLocationId);
          }}>
          <fbt desc="Caption for confirm button to move location">
            Move Location
          </fbt>
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default LocationMoveDialog;
