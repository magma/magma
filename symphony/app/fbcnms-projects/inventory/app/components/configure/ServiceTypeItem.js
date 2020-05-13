/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {ServiceTypeItem_serviceType} from './__generated__/ServiceTypeItem_serviceType.graphql';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';

import ConfigureExpansionPanel from './ConfigureExpansionPanel';
import DynamicPropertyTypesGrid from '../DynamicPropertyTypesGrid';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import LinearScaleIcon from '@material-ui/icons/LinearScale';
import React from 'react';
import RemoveServiceTypeMutation from '../../mutations/RemoveServiceTypeMutation';
import ServiceEndpointDefinitionStaticTable from './ServiceEndpointDefinitionStaticTable';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {ConnectionHandler} from 'relay-runtime';
import {createFragmentContainer, graphql} from 'react-relay';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  serviceType: ServiceTypeItem_serviceType,
  onEdit: () => void,
} & WithAlert &
  WithStyles<typeof styles>;

const styles = {
  detailsContainer: {
    width: '100%',
  },
  section: {
    marginBottom: '24px',
  },
};

class ServiceTypeItem extends React.Component<Props> {
  render() {
    const {classes, serviceType, onEdit} = this.props;
    return (
      <div>
        <ExpansionPanel>
          <ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
            <ConfigureExpansionPanel
              icon={<LinearScaleIcon />}
              name={serviceType.name}
              instanceCount={serviceType.numberOfServices}
              instanceNameSingular="service"
              instanceNamePlural="services"
              onDelete={this.onDelete}
              allowDelete={true}
              onEdit={onEdit}
            />
          </ExpansionPanelSummary>
          <ExpansionPanelDetails>
            <div className={classes.detailsContainer}>
              <div className={classes.section}>
                <DynamicPropertyTypesGrid
                  key={serviceType.id}
                  propertyTypes={serviceType.propertyTypes}
                />
              </div>
              <div className={classes.section}>
                <ServiceEndpointDefinitionStaticTable
                  serviceEndpointDefinitions={serviceType.endpointDefinitions}
                />
              </div>
            </div>
          </ExpansionPanelDetails>
        </ExpansionPanel>
      </div>
    );
  }

  onDelete = () => {
    const msg = `Are you sure you want to delete "${this.props.serviceType.name}"? The service type, and all it's instances will be deleted soon, in the background`;
    this.props.confirm(msg).then(confirm => {
      if (!confirm) {
        return;
      }
      RemoveServiceTypeMutation(
        {id: this.props.serviceType.id},
        {
          onError: (error: any) => {
            this.props.alert('Error: ' + error.source?.errors[0]?.message);
          },
        },
        store => {
          // $FlowFixMe (T62907961) Relay flow types
          const rootQuery = store.getRoot();
          const serviceTypes = ConnectionHandler.getConnection(
            rootQuery,
            'ServiceTypes_serviceTypes',
          );
          ConnectionHandler.deleteNode(
            // $FlowFixMe (T62907961) Relay flow types
            serviceTypes,
            this.props.serviceType.id,
          );
          // $FlowFixMe (T62907961) Relay flow types
          store.delete(this.props.serviceType.id);
        },
      );
    });
  };
}

export default withStyles(styles)(
  withAlert(
    createFragmentContainer(ServiceTypeItem, {
      serviceType: graphql`
        fragment ServiceTypeItem_serviceType on ServiceType {
          id
          name
          discoveryMethod
          propertyTypes {
            ...PropertyTypeFormField_propertyType
          }
          endpointDefinitions {
            ...ServiceEndpointDefinitionStaticTable_serviceEndpointDefinitions
          }
          numberOfServices
        }
      `,
    }),
  ),
);
