/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {LocationMoreActionsButton_location} from './__generated__/LocationMoreActionsButton_location.graphql';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {
  RemoveLocationMutationResponse,
  RemoveLocationMutationVariables,
} from '../../mutations/__generated__/RemoveLocationMutation.graphql';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithSnackbarProps} from 'notistack';

import MoreActionsButton from '@fbcnms/ui/components/MoreActionsButton';
import React from 'react';
import RemoveLocationMutation from '../../mutations/RemoveLocationMutation';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {ConnectionHandler} from 'relay-runtime';
import {createFragmentContainer, graphql} from 'react-relay';
import {withSnackbar} from 'notistack';

type Props = {
  location: LocationMoreActionsButton_location,
  onLocationRemoved: (
    removedLocation: LocationMoreActionsButton_location,
  ) => void,
} & WithAlert &
  WithSnackbarProps;

class LocationMoreActionsButton extends React.Component<Props> {
  render() {
    return (
      <MoreActionsButton
        variant="primary"
        items={[
          {
            name: 'Delete location',
            onClick: this.removeLocation,
          },
        ]}
      />
    );
  }

  removeLocation = () => {
    const {location} = this.props;
    if (
      location.children.filter(Boolean).length > 0 ||
      location.equipments.filter(Boolean).length > 0 ||
      location.images.filter(Boolean).length > 0 ||
      location.files.filter(Boolean).length > 0 ||
      location.surveys.filter(Boolean).length > 0
    ) {
      this.props.alert(
        'Cannot delete populated location (e.g. has equipment or files)',
      );
      return;
    }

    this.props
      .confirm({
        message: 'Are you sure you want to delete this location?',
        confirmLabel: 'Delete',
      })
      .then(confirmed => {
        if (!confirmed) {
          return;
        }

        const variables: RemoveLocationMutationVariables = {
          id: nullthrows(location.id),
        };

        const updater = store => {
          const {parentLocation} = location;
          if (!!parentLocation) {
            const parentProxy = store.get(parentLocation.id);
            const currNodes = parentProxy.getLinkedRecords('children') || [];
            const withoutCurrentLocation = currNodes.filter(
              child => child.getDataID() !== location.id,
            );
            parentProxy.setLinkedRecords(withoutCurrentLocation, 'children');
            parentProxy.setValue(
              parentProxy.getValue('numChildren') - 1,
              'numChildren',
            );
          } else {
            const rootQuery = store.getRoot();
            const locations = ConnectionHandler.getConnection(
              rootQuery,
              'LocationsTree_locations',
              {onlyTopLevel: true},
            );
            ConnectionHandler.deleteNode(locations, location.id);
          }

          this.props.onLocationRemoved(location);
          store.delete(location.id);
        };

        const callbacks: MutationCallbacks<RemoveLocationMutationResponse> = {
          onCompleted: (response, errors) => {
            if (errors && errors[0]) {
              this.props.enqueueSnackbar(errors[0].message, {
                children: key => (
                  <SnackbarItem
                    id={key}
                    message={errors[0].message}
                    variant="error"
                  />
                ),
              });
            }
          },
          onError: (_error: Error) => {
            this.props.alert('Failed removing location');
          },
        };

        RemoveLocationMutation(variables, callbacks, updater);
      });
  };
}

export default withAlert(
  withSnackbar(
    createFragmentContainer(LocationMoreActionsButton, {
      location: graphql`
        fragment LocationMoreActionsButton_location on Location {
          id
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
  ),
);
