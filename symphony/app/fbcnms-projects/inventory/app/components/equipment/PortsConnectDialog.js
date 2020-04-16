/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Equipment, EquipmentPort} from '../../common/Equipment';
import type {PowerSearchEquipmentResultsTable_equipment} from '../comparison_view/__generated__/PowerSearchEquipmentResultsTable_equipment.graphql';
import type {Property} from '../../common/Property';
import type {WithStyles} from '@material-ui/core';

import AvailablePortsTable from '../AvailablePortsTable';
import Button from '@fbcnms/ui/components/design-system/Button';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import EquipmentComparisonViewQueryRenderer from '../comparison_view/EquipmentComparisonViewQueryRenderer';
import InventoryQueryRenderer from '../InventoryQueryRenderer';
import NodePropertyInput from '../NodePropertyInput';
import PortsConnectConfirmation from './PortsConnectConfirmation';
import PowerSearchEquipmentResultsTable from '../comparison_view/PowerSearchEquipmentResultsTable';
import PropertiesAddEditSection from '../form/PropertiesAddEditSection';
import React from 'react';
import Step from '@material-ui/core/Step';
import StepConnector from '@material-ui/core/StepConnector/StepConnector.js';
import StepLabel from '@material-ui/core/StepLabel';
import Stepper from '@material-ui/core/Stepper';
import nullthrows from '@fbcnms/util/nullthrows';
import update from 'immutability-helper';
import {WizardContextProvider} from '@fbcnms/ui/components/design-system/Wizard/WizardContext';
import {getInitialPropertyFromType} from '../../common/PropertyType';
import {graphql} from 'react-relay';
import {sortPropertiesByIndex} from '../../common/Property';
import {uniqBy} from 'lodash';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  button: {
    marginTop: theme.spacing(),
    marginRight: theme.spacing(),
  },
  content: {
    height: '60vh',
    width: '100%',
  },
  portIdLabel: {
    marginRight: '8px',
    fontWeight: 500,
  },
  connectorActive: {
    '& $connectorLine': {
      borderColor: theme.palette.secondary.main,
    },
  },
  connectorCompleted: {
    '& $connectorLine': {
      borderColor: theme.palette.primary.main,
    },
  },
  connectorDisabled: {
    '& $connectorLine': {
      borderColor: theme.palette.grey[100],
    },
  },
  connectorLine: {
    transition: theme.transitions.create('border-color'),
  },
  root: {
    minWidth: '80vh',
  },
  searchResults: {
    flexGrow: 1,
    display: 'flex',
    flexDirection: 'column',
    backgroundColor: theme.palette.common.white,
    height: '100%',
  },
});

type Props = {
  equipment: Equipment,
  port: EquipmentPort,
  onConnectPorts: (EquipmentPort, Array<Property>) => void,
} & WithStyles<typeof styles>;

type State = {
  activeEquipement: ?Equipment,
  activeStep: number,
  targetPort: ?EquipmentPort,
  linkProperties: Array<Property>,
};

const steps = ['Select Equipment', 'Select Port', 'Link Properties', 'Confirm'];

const portsConnectDialogQuery = graphql`
  query PortsConnectDialogQuery($equipmentId: ID!) {
    equipment: node(id: $equipmentId) {
      ... on Equipment {
        id
        name
        equipmentType {
          id
          name
          portDefinitions {
            id
            name
            visibleLabel
            bandwidth
          }
        }
        descendentsIncludingSelf {
          ports(availableOnly: true) {
            id
            ...AvailablePortsTable_ports
          }
        }
      }
    }
  }
`;

