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
import type {ServiceTypeItem_serviceType} from './__generated__/ServiceTypeItem_serviceType.graphql';
import type {WithStyles} from '@material-ui/core';

import AddEditServiceTypeCard from './AddEditServiceTypeCard';
import Button from '@fbcnms/ui/components/design-system/Button';
import ConfigueTitle from '@fbcnms/ui/components/ConfigureTitle';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import InventoryQueryRenderer from '../InventoryQueryRenderer';
import React from 'react';
import ServiceTypeItem from './ServiceTypeItem';
import withInventoryErrorBoundary from '../../common/withInventoryErrorBoundary';
import {FormContextProvider} from '../../common/FormContext';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {graphql} from 'relay-runtime';
import {sortLexicographically} from '@fbcnms/ui/utils/displayUtils';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  root: {
    width: '100%',
    marginTop: '15px',
  },
  paper: {
    flexGrow: 1,
    overflowY: 'hidden',
  },
  content: {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'flex-start',
  },
  listItem: {
    marginBottom: theme.spacing(),
  },
  addButton: {
    marginLeft: 'auto',
  },
  addButtonContainer: {
    display: 'flex',
  },
  typesList: {
    padding: '24px',
  },
  title: {
    marginLeft: '10px',
  },
  firstRow: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
});

type Props = ContextRouter & WithStyles<typeof styles> & {};

type State = {
  dialogKey: number,
  errorMessage: ?string,
  editingServiceType: ?ServiceTypeItem_serviceType,
  showAddEditCard: boolean,
};

const serviceTypesQuery = graphql`
  query ServiceTypesQuery {
    serviceTypes(first: 500) @connection(key: "ServiceTypes_serviceTypes") {
      edges {
        node {
          ...ServiceTypeItem_serviceType
          ...AddEditServiceTypeCard_editingServiceType
          id
          name
          isDeleted
        }
      }
    }
  }
`;

class ServiceTypes extends React.Component<Props, State> {
  state = {
    dialogKey: 1,
    errorMessage: null,
    showAddEditCard: false,
    editingServiceType: null,
  };

  render() {
    const {classes} = this.props;
    const {showAddEditCard, editingServiceType} = this.state;

    return (
      <InventoryQueryRenderer
        query={serviceTypesQuery}
        variables={{}}
        render={props => {
          const {serviceTypes} = props;

          if (showAddEditCard) {
            return (
              <div className={classes.paper}>
                <AddEditServiceTypeCard
                  key={'new_service_type@' + this.state.dialogKey}
                  open={showAddEditCard}
                  onClose={this.hideAddEditServiceTypeCard}
                  onSave={this.saveService}
                  editingServiceType={editingServiceType}
                />
              </div>
            );
          }

          return (
            <FormContextProvider>
              <div className={classes.typesList}>
                <div className={classes.firstRow}>
                  <ConfigueTitle
                    className={classes.title}
                    title={'Service Types'}
                    subtitle={'Manage the types of services in your inventory'}
                  />
                  <div className={classes.addButtonContainer}>
                    <FormAction>
                      <Button
                        className={classes.addButton}
                        onClick={() => this.showAddEditServiceTypeCard(null)}>
                        Add Service Type
                      </Button>
                    </FormAction>
                  </div>
                </div>
                <div className={classes.root}>
                  {serviceTypes.edges
                    .map(edge => edge.node)
                    .filter(Boolean)
                    .sort((serviceTypeA, serviceTypeB) =>
                      sortLexicographically(
                        serviceTypeA.name,
                        serviceTypeB.name,
                      ),
                    )
                    .filter(s => !s.isDeleted)
                    .map(srvType => (
                      <div
                        className={classes.listItem}
                        key={`srvType_${srvType.id}`}>
                        <ServiceTypeItem
                          serviceType={srvType}
                          onEdit={() =>
                            this.showAddEditServiceTypeCard(srvType)
                          }
                        />
                      </div>
                    ))}
                </div>
              </div>
            </FormContextProvider>
          );
        }}
      />
    );
  }

  showAddEditServiceTypeCard = (serviceType: ?ServiceTypeItem_serviceType) => {
    ServerLogger.info(LogEvents.ADD_SERVICE_TYPE_BUTTON_CLICKED);
    this.setState({
      showAddEditCard: true,
      editingServiceType: serviceType,
    });
  };

  hideAddEditServiceTypeCard = () =>
    this.setState(state => ({
      editingServiceType: null,
      dialogKey: state.dialogKey + 1,
      showAddEditCard: false,
    }));

  saveService = (serviceType: ServiceTypeItem_serviceType) => {
    ServerLogger.info(LogEvents.SAVE_SERVICE_TYPE_BUTTON_CLICKED);
    this.setState(state => {
      if (serviceType) {
        return {
          showAddEditCard: false,
          dialogKey: state.dialogKey + 1,
          editingServiceType: null,
        };
      }
    });
  };
}

export default withStyles(styles)(
  withRouter(withInventoryErrorBoundary(ServiceTypes)),
);
