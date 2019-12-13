/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {WithStyles} from '@material-ui/core';
import type {WorkOrder} from '../../common/WorkOrder';

import DescriptionIcon from '@material-ui/icons/Description';
import IconButton from '@material-ui/core/IconButton';
import InventoryQueryRenderer from '../../components/InventoryQueryRenderer';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
import React from 'react';
import {graphql} from 'relay-runtime';
import {withStyles} from '@material-ui/core/styles';

type Props = WithStyles<typeof styles> & {
  onSelect: (workOrder: WorkOrder) => void,
  onNavigateToWorkOrder: (workOrderId: string) => void,
};

const styles = theme => ({
  root: {
    backgroundColor: theme.palette.background.paper,
    minWidth: '200px',
  },
  heading: {
    display: 'inline-flex',
    paddingRight: theme.spacing(3),
    width: '100%',
    justifyContent: 'space-between',
  },
  button: {
    margin: theme.spacing(),
  },
  title: {
    lineHeight: '100%',
    marginBottom: '12px',
    marginLeft: '12px',
    fontWeight: 'bold',
  },
  listItem: {
    paddingLeft: '30px',
  },
  listItemText: {
    color: theme.palette.dark,
  },
});

graphql`
  fragment WorkOrdersPane_workOrder on WorkOrder @relay(mask: false) {
    id
    name
  }
`;

const WorkOrdersPaneQuery = graphql`
  query WorkOrdersPaneQuery {
    workOrders(first: 50, showCompleted: false)
      @connection(key: "WorkOrdersPane_workOrders") {
      edges {
        node {
          ...WorkOrdersPane_workOrder @relay(mask: false)
        }
      }
    }
  }
`;

// This is a QueryRenderer that uses the query in it.
// eslint-disable-next-line relay/generated-flow-types
class WorkOrdersPane extends React.Component<Props> {
  render() {
    const {classes, onNavigateToWorkOrder} = this.props;
    return (
      <InventoryQueryRenderer
        query={WorkOrdersPaneQuery}
        variables={{}}
        render={props => {
          const workOrders = props.workOrders.edges.map(edge => edge.node);
          return (
            <div>
              <List component="nav" dense={true} className={classes.root}>
                {workOrders.length > 0 ? (
                  workOrders.map(workOrder => (
                    <ListItem
                      button
                      key={workOrder.id}
                      className={classes.listItem}
                      onClick={() => this.props.onSelect(workOrder)}>
                      <ListItemText
                        classes={{primary: classes.listItemText}}
                        primary={workOrder.name}
                      />
                      <ListItemSecondaryAction>
                        <IconButton
                          onClick={_ => onNavigateToWorkOrder(workOrder.id)}>
                          <DescriptionIcon />
                        </IconButton>
                      </ListItemSecondaryAction>
                    </ListItem>
                  ))
                ) : (
                  <ListItem
                    button
                    key={'placeholder'}
                    className={classes.listItem}>
                    <ListItemText
                      primary={'No work orders found'}
                      classes={{primary: classes.listItemText}}
                    />
                  </ListItem>
                )}
              </List>
            </div>
          );
        }}
      />
    );
  }
}

export default withStyles(styles)(WorkOrdersPane);
