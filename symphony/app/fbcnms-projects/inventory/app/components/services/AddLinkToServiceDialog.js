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
import type {PowerSearchLinkFirstEquipmentResultsTable_equipment} from './__generated__/PowerSearchLinkFirstEquipmentResultsTable_equipment.graphql';
import type {WithStyles} from '@material-ui/core';

import AvailableLinksTable from './AvailableLinksTable';
import Button from '@fbcnms/ui/components/design-system/Button';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import EquipmentComparisonViewQueryRenderer from '../comparison_view/EquipmentComparisonViewQueryRenderer';
import InventoryQueryRenderer from '../InventoryQueryRenderer';
import PowerSearchLinkFirstEquipmentResultsTable from './PowerSearchLinkFirstEquipmentResultsTable';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import nullthrows from '@fbcnms/util/nullthrows';
import symphony from '@fbcnms/ui/theme/symphony';
import {WizardContextProvider} from '@fbcnms/ui/components/design-system/Wizard/WizardContext';
import {graphql} from 'react-relay';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  button: {
    marginTop: theme.spacing(),
    marginRight: theme.spacing(),
  },
  content: {
    height: '100%',
    width: '100%',
  },
  portIdLabel: {
    marginRight: '8px',
    fontWeight: 500,
  },
  root: {
    minWidth: '80vh',
    paddingTop: '0px',
    paddingLeft: '32px',
    paddingRight: '32px',
  },
  searchResults: {
    flexGrow: 1,
    display: 'flex',
    flexDirection: 'column',
    backgroundColor: theme.palette.common.white,
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
});

type Props = {
  service: {id: string, name: string},
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
    this.setState({
      activeEquipement: equipment,
    });
  };

  handleLinkSelected = (link: Link) => {
    this.setState({
      activeLink: link,
    });
  };

  getStepContent = () => {
    const {classes} = this.props;
    const EquipmentTable = (props: {
      equipment: PowerSearchLinkFirstEquipmentResultsTable_equipment,
    }) => {
      return (
        <div className={classes.searchResults}>
          <PowerSearchLinkFirstEquipmentResultsTable
            equipment={props.equipment}
            onEquipmentSelected={this.handleElementSelected}
            selectedEquipment={this.state.activeEquipement}
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
              return (
                <AvailableLinksTable
                  equipment={nullthrows(this.state.activeEquipement)}
                  links={linkSearch.links}
                  selectedLink={this.state.activeLink}
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
    const {classes, onAddLink, service, onClose} = this.props;
    const {activeStep, activeLink, activeEquipement} = this.state;
    const lastStep = activeStep == steps.length - 1;

    return (
      <>
        <DialogTitle>
          <Text className={classes.title} variant="h6">
            Add link to {service.name}
          </Text>
          {lastStep ? (
            <Text
              className={classes.subtitle}
              variant="subtitle2"
              color="light">
              Select the link you want to add to this service.
            </Text>
          ) : (
            <Text
              className={classes.subtitle}
              variant="subtitle2"
              color="light">
              Select the equipment associated with the link.
            </Text>
          )}
        </DialogTitle>
        <DialogContent div className={classes.root}>
          <WizardContextProvider>
            <div className={classes.content}>{this.getStepContent()}</div>
          </WizardContextProvider>
        </DialogContent>
        <DialogActions className={classes.footer}>
          {!lastStep && (
            <Button skin="gray" onClick={onClose}>
              Cancel
            </Button>
          )}
          {lastStep && (
            <Button skin="gray" onClick={this.handleBack}>
              Back
            </Button>
          )}
          {!lastStep && (
            <Button
              disabled={activeEquipement === null}
              onClick={this.handleNext}
              className={classes.actionButton}>
              Next
            </Button>
          )}
          {lastStep && (
            <Button
              disabled={activeLink === null}
              color="primary"
              onClick={() => onAddLink(nullthrows(activeLink))}
              className={classes.actionButton}>
              Add
            </Button>
          )}
        </DialogActions>
      </>
    );
  }
}

export default withRouter(withStyles(styles)(AddLinkToServiceDialog));
