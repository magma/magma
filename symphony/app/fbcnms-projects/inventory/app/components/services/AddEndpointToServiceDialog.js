/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {EndpointDef} from './ServiceEndpointsMenu';
import type {Equipment, EquipmentPort} from '../../common/Equipment';
import type {ServicePanel_service} from './__generated__/ServicePanel_service.graphql';
import type {WithStyles} from '@material-ui/core';

import AvailablePortsTable from '../AvailablePortsTable';
import Button from '@fbcnms/ui/components/design-system/Button';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import EquipmentComparisonViewQueryRenderer from '../comparison_view/EquipmentComparisonViewQueryRenderer';
import InventoryQueryRenderer from '../InventoryQueryRenderer';
import PowerSearchLinkFirstEquipmentResultsTable from './PowerSearchLinkFirstEquipmentResultsTable';
import React from 'react';
import Strings from '@fbcnms/strings/Strings';
import Text from '@fbcnms/ui/components/design-system/Text';
import fbt from 'fbt';
import nullthrows from '@fbcnms/util/nullthrows';
import symphony from '@fbcnms/ui/theme/symphony';
import {OperatorMap} from '../comparison_view/ComparisonViewTypes';
import {WizardContextProvider} from '@fbcnms/ui/components/design-system/Wizard/WizardContext';
import {generateTempId} from '../../common/EntUtils';
import {graphql} from 'react-relay';

import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = {
  root: {
    minWidth: '80vh',
    paddingTop: '0px',
    paddingLeft: '32px',
    paddingRight: '32px',
  },
  searchResults: {
    backgroundColor: symphony.palette.white,
    height: '100%',
  },
  title: {
    display: 'block',
  },
  subtitle: {
    display: 'block',
    color: symphony.palette.D500,
  },
  footer: {
    padding: '16px 24px',
    boxShadow: '0px 1px 4px 0px rgba(0, 0, 0, 0.17)',
  },
  actionButton: {
    '&&': {
      marginLeft: '12px',
    },
  },
};

type Props = {
  service: ServicePanel_service,
  onClose: () => void,
  onAddEndpoint: (port: EquipmentPort) => void,
  endpointDef: ?EndpointDef,
} & WithStyles<typeof styles>;

type State = {
  activeEquipement: ?Equipment,
  activeStep: number,
  activePort: ?EquipmentPort,
};

const steps = [
  fbt(
    'Select the equipment of the endpoint port you want to add to the service.',
    'Subtitle in dialog to add endpoint to service',
  ),
  fbt(
    'Select the port of the endpoint port.',
    'Subtitle in dialog to add endpoint to service',
  ),
];

const addEndpointToServiceDialogQuery = graphql`
  query AddEndpointToServiceDialogQuery($filters: [PortFilterInput!]!) {
    portSearch(filters: $filters, limit: 50) {
      ports {
        id
        definition {
          id
          name
        }
        ...AvailablePortsTable_ports
      }
    }
  }
`;

class AddEndpointToServiceDialog extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {
      activeEquipement: null,
      activeStep: 0,
      activePort: null,
    };
  }

  handleElementSelected = equipment => {
    this.setState({
      activeEquipement: equipment,
    });
  };

  handlePortSelected = (port: EquipmentPort) => {
    this.setState({
      activePort: port,
    });
  };

  getStepContent = () => {
    const {classes, endpointDef} = this.props;
    switch (this.state.activeStep) {
      case 0:
        return (
          <EquipmentComparisonViewQueryRenderer
            limit={50}
            initialFilters={
              endpointDef
                ? [
                    {
                      id: generateTempId(),
                      key: 'equipment_type',
                      name: 'equipment_type',
                      operator: OperatorMap.is_one_of,
                      idSet: [endpointDef.equipmentTypeID],
                    },
                  ]
                : []
            }>
            {props => (
              <div className={classes.searchResults}>
                <PowerSearchLinkFirstEquipmentResultsTable
                  equipment={props.equipment}
                  onEquipmentSelected={this.handleElementSelected}
                  selectedEquipment={this.state.activeEquipement}
                />
              </div>
            )}
          </EquipmentComparisonViewQueryRenderer>
        );
      case 1:
        return (
          <InventoryQueryRenderer
            query={addEndpointToServiceDialogQuery}
            variables={{
              filters: [
                {
                  filterType: 'SERVICE_INST',
                  operator: 'IS_NOT_ONE_OF',
                  idSet: [this.props.service.id],
                },
                {
                  filterType: 'PORT_INST_EQUIPMENT',
                  operator: 'IS_ONE_OF',
                  idSet: [nullthrows(this.state.activeEquipement).id],
                },
              ],
            }}
            render={props => {
              const {portSearch} = props;
              return (
                <AvailablePortsTable
                  equipment={nullthrows(this.state.activeEquipement)}
                  ports={portSearch.ports}
                  selectedPort={this.state.activePort}
                  onPortSelected={this.handlePortSelected}
                />
              );
            }}
          />
        );
      default:
        return '';
    }
  };

  handleNext = () => {
    this.setState(state => ({
      activeStep: state.activeStep + 1,
    }));
  };

  handleBack = () => {
    this.setState(state => ({
      activeStep: state.activeStep - 1,
    }));
  };

  handleReset = () => {
    this.setState({
      activeStep: 0,
    });
  };

  getSubtitle = (activeStep: number) => {
    return steps[activeStep];
  };

  render() {
    const {classes, onAddEndpoint, service, onClose} = this.props;
    const {activeStep, activePort, activeEquipement} = this.state;

    return (
      <>
        <DialogTitle>
          <Text className={classes.title} variant="h6">
            {`${fbt(
              'Add endpoint to ' + fbt.param('service name', service.name),
              'Title of dialog for adding an endpoint to a service',
            )}
            `}
          </Text>
          <Text className={classes.subtitle} variant="subtitle2" color="light">
            {this.getSubtitle(activeStep)}
          </Text>
        </DialogTitle>
        <DialogContent div className={classes.root}>
          <WizardContextProvider>{this.getStepContent()}</WizardContextProvider>
        </DialogContent>
        <DialogActions className={classes.footer}>
          {activeStep === 0 ? (
            <>
              <Button skin="gray" onClick={onClose}>
                {Strings.common.cancelButton}
              </Button>
              <Button
                disabled={activeEquipement === null}
                onClick={this.handleNext}
                className={classes.actionButton}>
                {Strings.common.nextButton}
              </Button>
            </>
          ) : (
            <>
              <Button skin="gray" onClick={this.handleBack}>
                {Strings.common.backButton}
              </Button>
              <Button
                disabled={activePort === null}
                color="primary"
                onClick={() => onAddEndpoint(nullthrows(activePort))}
                className={classes.actionButton}>
                {Strings.common.addButton}
              </Button>
            </>
          )}
        </DialogActions>
      </>
    );
  }
}

export default withRouter(withStyles(styles)(AddEndpointToServiceDialog));
