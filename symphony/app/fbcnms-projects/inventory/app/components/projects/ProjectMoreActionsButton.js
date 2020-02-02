/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {ProjectMoreActionsButton_project} from './__generated__/ProjectMoreActionsButton_project.graphql.js';
import type {
  RemoveProjectMutationResponse,
  RemoveProjectMutationVariables,
} from '../../mutations/__generated__/RemoveProjectMutation.graphql';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithSnackbarProps} from 'notistack';
import type {WithStyles} from '@material-ui/core';

import Button from '@fbcnms/ui/components/design-system/Button';
import DeleteOutlineIcon from '@material-ui/icons/DeleteOutline';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import React from 'react';
import RemoveProjectMutation from '../../mutations/RemoveProjectMutation';
import classNames from 'classnames';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {createFragmentContainer, graphql} from 'react-relay';
import {withSnackbar} from 'notistack';
import {withStyles} from '@material-ui/core/styles';

const styles = () => ({
  deleteButton: {
    minWidth: 'unset',
    paddingTop: '2px',
  },
});

type Props = {
  className?: string,
  project: ProjectMoreActionsButton_project,
  onProjectRemoved: () => void,
} & WithStyles<typeof styles> &
  WithAlert &
  WithSnackbarProps;

class ProjectMoreActionsButton extends React.Component<Props> {
  render() {
    const {className, classes} = this.props;
    return (
      <FormAction>
        <Button
          className={classNames(className, classes.deleteButton)}
          variant="text"
          skin="gray"
          onClick={this.removeProject}>
          <DeleteOutlineIcon />
        </Button>
      </FormAction>
    );
  }

  removeProject = () => {
    ServerLogger.info(LogEvents.DELETE_PROJECT_BUTTON_CLICKED, {
      source: 'project_details',
    });
    const {project} = this.props;
    if (project.numberOfWorkOrders > 0) {
      this.props.alert('Cannot delete a project that is in use');
    }
    const projectId = project.id;
    this.props
      .confirm({
        message: 'Are you sure you want to delete this project?',
        confirmLabel: 'Delete',
      })
      .then(confirmed => {
        if (!confirmed) {
          return;
        }

        const variables: RemoveProjectMutationVariables = {
          id: projectId,
        };

        const updater = store => {
          this.props.onProjectRemoved();
          store.delete(projectId);
        };

        const callbacks: MutationCallbacks<RemoveProjectMutationResponse> = {
          onCompleted: (response, errors) => {
            if (errors && errors[0]) {
              this.props.alert('Failed removing project');
            }
          },
          onError: (_error: Error) => {
            this.props.alert('Failed removing project');
          },
        };

        RemoveProjectMutation(variables, callbacks, updater);
      });
  };
}

export default withStyles(styles)(
  withAlert(
    withSnackbar(
      createFragmentContainer(ProjectMoreActionsButton, {
        project: graphql`
          fragment ProjectMoreActionsButton_project on Project {
            id
            name
            numberOfWorkOrders
          }
        `,
      }),
    ),
  ),
);
