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
import type {Equipment, EquipmentPort} from '../../common/Equipment';
import type {PowerSearchEquipmentResultsTable_equipment} from '../comparison_view/__generated__/PowerSearchEquipmentResultsTable_equipment.graphql';
import type {Property} from '../../common/Property';
import type {WithStyles} from '@material-ui/core';

import AvailablePortsTable from './AvailablePortsTable';
import Button from '@fbcnms/ui/components/design-system/Button';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import EquipmentComparisonViewQueryRenderer from '../comparison_view/EquipmentComparisonViewQueryRenderer';
import InventoryQueryRenderer from '../InventoryQueryRenderer';
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
import {getInitialPropertyFromType} from '../../common/PropertyType';
import {graphql} from 'react-relay';
import {sortPropertiesByIndex} from '../../common/Property';
import {uniqBy} from 'lodash';
import {withRouter} from 'react-router-dom';
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
} & WithStyles<typeof styles> &
  ContextRouter;

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
            type
            bandwidth
          }
        }
        positions {
          ...EquipmentPortsTable_position @relay(mask: false)
        }
        ports {
          ...EquipmentPortsTable_port @relay(mask: false)
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
    this.setState(state => ({
      activeStep: state.activeStep + 1,
      activeEquipement: equipment,
    }));
  };

  handlePortSelected = (port: EquipmentPort) => {
    this.setState(state => ({
      activeStep: state.activeStep + 1,
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
    const {history, classes} = this.props;
    const {linkProperties} = this.state;
    const EquipmentTable = (props: {
      equipment: PowerSearchEquipmentResultsTable_equipment,
    }) => {
      return (
        <div className={classes.searchResults}>
          <PowerSearchEquipmentResultsTable
            equipment={props.equipment}
            onEquipmentSelected={this.handleElementSelected}
            onWorkOrderSelected={workOrderId =>
              history.replace(`inventory?workorder=${workOrderId}`)
            }
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
              equipmentId: nullthrows(this.state.activeEquipement).id,
            }}
            render={props => {
              const {equipment} = props;
              return (
                <AvailablePortsTable
                  equipment={equipment}
                  onPortClicked={this.handlePortSelected}
                  sourcePortId={this.props.port.id}
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
          />
        );
      case 3:
        return (
          <PortsConnectConfirmation
            aSideEquipment={this.props.equipment}
            aSidePort={this.props.port}
            zSideEquipment={nullthrows(this.state.activeEquipement)}
            zSidePort={nullthrows(this.state.targetPort)}
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
    const {activeStep} = this.state;
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
        </DialogContent>
        <DialogActions>
          <Button
            disabled={activeStep === 0}
            skin="regular"
            onClick={this.handleBack}>
            Back
          </Button>
          {!lastStep && (
            <Button disabled={activeStep != 2} onClick={this.handleNext}>
              Next
            </Button>
          )}
          {lastStep && (
            <Button
              disabled={activeStep < steps.length - 1}
              color="primary"
              onClick={() =>
                this.props.onConnectPorts(
                  nullthrows(this.state.targetPort),
                  this.state.linkProperties,
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

export default withRouter(withStyles(styles)(PortsConnectDialog));
