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
import type {
  RemoveProjectMutationResponse,
  RemoveProjectMutationVariables,
} from '../../mutations/__generated__/RemoveProjectMutation.graphql';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithSnackbarProps} from 'notistack';
import type {WithStyles} from '@material-ui/core';

import DeleteOutlineIcon from '@material-ui/icons/DeleteOutline';
import React from 'react';
import RemoveProjectTypeMutation from '../../mutations/RemoveProjectTypeMutation';
import SymphonyTheme from '@fbcnms/ui/theme/symphony';
import classNames from 'classnames';
import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {withSnackbar} from 'notistack';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  deleteButton: {
    cursor: 'pointer',
    color: theme.palette.primary.main,
    width: '32px',
    height: '32px',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    cursor: 'pointer',
    '&:hover': {
      color: SymphonyTheme.palette.B700,
    },
  },
});

type Props = {
  className?: string,
  projectType: {id: string, name: string},
} & WithStyles<typeof styles> &
  WithAlert &
  WithSnackbarProps;

class ProjectTypeDeleteButton extends React.Component<Props> {
  render() {
    const {classes, className} = this.props;
    return (
      <div className={classNames(classes.deleteButton, className)}>
        <DeleteOutlineIcon onClick={this.removeProject} />
      </div>
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

export default withStyles(styles)(
  withAlert(withSnackbar(ProjectTypeDeleteButton)),
);
