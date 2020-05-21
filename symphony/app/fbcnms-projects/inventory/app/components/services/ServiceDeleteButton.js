/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {
  RemoveServiceMutationResponse,
  RemoveServiceMutationVariables,
} from '../../mutations/__generated__/RemoveServiceMutation.graphql';
import type {Service} from '../../common/Service';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithSnackbarProps} from 'notistack';

import DeleteIcon from '@fbcnms/ui/components/design-system/Icons/Actions/DeleteIcon';
import FormActionWithPermissions from '../../common/FormActionWithPermissions';
import IconButton from '@fbcnms/ui//components/design-system/IconButton';
import React from 'react';
import RemoveServiceMutation from '../../mutations/RemoveServiceMutation';
import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {withSnackbar} from 'notistack';

type Props = $ReadOnly<{|
  className?: string,
  service: Service,
  onServiceRemoved: () => void,
  ...WithAlert,
  ...WithSnackbarProps,
|}>;

class ServiceDeleteButton extends React.Component<Props> {
  render() {
    const {className} = this.props;
    return (
      <FormActionWithPermissions
        permissions={{entity: 'service', action: 'delete'}}>
        <IconButton
          icon={DeleteIcon}
          skin="gray"
          className={className}
          onClick={this.removeService}
        />
      </FormActionWithPermissions>
    );
  }

  removeService = () => {
    ServerLogger.info(LogEvents.DELETE_SERVICE_BUTTON_CLICKED, {
      source: 'service_details',
    });
    const {service} = this.props;
    const serviceId = service.id;
    this.props
      .confirm({
        message: 'Are you sure you want to delete this service?',
        confirmLabel: 'Delete',
      })
      .then(confirmed => {
        if (!confirmed) {
          return;
        }

        const variables: RemoveServiceMutationVariables = {
          id: nullthrows(serviceId),
        };

        const updater = store => {
          this.props.onServiceRemoved();
          // $FlowFixMe (T62907961) Relay flow types
          store.delete(serviceId);
        };

        const callbacks: MutationCallbacks<RemoveServiceMutationResponse> = {
          onCompleted: (response, errors) => {
            if (errors && errors[0]) {
              this.props.alert('Failed removing service');
            }
          },
          onError: (_error: Error) => {
            this.props.alert('Failed removing service');
          },
        };

        RemoveServiceMutation(variables, callbacks, updater);
      });
  };
}

export default withAlert(withSnackbar(ServiceDeleteButton));
