/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ContextRouter} from 'react-router-dom';
import type {WithStyles} from '@material-ui/core';

import InventoryQueryRenderer from '../InventoryQueryRenderer';
import ProjectDetails from './ProjectDetails';
import React from 'react';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {graphql} from 'react-relay';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  projectId: ?string,
  onProjectExecuted: () => void,
  onProjectRemoved: () => void,
} & WithStyles<typeof styles> &
  ContextRouter;

type State = {
  isLoadingDocument: boolean,
};

const styles = _theme => ({
  root: {
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
    padding: '40px 32px',
  },
});

const projectQuery = graphql`
  query ProjectCardQuery($projectId: ID!) {
    project: node(id: $projectId) {
      ... on Project {
        ...ProjectMoreActionsButton_project
        ...ProjectDetails_project
      }
    }
  }
`;

class ProjectCard extends React.Component<Props, State> {
  state = {
    isLoadingDocument: false,
  };

  render() {
    const {
      classes,
      projectId,
      onProjectExecuted,
      onProjectRemoved,
    } = this.props;
    return (
      <>
        <InventoryQueryRenderer
          query={projectQuery}
          variables={{
            projectId,
          }}
          render={props => {
            const {project} = props;
            return (
              <div className={classes.root}>
                <ProjectDetails
                  project={project}
                  onProjectExecuted={onProjectExecuted}
                  navigateToWorkOrder={this.navigateToWorkOrder}
                  onProjectRemoved={onProjectRemoved}
                />
              </div>
            );
          }}
        />
      </>
    );
  }

  navigateToMainPage = () => {
    ServerLogger.info(LogEvents.PROJECTS_SEARCH_NAV_CLICKED, {
      source: 'project_details',
    });
    const {match} = this.props;
    this.props.history.push(match.url);
  };

  navigateToWorkOrder = (WorkOrderId: ?string) => {
    const {history} = this.props;
    if (WorkOrderId) {
      ServerLogger.info(LogEvents.WORK_ORDER_DETAILS_NAV_CLICKED, {
        source: 'project_details',
      });
      history.push(`/workorders/search?workorder=${WorkOrderId}`);
    }
  };
}

export default withRouter(withStyles(styles)(ProjectCard));
