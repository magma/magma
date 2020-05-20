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
import type {
  EditLocationTypesIndexMutationResponse,
  EditLocationTypesIndexMutationVariables,
} from '../../mutations/__generated__/EditLocationTypesIndexMutation.graphql';
import type {LocationTypeItem_locationType} from './__generated__/LocationTypeItem_locationType.graphql';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithSnackbarProps} from 'notistack';
import type {WithStyles} from '@material-ui/core';

import AddEditLocationTypeCard from './AddEditLocationTypeCard';
import Button from '@fbcnms/ui/components/design-system/Button';
import CircularProgress from '@material-ui/core/CircularProgress';
import ConfigueTitle from '@fbcnms/ui/components/ConfigureTitle';
import DroppableTableBody from '../draggable/DroppableTableBody';
import EditLocationTypesIndexMutation from '../../mutations/EditLocationTypesIndexMutation';
import FormActionWithPermissions from '../../common/FormActionWithPermissions';
import InventoryQueryRenderer from '../InventoryQueryRenderer';
import LocationTypeItem from './LocationTypeItem';
import React from 'react';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import withInventoryErrorBoundary from '../../common/withInventoryErrorBoundary';
import {FormContextProvider} from '../../common/FormContext';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {getGraphError} from '../../common/EntUtils';
import {graphql} from 'relay-runtime';
import {reorder, sortByIndex} from '../draggable/DraggableUtils';
import {withRouter} from 'react-router-dom';
import {withSnackbar} from 'notistack';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  root: {
    display: 'flex',
    width: '100%',
    flexDirection: 'column',
  },
  table: {
    width: '100%',
    marginTop: '15px',
  },
  paper: {
    flexGrow: 1,
    overflowY: 'hidden',
  },
  typesList: {
    padding: '24px',
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
  progress: {
    alignSelf: 'center',
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

type Props = ContextRouter &
  WithStyles<typeof styles> &
  WithSnackbarProps &
  WithAlert & {};

type State = {
  dialogKey: number,
  errorMessage: ?string,
  showAddEditCard: boolean,
  editingLocationType: ?LocationTypeItem_locationType,
  isSaving: boolean,
};

const locationTypesQuery = graphql`
  query LocationTypesQuery {
    locationTypes(first: 500) @connection(key: "Catalog_locationTypes") {
      edges {
        node {
          ...LocationTypeItem_locationType
          ...AddEditLocationTypeCard_editingLocationType
          id
          name
          index
        }
      }
    }
  }
`;

class LocationTypes extends React.Component<Props, State> {
  state = {
    dialogKey: 1,
    errorMessage: null,
    showAddEditCard: false,
    editingLocationType: null,
    isSaving: false,
  };

  render() {
    const {classes} = this.props;
    const {showAddEditCard, editingLocationType} = this.state;

    return (
      <InventoryQueryRenderer
        query={locationTypesQuery}
        variables={{}}
        render={props => {
          const {locationTypes} = props;
          if (showAddEditCard) {
            return (
              <div className={classes.paper}>
                <AddEditLocationTypeCard
                  key={'new_location_type@' + this.state.dialogKey}
                  open={showAddEditCard}
                  onClose={this.hideAddEditLocationTypeCard}
                  onSave={this.saveLocation}
                  editingLocationType={editingLocationType}
                />
              </div>
            );
          }
          return (
            <FormContextProvider
              permissions={{
                entity: 'locationType',
              }}>
              <div className={classes.typesList}>
                <div className={classes.firstRow}>
                  <ConfigueTitle
                    className={classes.title}
                    title={'Location Types'}
                    subtitle={
                      'Drag and drop location types to arrange them by size, from largest to smallest'
                    }
                  />
                  <div className={classes.addButtonContainer}>
                    {this.state.isSaving ? (
                      <CircularProgress className={classes.progress} />
                    ) : null}
                    <FormActionWithPermissions
                      permissions={{entity: 'locationType', action: 'create'}}>
                      <Button
                        className={classes.addButton}
                        onClick={() => this.showAddEditLocationTypeCard(null)}>
                        Add Location Type
                      </Button>
                    </FormActionWithPermissions>
                  </div>
                </div>
                <div className={classes.root}>
                  <DroppableTableBody
                    isDragDisabled={this.state.isSaving}
                    className={classes.table}
                    onDragEnd={res => this._onDragEnd(res, locationTypes)}>
                    {locationTypes.edges
                      .map(edge => edge.node)
                      .sort(sortByIndex)
                      .map((locType, i) => {
                        return (
                          <div
                            className={classes.listItem}
                            key={`${locType.id}_${locType.index}`}>
                            <LocationTypeItem
                              locationType={locType}
                              position={i}
                              onEdit={() =>
                                this.showAddEditLocationTypeCard(locType)
                              }
                            />
                          </div>
                        );
                      })}
                  </DroppableTableBody>
                </div>
              </div>
            </FormContextProvider>
          );
        }}
      />
    );
  }

  _onDragEnd = (result, locationTypes) => {
    if (!result.destination) {
      return;
    }
    locationTypes = locationTypes.edges
      .map(edge => edge.node)
      .sort(sortByIndex);

    ServerLogger.info(LogEvents.LOCATION_TYPE_REORDERED);
    const items = reorder(
      locationTypes,
      result.source.index,
      result.destination.index,
    );
    const newItems = items.map((locTyp, i) => ({...locTyp, index: i}));
    this.saveOrder(newItems);
  };

  saveOrder = newItems => {
    const variables: EditLocationTypesIndexMutationVariables = {
      locationTypeIndex: this.buildMutationVariables(newItems),
    };
    this.setState({isSaving: true});
    // eslint-disable-next-line max-len
    const callbacks: MutationCallbacks<EditLocationTypesIndexMutationResponse> = {
      onCompleted: (response, errors) => {
        this.setState({isSaving: false});
        if (errors && errors[0]) {
          this.props.enqueueSnackbar(errors[0].message, {
            children: key => (
              <SnackbarItem
                id={key}
                message={errors[0].message}
                variant="error"
              />
            ),
          });
        } else {
          this.setState({isSaving: false});
        }
      },
      onError: (error: Error) => {
        this.setState({errorMessage: getGraphError(error), isSaving: false});
      },
    };
    EditLocationTypesIndexMutation(variables, callbacks);
  };

  showAddEditLocationTypeCard = (locType: ?LocationTypeItem_locationType) => {
    ServerLogger.info(LogEvents.ADD_LOCATION_TYPE_BUTTON_CLICKED);
    this.setState({editingLocationType: locType, showAddEditCard: true});
  };

  hideAddEditLocationTypeCard = () =>
    this.setState(prevState => ({
      editingLocationType: null,
      showAddEditCard: false,
      dialogKey: prevState.dialogKey + 1,
    }));

  saveLocation = (locationType: LocationTypeItem_locationType) => {
    ServerLogger.info(LogEvents.SAVE_LOCATION_TYPE_BUTTON_CLICKED);
    this.setState(prevState => {
      if (locationType) {
        return {
          dialogKey: prevState.dialogKey + 1,
          showAddEditCard: false,
        };
      }
    });
  };

  buildMutationVariables = newItems => {
    return newItems.map(item => {
      return {
        locationTypeID: item.id,
        index: item.index,
      };
    });
  };
}

export default withStyles(styles)(
  withAlert(
    withSnackbar(withRouter(withInventoryErrorBoundary(LocationTypes))),
  ),
);
