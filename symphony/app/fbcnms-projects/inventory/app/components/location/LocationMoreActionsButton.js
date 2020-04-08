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

import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import IconButton from '@fbcnms/ui/components/design-system/IconButton';
import React from 'react';
import RemoveLocationMutation from '../../mutations/RemoveLocationMutation';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {ConnectionHandler} from 'relay-runtime';
import {DeleteIcon} from '@fbcnms/ui/components/design-system/Icons';
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
      <FormAction>
        <IconButton
          skin="gray"
          onClick={this.removeLocation}
          icon={DeleteIcon}
        />
      </FormAction>
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
            // $FlowFixMe (T62907961) Relay flow types
            const parentProxy = store.get(parentLocation.id);
            // $FlowFixMe (T62907961) Relay flow types
            const currNodes = parentProxy.getLinkedRecords('children') || [];
            const withoutCurrentLocation = currNodes.filter(
              // $FlowFixMe (T62907961) Relay flow types
              child => child.getDataID() !== location.id,
            );
            // $FlowFixMe (T62907961) Relay flow types
            parentProxy.setLinkedRecords(withoutCurrentLocation, 'children');
            // $FlowFixMe (T62907961) Relay flow types
            parentProxy.setValue(
              // $FlowFixMe (T62907961) Relay flow types
              parentProxy.getValue('numChildren') - 1,
              'numChildren',
            );
          } else {
            // $FlowFixMe (T62907961) Relay flow types
            const rootQuery = store.getRoot();
            const locations = ConnectionHandler.getConnection(
              rootQuery,
              'LocationsTree_locations',
              {onlyTopLevel: true},
            );
            // $FlowFixMe (T62907961) Relay flow types
            ConnectionHandler.deleteNode(locations, location.id);
          }

          this.props.onLocationRemoved(location);
          // $FlowFixMe (T62907961) Relay flow types
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
