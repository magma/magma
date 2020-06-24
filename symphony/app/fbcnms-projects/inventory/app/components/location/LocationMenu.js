/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {LocationMenu_location} from './__generated__/LocationMenu_location.graphql';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithSnackbarProps} from 'notistack';

import LocationMoveDialog from './LocationMoveDialog';
import OptionsPopoverButton from '../OptionsPopoverButton';
import React, {useState} from 'react';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import fbt from 'fbt';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {ThreeDotsHorizontalIcon} from '@fbcnms/ui/components/design-system/Icons';
import {createFragmentContainer, graphql} from 'react-relay';
import {makeStyles} from '@material-ui/styles';
import {moveLocation} from '../../mutations/MoveLocationMutation';
import {removeLocation} from '../../mutations/RemoveLocationMutation';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';

type Props = $ReadOnly<
  {|
    location: LocationMenu_location,
    onLocationMoved: (location: LocationMenu_location) => void,
    onLocationRemoved: (location: LocationMenu_location) => void,
    popoverMenuClassName?: ?string,
    onVisibilityChange?: (isVisible: boolean) => void,
  |} & WithSnackbarProps &
    WithAlert,
>;

const useStyles = makeStyles(() => ({
  popoverMenu: {
    minWidth: '150px',
  },
}));

const LocationMenu = (props: Props) => {
  const {
    location,
    onLocationMoved,
    onLocationRemoved,
    onVisibilityChange,
    confirm,
  } = props;
  const [moveDialogOpen, setMoveDialogOpen] = useState(false);
  const enqueueSnackbar = useEnqueueSnackbar();

  const handleMove = (targetLocationId: ?string) => {
    moveLocation(location.id, location.parentLocation?.id, targetLocationId)
      .then(() => onLocationMoved(location))
      .catch((errorMessage: string) =>
        enqueueSnackbar(errorMessage, {
          children: (key: string) => (
            <SnackbarItem id={key} message={errorMessage} variant="error" />
          ),
        }),
      );
  };

  const handleRemove = () => {
    confirm(
      fbt(
        `Are you sure you want to delete ${fbt.param('name', location.name)}`,
        '',
      ).toString(),
    )
      .then(confirmed => {
        if (!confirmed) {
          return;
        }
        removeLocation(location.id, location.parentLocation?.id);
      })
      .then(() => onLocationRemoved(location))
      .catch((errorMessage: string) =>
        enqueueSnackbar(errorMessage, {
          children: (key: string) => (
            <SnackbarItem id={key} message={errorMessage} variant="error" />
          ),
        }),
      );
  };

  const menuOptions = [
    {
      onClick: () => setMoveDialogOpen(true),
      caption: fbt(
        'Move Location',
        'Caption for menu option for moving a location to a different location',
      ),
      permissions: {
        entity: 'location',
        action: 'update',
        locationTypeId: location.locationType.id,
        hideOnMissingPermissions: true,
      },
    },
  ];

  if (
    location.children.filter(Boolean).length === 0 &&
    location.equipments.filter(Boolean).length === 0 &&
    location.images.filter(Boolean).length === 0 &&
    location.files.filter(Boolean).length === 0 &&
    location.surveys.filter(Boolean).length === 0
  ) {
    menuOptions.push({
      onClick: handleRemove,
      caption: fbt(
        'Delete Location',
        'Caption for menu option for deleting a location',
      ),
      permissions: {
        entity: 'location',
        action: 'delete',
        hideOnMissingPermissions: true,
        locationTypeId: location.locationType.id,
      },
    });
  }

  const classes = useStyles();

  return (
    <>
      <OptionsPopoverButton
        options={menuOptions}
        popoverMenuClassName={classes.popoverMenu}
        onVisibilityChange={onVisibilityChange}
        menuIcon={<ThreeDotsHorizontalIcon color="gray" />}
      />
      <LocationMoveDialog
        locationId={location.id}
        locationParentId={location.parentLocation?.id}
        open={moveDialogOpen}
        onClose={() => setMoveDialogOpen(false)}
        onLocationSelected={locationId => {
          setMoveDialogOpen(false);
          handleMove(locationId);
        }}
      />
    </>
  );
};

export default withAlert(
  createFragmentContainer(LocationMenu, {
    location: graphql`
      fragment LocationMenu_location on Location {
        id
        name
        locationType {
          id
        }
        parentLocation {
          id
        }
        children {
          id
        }
        equipments {
          id
        }
        images {
          id
        }
        files {
          id
        }
        surveys {
          id
        }
      }
    `,
  }),
);
