/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {LocationTypeItem_locationType} from './__generated__/LocationTypeItem_locationType.graphql';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithSnackbarProps} from 'notistack';
import type {WithStyles} from '@material-ui/core';

import ConfigureExpansionPanel from './ConfigureExpansionPanel';
import DraggableTableRow from '../draggable/DraggableTableRow';
import DynamicPropertyTypesGrid from '../DynamicPropertyTypesGrid';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import React from 'react';
import RemoveLocationTypeMutation from '../../mutations/RemoveLocationTypeMutation';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {ConnectionHandler} from 'relay-runtime';
import {createFragmentContainer, graphql} from 'react-relay';
import {withSnackbar} from 'notistack';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  locationType: LocationTypeItem_locationType,
  onEdit: () => void,
  position: number,
} & WithAlert &
  WithStyles<typeof styles> &
  WithSnackbarProps;

const styles = theme => ({
  properties: {
    marginBottom: '24px',
    width: '100%',
  },
  draggableRow: {
    display: 'flex',
    paddingLeft: '10px',
    alignItems: 'center',
    boxShadow: theme.shadows[1],
    borderRadius: 4,
  },
  row: {
    flexGrow: 1,
  },
  panel: {
    width: '100%',
    boxShadow: 'none',
  },
  cell: {
    border: 'none',
    paddingLeft: '10px',
  },
  removeBefore: {
    '&:before': {
      backgroundColor: 'transparent',
    },
  },
});

class LocationTypeItem extends React.Component<Props> {
  render() {
    const {classes, locationType, onEdit, position} = this.props;
    return (
      <div>
        <DraggableTableRow
          className={classes.draggableRow}
          draggableCellClassName={classes.cell}
          id={locationType.id}
          index={position}
          key={locationType.id}>
          <ExpansionPanel
            className={classes.panel}
            classes={{root: classes.removeBefore}}>
            <ExpansionPanelSummary
              className={classes.row}
              expandIcon={<ExpandMoreIcon />}>
              <ConfigureExpansionPanel
                entityName="locationType"
                icon={<div>{position + 1}</div>}
                name={locationType.name}
                instanceCount={locationType.numberOfLocations}
                instanceNameSingular="location"
                instanceNamePlural="locations"
                onDelete={this.onDelete}
                onEdit={onEdit}
              />
            </ExpansionPanelSummary>
            <ExpansionPanelDetails>
              <div className={classes.properties}>
                <DynamicPropertyTypesGrid
                  key={locationType.id}
                  propertyTypes={locationType.propertyTypes}
                />
              </div>
            </ExpansionPanelDetails>
          </ExpansionPanel>
        </DraggableTableRow>
      </div>
    );
  }

  onDelete = () => {
    this.props
      .confirm(
        `Are you sure you want to delete "${this.props.locationType.name}"?`,
      )
      .then(confirm => {
        if (!confirm) {
          return;
        }
        RemoveLocationTypeMutation(
          {id: this.props.locationType.id},
          {
            onError: (error: any) => {
              this.props.alert('Error: ' + error.source?.errors[0]?.message);
            },
          },
          store => {
            // $FlowFixMe (T62907961) Relay flow types
            const rootQuery = store.getRoot();
            const locationTypes = ConnectionHandler.getConnection(
              rootQuery,
              'Catalog_locationTypes',
            );
            ConnectionHandler.deleteNode(
              // $FlowFixMe (T62907961) Relay flow types
              locationTypes,
              this.props.locationType.id,
            );
            // $FlowFixMe (T62907961) Relay flow types
            store.delete(this.props.locationType.id);
          },
        );
      });
  };
}

export default withSnackbar(
  withStyles(styles)(
    withAlert(
      createFragmentContainer(LocationTypeItem, {
        locationType: graphql`
          fragment LocationTypeItem_locationType on LocationType {
            id
            name
            index
            propertyTypes {
              ...DynamicPropertyTypesGrid_propertyTypes
            }
            numberOfLocations
          }
        `,
      }),
    ),
  ),
);