class PortsConnectDialog extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {
      activeEquipement: null,
      activeStep: 0,
      targetPort: null,
      linkProperties: this._getPortProperties(props.port),
    };
  }

  handleElementSelected = equipment => {
    this.setState({
      activeEquipement: equipment,
    });
  };

  handlePortSelected = (port: EquipmentPort) => {
    this.setState(state => ({
      targetPort: port,
      linkProperties: uniqBy(
        [...state.linkProperties, ...this._getPortProperties(port)],
        'id',
      ).sort(sortPropertiesByIndex),
    }));
  };

  _getPortProperties = (port: EquipmentPort): Array<Property> => {
    const propTypes = port.definition.portType?.linkPropertyTypes ?? [];
    return propTypes.map(propType => getInitialPropertyFromType(propType));
  };

  _propertyChangedHandler = index => property =>
    this.setState(prevState => {
      return {
        linkProperties: update(prevState.linkProperties, {
          [index]: {$set: property},
        }),
      };
    });

  getStepContent = () => {
    const {classes} = this.props;
    const {linkProperties, activeEquipement, targetPort} = this.state;
    const EquipmentTable = (props: {
      equipment: PowerSearchEquipmentResultsTable_equipment,
    }) => {
      return (
        <div className={classes.searchResults}>
          <PowerSearchEquipmentResultsTable
            equipment={props.equipment}
            selectedEquipment={activeEquipement}
            onRowSelected={this.handleElementSelected}
          />
        </div>
      );
    };
    switch (this.state.activeStep) {
      case 0:
        return (
          <div className={classes.searchResults}>
            <EquipmentComparisonViewQueryRenderer limit={50}>
              {props => <EquipmentTable {...props} />}
            </EquipmentComparisonViewQueryRenderer>
          </div>
        );
      case 1:
        return (
          <InventoryQueryRenderer
            query={portsConnectDialogQuery}
            variables={{
              equipmentId: nullthrows(activeEquipement).id,
            }}
            render={props => {
              const {equipment} = props;
              const availablePorts = equipment.descendentsIncludingSelf
                .map(a => a.ports)
                .flat()
                .filter(p => p.id != this.props.port.id);
              return (
                <AvailablePortsTable
                  equipment={nullthrows(activeEquipement)}
                  ports={availablePorts}
                  selectedPort={targetPort}
                  onPortSelected={this.handlePortSelected}
                />
              );
            }}
          />
        );
      case 2:
        if (linkProperties.length == 0) {
          this.setState(state => ({
            activeStep: state.activeStep + 1,
          }));
        }
        return (
          <PropertiesAddEditSection
            properties={linkProperties}
            onChange={index => this._propertyChangedHandler(index)}
            nodeInput={NodePropertyInput}
          />
        );
      case 3:
        return (
          <PortsConnectConfirmation
            aSideEquipment={this.props.equipment}
            aSidePort={this.props.port}
            zSideEquipment={nullthrows(activeEquipement)}
            zSidePort={nullthrows(targetPort)}
          />
        );
      default:
        return 'Unknown step';
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

  render() {
    const {classes} = this.props;
    const {
      activeStep,
      activeEquipement,
      targetPort,
      linkProperties,
    } = this.state;
    const lastStep = activeStep == steps.length - 1;
    const connector = (
      <StepConnector
        classes={{
          active: classes.connectorActive,
          completed: classes.connectorCompleted,
          disabled: classes.connectorDisabled,
          line: classes.connectorLine,
        }}
      />
    );

    return (
      <>
        <DialogContent div className={classes.root}>
          <WizardContextProvider>
            <Stepper
              alternativeLabel
              activeStep={activeStep}
              connector={connector}>
              {steps.map(label => (
                <Step key={label}>
                  <StepLabel>{label}</StepLabel>
                </Step>
              ))}
            </Stepper>
            <div className={classes.content}>{this.getStepContent()}</div>
          </WizardContextProvider>
        </DialogContent>
        <DialogActions>
          <Button
            disabled={activeStep === 0}
            skin="gray"
            onClick={this.handleBack}>
            Back
          </Button>
          {!lastStep && (
            <Button
              disabled={
                (activeStep === 0 && !activeEquipement) ||
                (activeStep === 1 && !targetPort)
              }
              onClick={this.handleNext}>
              Next
            </Button>
          )}
          {lastStep && (
            <Button
              disabled={activeStep < steps.length - 1}
              color="primary"
              onClick={() =>
                this.props.onConnectPorts(
                  nullthrows(targetPort),
                  linkProperties,
                )
              }>
              Connect
            </Button>
          )}
        </DialogActions>
      </>
    );
  }
}

export default withStyles(styles)(PortsConnectDialog);
