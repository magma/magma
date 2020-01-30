/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {AddEditWorkOrderTypeCard_editingWorkOrderType} from './__generated__/AddEditWorkOrderTypeCard_editingWorkOrderType.graphql';
import type {ContextRouter} from 'react-router-dom';
import type {EditWorkOrderTypeMutationResponse} from '../../mutations/__generated__/EditWorkOrderTypeMutation.graphql';
import type {WithStyles} from '@material-ui/core';

import AddEditWorkOrderTypeCard from './AddEditWorkOrderTypeCard';
import Button from '@fbcnms/ui/components/design-system/Button';
import ConfigureTitle from '@fbcnms/ui/components/ConfigureTitle';
import InventoryQueryRenderer from '../InventoryQueryRenderer';
import React from 'react';
import Table from '@fbcnms/ui/components/design-system/Table/Table';
import withInventoryErrorBoundary from '../../common/withInventoryErrorBoundary';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {graphql} from 'relay-runtime';
import {sortLexicographically} from '@fbcnms/ui/utils/displayUtils';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  root: {
    width: '100%',
  },
  paper: {
    flexGrow: 1,
    overflowY: 'hidden',
  },
  typesList: {
    flexGrow: 1,
    padding: '24px',
  },
  content: {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'flex-start',
  },
  listItem: {
    marginBottom: theme.spacing(),
  },
  addButton: {
    alignSelf: 'flex-end',
    marginLeft: 'auto',
  },
  addButtonContainer: {
    marginBottom: '20px',
    display: 'flex',
  },
});

type Props = ContextRouter & WithStyles<typeof styles> & {};

type State = {
  dialogKey: number,
  showAddEditCard: boolean,
  editingWorkOrderType: ?AddEditWorkOrderTypeCard_editingWorkOrderType,
};

const workOrderTypesQuery = graphql`
  query WorkOrderTypesQuery {
    workOrderTypes(first: 500) @connection(key: "Configure_workOrderTypes") {
      edges {
        node {
          id
          name
          description
          ...AddEditWorkOrderTypeCard_editingWorkOrderType
        }
      }
    }
  }
`;

class WorkOrderTypes extends React.Component<Props, State> {
  state = {
    dialogKey: 1,
    showAddEditCard: false,
    editingWorkOrderType: null,
  };

  render() {
    const {classes} = this.props;
    const {showAddEditCard, editingWorkOrderType} = this.state;
    return (
      <InventoryQueryRenderer
        query={workOrderTypesQuery}
        variables={{}}
        render={props => {
          const {workOrderTypes} = props;
          if (showAddEditCard) {
            return (
              <div className={classes.paper}>
                <AddEditWorkOrderTypeCard
                  key={'new_work_order_type@' + this.state.dialogKey}
                  open={showAddEditCard}
                  onClose={this.hideAddEditWorkOrderTypeCard}
                  onSave={this.saveWorkOrder}
                  editingWorkOrderType={editingWorkOrderType}
                />
              </div>
            );
          }
          return (
            <div className={classes.typesList}>
              <div className={classes.addButtonContainer}>
                <ConfigureTitle
                  title="Work Order Templates"
                  subtitle="Create and manage reusable work orders."
                />
                <Button
                  className={classes.addButton}
                  onClick={() => this.showAddEditWorkOrderTypeCard(null)}>
                  Add Work Order Template
                </Button>
              </div>
              <div className={classes.root}>
                <Table
                  className={classes.table}
                  data={workOrderTypes.edges
                    .map(edge => edge.node)
                    .sort((woTypeA, woTypeB) =>
                      sortLexicographically(woTypeA.name, woTypeB.name),
                    )}
                  columns={[
                    {
                      key: 'name',
                      title: 'Work order template',
                      render: row => (
                        <Button
                          variant="text"
                          onClick={() =>
                            this.showAddEditWorkOrderTypeCard(row)
                          }>
                          {row.name}
                        </Button>
                      ),
                    },
                    {
                      key: 'description',
                      title: 'Description',
                      render: row => row.description,
                    },
                  ]}
                />
              </div>
            </div>
          );
        }}
      />
    );
  }

  showAddEditWorkOrderTypeCard = (
    woType: ?AddEditWorkOrderTypeCard_editingWorkOrderType,
  ) => {
    ServerLogger.info(LogEvents.ADD_WORK_ORDER_TYPE_BUTTON_CLICKED);
    this.setState({editingWorkOrderType: woType, showAddEditCard: true});
  };

  hideAddEditWorkOrderTypeCard = () =>
    this.setState(prevState => ({
      editingWorkOrderType: null,
      showAddEditCard: false,
      dialogKey: prevState.dialogKey + 1,
    }));

  saveWorkOrder = (
    workOrderType: $PropertyType<
      EditWorkOrderTypeMutationResponse,
      'editWorkOrderType',
    >,
  ) => {
    ServerLogger.info(LogEvents.SAVE_WORK_ORDER_TYPE_BUTTON_CLICKED);
    this.setState(prevState => {
      if (workOrderType) {
        return {
          dialogKey: prevState.dialogKey + 1,
          showAddEditCard: false,
        };
      }
    });
  };
}

export default withStyles(styles)(
  withRouter(withInventoryErrorBoundary(WorkOrderTypes)),
);
