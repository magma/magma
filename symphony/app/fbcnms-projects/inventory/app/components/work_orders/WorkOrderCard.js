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
import React from 'react';
import WorkOrderDetails from './WorkOrderDetails';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {graphql} from 'react-relay';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  workOrderId: ?string,
  onWorkOrderExecuted: () => void,
  onWorkOrderRemoved: () => void,
} & WithStyles<typeof styles> &
  ContextRouter;

type State = {
  isLoadingDocument: boolean,
};

const styles = theme => ({
  root: {
    height: '100%',
    width: '100%',
    padding: '40px 32px',
    overflow: 'hidden',
  },
  tabs: {
    backgroundColor: theme.palette.common.white,
  },
  docs: {
    margin: '24px 24px 0px 24px',
    backgroundColor: theme.palette.common.white,
  },
  titleText: {
    fontWeight: 500,
  },
  section: {
    marginBottom: theme.spacing(3),
  },
  tabContainer: {
    width: 'auto',
  },
  cardContentRoot: {
    '&:last-child': {
      paddingBottom: '0px',
    },
  },
  iconButton: {
    padding: '0px',
    marginLeft: theme.spacing(),
  },
});

const workOrderQuery = graphql`
  query WorkOrderCardQuery($workOrderId: ID!) {
    workOrder(id: $workOrderId) {
      id
      name
      ...WorkOrderDetails_workOrder
    }
  }
`;

class WorkOrderCard extends React.Component<Props, State> {
  state = {
    isLoadingDocument: false,
  };

  render() {
    const {
      classes,
      workOrderId,
      onWorkOrderExecuted,
      onWorkOrderRemoved,
    } = this.props;
    return (
      <InventoryQueryRenderer
        query={workOrderQuery}
        variables={{
          workOrderId,
        }}
        render={props => {
          const {workOrder} = props;
          return (
            <div className={classes.root}>
              <WorkOrderDetails
                workOrder={workOrder}
                onWorkOrderRemoved={onWorkOrderRemoved}
                onWorkOrderExecuted={onWorkOrderExecuted}
                onCancelClicked={this.navigateToMainPage}
              />
            </div>
          );
        }}
      />
    );
  }

  navigateToMainPage = () => {
    ServerLogger.info(LogEvents.WORK_ORDERS_SEARCH_NAV_CLICKED, {
      source: 'work_order_details',
    });
    const {match} = this.props;
    this.props.history.push(match.url);
  };
}

export default withRouter(withStyles(styles)(WorkOrderCard));
