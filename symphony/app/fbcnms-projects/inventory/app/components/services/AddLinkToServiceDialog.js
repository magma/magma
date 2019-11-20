/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Equipment, Link} from '../../common/Equipment';
import type {Service} from '../../common/Service';
import type {WithStyles} from '@material-ui/core';

import AvailableLinksTable from './AvailableLinksTable';
import Button from '@fbcnms/ui/components/design-system/Button';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import EquipmentComparisonViewQueryRenderer from '../comparison_view/EquipmentComparisonViewQueryRenderer';
import InventoryQueryRenderer from '../InventoryQueryRenderer';
import React from 'react';
import Step from '@material-ui/core/Step';
import StepConnector from '@material-ui/core/StepConnector/StepConnector.js';
import StepLabel from '@material-ui/core/StepLabel';
import Stepper from '@material-ui/core/Stepper';
import nullthrows from '@fbcnms/util/nullthrows';
import {graphql} from 'react-relay';
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
  service: Service,
  onClose: () => void,
  onAddLink: (link: Link) => void,
} & WithStyles<typeof styles>;

type State = {
  activeEquipement: ?Equipment,
  activeStep: number,
  activeLink: ?Link,
};

const steps = ['Select Equipment', 'Select Link'];

const addLinkToServiceDialogQuery = graphql`
  query AddLinkToServiceDialogQuery($filters: [LinkFilterInput!]!) {
    linkSearch(filters: $filters, limit: 50) {
      links {
        id
        ports {
          parentEquipment {
            id
            name
          }
          definition {
            id
            name
            type
          }
        }
        ...AvailableLinksTable_links
      }
    }
  }
`;

class AddLinkToServiceDialog extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {
      activeEquipement: null,
      activeStep: 0,
      activeLink: null,
    };
  }

  handleElementSelected = equipment => {
    this.setState(state => ({
      activeStep: state.activeStep + 1,
      activeEquipement: equipment,
    }));
  };

  handleLinkSelected = (link: Link) => {
    const {onAddLink} = this.props;
    this.setState(_ => ({
      activeLink: link,
    }));
    onAddLink(nullthrows(link));
  };

  getStepContent = () => {
    const {classes, service} = this.props;
    switch (this.state.activeStep) {
      case 0:
        return (
          <div className={classes.searchResults}>
            <EquipmentComparisonViewQueryRenderer
              limit={50}
              onEquipmentSelected={this.handleElementSelected}
            />
          </div>
        );
      case 1:
        return (
          <InventoryQueryRenderer
            query={addLinkToServiceDialogQuery}
            variables={{
              filters: [
                {
                  filterType: 'SERVICE_INST',
                  operator: 'IS_NOT_ONE_OF',
                  idSet: [this.props.service.id],
                },
                {
                  filterType: 'EQUIPMENT_INST',
                  operator: 'IS_ONE_OF',
                  idSet: [nullthrows(this.state.activeEquipement).id],
                },
              ],
            }}
            render={props => {
              const {linkSearch} = props;
              // TODO: Remove this filtering after commiting links change
              //       on every link add
              const links = linkSearch.links.filter(
                link => !(service.links.map(l => l.id) || []).includes(link.id),
              );
              return (
                <AvailableLinksTable
                  equipment={nullthrows(this.state.activeEquipement)}
                  links={links}
                  onLinkSelected={this.handleLinkSelected}
                />
              );
            }}
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
    const {classes, onAddLink} = this.props;
    const {activeStep, activeLink} = this.state;
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
            <Button disabled={activeStep != 1} onClick={this.handleNext}>
              Next
            </Button>
          )}
          {lastStep && (
            <Button
              disabled={activeStep < steps.length - 1}
              color="primary"
              onClick={() => onAddLink(nullthrows(activeLink))}>
              Add
            </Button>
          )}
        </DialogActions>
      </>
    );
  }
}

export default withStyles(styles)(AddLinkToServiceDialog);
