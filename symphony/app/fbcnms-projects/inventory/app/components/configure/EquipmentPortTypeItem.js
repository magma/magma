/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {EquipmentPortTypeItem_equipmentPortType} from './__generated__/EquipmentPortTypeItem_equipmentPortType.graphql';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';

import CardSection from '../CardSection';
import ConfigureExpansionPanel from './ConfigureExpansionPanel';
import DynamicPropertyTypesGrid from '../DynamicPropertyTypesGrid';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import React from 'react';
import RemoveEquipmentPortTypeMutation from '../../mutations/RemoveEquipmentPortTypeMutation';
import SettingsEthernetIcon from '@material-ui/icons/SettingsEthernet';
import {createFragmentContainer, graphql} from 'react-relay';
import {withStyles} from '@material-ui/core/styles';

import withAlert from '@fbcnms/ui/components/Alert/withAlert';

type Props = {
  equipmentPortType: EquipmentPortTypeItem_equipmentPortType,
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

class EquipmentPortTypeItem extends React.Component<Props> {
  render() {
    const {classes, equipmentPortType, onEdit} = this.props;
    return (
      <div>
        <ExpansionPanel>
          <ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
            <ConfigureExpansionPanel
              entityName="portType"
              icon={<SettingsEthernetIcon />}
              name={equipmentPortType.name}
              instanceCount={equipmentPortType.numberOfPortDefinitions}
              instanceNameSingular="port type"
              instanceNamePlural="port types"
              onDelete={this.onDelete}
              onEdit={onEdit}
            />
          </ExpansionPanelSummary>
          <ExpansionPanelDetails className={classes.detailsRoot}>
            <div className={classes.detailsContainer}>
              <CardSection title="Port Properties">
                <DynamicPropertyTypesGrid
                  key={equipmentPortType.id}
                  propertyTypes={equipmentPortType.propertyTypes}
                />
              </CardSection>
              <CardSection title="Link Properties">
                <DynamicPropertyTypesGrid
                  key={equipmentPortType.id}
                  propertyTypes={equipmentPortType.linkPropertyTypes}
                />
              </CardSection>
            </div>
          </ExpansionPanelDetails>
        </ExpansionPanel>
      </div>
    );
  }

  onDelete = () => {
    this.props
      .confirm(
        `Are you sure you want to delete "${this.props.equipmentPortType.name}"?`,
      )
      .then(confirm => {
        if (confirm) {
          RemoveEquipmentPortTypeMutation(
            {id: this.props.equipmentPortType.id},
            {
              onError: (error: any) => {
                this.props.alert('Error: ' + error.source?.errors[0]?.message);
              },
            },
            // $FlowFixMe (T62907961) Relay flow types
            store => store.delete(this.props.equipmentPortType.id),
          );
        }
      });
  };
}

export default withAlert(
  withStyles(styles)(
    createFragmentContainer(EquipmentPortTypeItem, {
      equipmentPortType: graphql`
        fragment EquipmentPortTypeItem_equipmentPortType on EquipmentPortType {
          id
          name
          numberOfPortDefinitions
          propertyTypes {
            ...DynamicPropertyTypesGrid_propertyTypes
          }
          linkPropertyTypes {
            ...DynamicPropertyTypesGrid_propertyTypes
          }
        }
      `,
    }),
  ),
);
