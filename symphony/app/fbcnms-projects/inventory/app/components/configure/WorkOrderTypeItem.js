/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithSnackbarProps} from 'notistack';
import type {WithStyles} from '@material-ui/core';
import type {WorkOrderTypeItem_workOrderType} from './__generated__/WorkOrderTypeItem_workOrderType.graphql';

import ConfigureExpansionPanel from './ConfigureExpansionPanel';
import DynamicPropertyTypesGrid from '../DynamicPropertyTypesGrid';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import React from 'react';
import RemoveWorkOrderTypeMutation from '../../mutations/RemoveWorkOrderTypeMutation';
import WorkIcon from '@material-ui/icons/Work';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {ConnectionHandler} from 'relay-runtime';
import {createFragmentContainer, graphql} from 'react-relay';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  workOrderType: WorkOrderTypeItem_workOrderType,
  onEdit: () => void,
} & WithAlert &
  WithStyles<typeof styles> &
  WithSnackbarProps;

const styles = {
  properties: {
    marginBottom: '24px',
  },
};

class WorkOrderTypeItem extends React.Component<Props> {
  render() {
    const {classes, workOrderType, onEdit} = this.props;
    return (
      <div>
        <ExpansionPanel>
          <ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
            <ConfigureExpansionPanel
              icon={<WorkIcon />}
              name={workOrderType.name}
              instanceCount={workOrderType.numberOfWorkOrders}
              instanceNameSingular="work order"
              instanceNamePlural="work orders"
              onEdit={onEdit}
              onDelete={this.onDelete}
            />
          </ExpansionPanelSummary>
          <ExpansionPanelDetails>
            <div className={classes.properties}>
              <DynamicPropertyTypesGrid
                key={workOrderType.id}
                propertyTypes={workOrderType.propertyTypes}
              />
            </div>
          </ExpansionPanelDetails>
        </ExpansionPanel>
      </div>
    );
  }

  onDelete = () => {
    this.props
      .confirm(
        `Are you sure you want to delete "${this.props.workOrderType.name}"?`,
      )
      .then(confirm => {
        if (!confirm) {
          return;
        }
        RemoveWorkOrderTypeMutation(
          {id: this.props.workOrderType.id},
          {
            onError: (error: any) => {
              this.props.alert(`Error: ${error.source?.errors[0]?.message}`);
            },
          },
          store => {
            const rootQuery = store.getRoot();
            const workOrderTypes = ConnectionHandler.getConnection(
              rootQuery,
              'Configure_workOrderTypes',
            );
            ConnectionHandler.deleteNode(
              workOrderTypes,
              this.props.workOrderType.id,
            );
            store.delete(this.props.workOrderType.id);
          },
        );
      });
  };
}

export default withStyles(styles)(
  withAlert(
    createFragmentContainer(WorkOrderTypeItem, {
      workOrderType: graphql`
        fragment WorkOrderTypeItem_workOrderType on WorkOrderType {
          id
          name
          propertyTypes {
            ...DynamicPropertyTypesGrid_propertyTypes
          }
          numberOfWorkOrders
        }
      `,
    }),
  ),
);
