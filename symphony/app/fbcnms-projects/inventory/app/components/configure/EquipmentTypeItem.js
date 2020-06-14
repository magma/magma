/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {EquipmentTypeItem_equipmentType} from './__generated__/EquipmentTypeItem_equipmentType.graphql';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';

import ConfigureExpansionPanel from './ConfigureExpansionPanel';
import DynamicPropertyTypesGrid from '../DynamicPropertyTypesGrid';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import PortDefinitionsTable from './PortDefinitionsTable';
import PositionDefinitionsTable from './PositionDefinitionsTable';
import React from 'react';
import RemoveEquipmentTypeMutation from '../../mutations/RemoveEquipmentTypeMutation';
import RouterIcon from '@material-ui/icons/Router';
import {createFragmentContainer, graphql} from 'react-relay';
import {withStyles} from '@material-ui/core/styles';

import withAlert from '@fbcnms/ui/components/Alert/withAlert';

type Props = {
  equipmentType: EquipmentTypeItem_equipmentType,
  onEdit: () => void,
} & WithAlert &
  WithStyles<typeof styles>;

const styles = {
  detailsRoot: {
    display: 'block',
  },
  detailsContainer: {
    width: '100%',
  },
  section: {
    marginBottom: '24px',
  },
};

class EquipmentTypeItem extends React.Component<Props> {
  render() {
    const {classes, equipmentType, onEdit} = this.props;

    if (equipmentType == null) {
      return null;
    }

    return (
      <div>
        <ExpansionPanel>
          <ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
            <ConfigureExpansionPanel
              entityName="equipmentType"
              icon={<RouterIcon />}
              name={equipmentType.name}
              instanceCount={equipmentType.numberOfEquipment}
              instanceNameSingular="equipment instance"
              instanceNamePlural="equipment instances"
              onDelete={this.onDelete}
              onEdit={onEdit}
            />
          </ExpansionPanelSummary>
          <ExpansionPanelDetails className={classes.detailsRoot}>
            <div className={classes.detailsContainer}>
              <div className={classes.section}>
                <DynamicPropertyTypesGrid
                  key={equipmentType.id}
                  propertyTypes={equipmentType.propertyTypes}
                />
              </div>
              <div className={classes.section}>
                <PositionDefinitionsTable
                  positionDefinitions={equipmentType.positionDefinitions}
                />
              </div>
              <div className={classes.section}>
                <PortDefinitionsTable
                  portDefinitions={equipmentType.portDefinitions}
                />
              </div>
            </div>
          </ExpansionPanelDetails>
        </ExpansionPanel>
      </div>
    );
  }

  onDelete = () => {
    this.props
      .confirm(
        `Are you sure you want to delete "${this.props.equipmentType.name}"?`,
      )
      .then(confirm => {
        if (confirm) {
          RemoveEquipmentTypeMutation(
            {id: this.props.equipmentType.id},
            {
              onError: (error: any) => {
                this.props.alert('Error: ' + error.source?.errors[0]?.message);
              },
            },
            // $FlowFixMe (T62907961) Relay flow types
            store => store.delete(this.props.equipmentType.id),
          );
        }
      });
  };
}

export default withAlert(
  withStyles(styles)(
    createFragmentContainer(EquipmentTypeItem, {
      equipmentType: graphql`
        fragment EquipmentTypeItem_equipmentType on EquipmentType {
          id
          name
          propertyTypes {
            ...DynamicPropertyTypesGrid_propertyTypes
          }
          positionDefinitions {
            ...PositionDefinitionsTable_positionDefinitions
          }
          portDefinitions {
            ...PortDefinitionsTable_portDefinitions
          }
          numberOfEquipment
        }
      `,
    }),
  ),
);
