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
  RemoveProjectMutationResponse,
  RemoveProjectMutationVariables,
} from '../../mutations/__generated__/RemoveProjectMutation.graphql';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithSnackbarProps} from 'notistack';

import Button from '@fbcnms/ui/components/design-system/Button';
import DeleteOutlineIcon from '@material-ui/icons/DeleteOutline';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import React from 'react';
import RemoveProjectTypeMutation from '../../mutations/RemoveProjectTypeMutation';
import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {withSnackbar} from 'notistack';

type Props = {
  className?: string,
  projectType: {id: string, name: string},
} & WithAlert &
  WithSnackbarProps;

class ProjectTypeDeleteButton extends React.Component<Props> {
  render() {
    return (
      <FormAction>
        <Button
          className={this.props.className}
          variant="text"
          skin="primary"
          onClick={this.removeProject}>
          <DeleteOutlineIcon />
        </Button>
      </FormAction>
    );
  }

  removeProject = () => {
    ServerLogger.info(LogEvents.DELETE_PROJECT_TYPE_BUTTON_CLICKED, {
      source: 'project_templates',
    });
    const {projectType} = this.props;
    const projectTypeId = projectType.id;
    this.props
      .confirm({
        message: 'Are you sure you want to delete this project template?',
        confirmLabel: 'Delete',
      })
      .then(confirmed => {
        if (!confirmed) {
          return;
        }

        const variables: RemoveProjectMutationVariables = {
          id: nullthrows(projectTypeId),
        };

        const updater = store => {
          // $FlowFixMe (T62907961) Relay flow types
          store.delete(projectTypeId);
        };

        const callbacks: MutationCallbacks<RemoveProjectMutationResponse> = {
          onCompleted: (response, errors) => {
            if (errors && errors[0]) {
              this.props.alert('Failed removing project template');
            }
          },
          onError: (_error: Error) => {
            this.props.alert('Failed removing project template');
          },
        };

        RemoveProjectTypeMutation(variables, callbacks, updater);
      });
  };
}

export default withAlert(withSnackbar(ProjectTypeDeleteButton));
