/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  AddWorkOrderMutationResponse,
  AddWorkOrderMutationVariables,
  CheckListCategoryInput,
} from '../../mutations/__generated__/AddWorkOrderMutation.graphql';
import type {ChecklistCategoriesMutateStateActionType} from '../checklist/ChecklistCategoriesMutateAction';
import type {ChecklistCategoriesStateType} from '../checklist/ChecklistCategoriesMutateState';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {WorkOrder, WorkOrderType} from '../../common/WorkOrder';

import AddWorkOrderMutation from '../../mutations/AddWorkOrderMutation';
import Breadcrumbs from '@fbcnms/ui/components/Breadcrumbs';
import CheckListCategoryExpandingPanel from '../checklist/checkListCategory/CheckListCategoryExpandingPanel';
import ChecklistCategoriesMutateDispatchContext from '../checklist/ChecklistCategoriesMutateDispatchContext';
import CircularProgress from '@material-ui/core/CircularProgress';
import ExpandingPanel from '@fbcnms/ui/components/ExpandingPanel';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import FormSaveCancelPanel from '@fbcnms/ui/components/design-system/Form/FormSaveCancelPanel';
import Grid from '@material-ui/core/Grid';
import InventoryQueryRenderer from '../InventoryQueryRenderer';
import LocationTypeahead from '../typeahead/LocationTypeahead';
import MenuItem from '@material-ui/core/MenuItem';
import NameDescriptionSection from '@fbcnms/ui/components/NameDescriptionSection';
import ProjectTypeahead from '../typeahead/ProjectTypeahead';
import PropertyValueInput from '../form/PropertyValueInput';
import React, {useCallback, useReducer, useState} from 'react';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import TextField from '@material-ui/core/TextField';
import UserTypeahead from '../typeahead/UserTypeahead';
import nullthrows from '@fbcnms/util/nullthrows';
import {FormValidationContextProvider} from '@fbcnms/ui/components/design-system/Form/FormValidationContext';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {getInitialPropertyFromType} from '../../common/PropertyType';
import {graphql} from 'relay-runtime';
import {makeStyles} from '@material-ui/styles';
import {priorityValues, statusValues} from '../../common/WorkOrder';
import {reducer} from '../checklist/ChecklistCategoriesMutateReducer';
import {removeTempIDs} from '../../common/EntUtils';
import {sortPropertiesByIndex, toPropertyInput} from '../../common/Property';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useHistory, useRouteMatch} from 'react-router';

type Props = {
  workOrderTypeId: ?string,
};

const useStyles = makeStyles(theme => ({
  root: {
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
    padding: '40px 32px',
  },
  contentRoot: {
    display: 'flex',
    flexDirection: 'column',
    position: 'relative',
    flexGrow: 1,
    overflow: 'auto',
  },
  cards: {
    flexGrow: 1,
    overflow: 'hidden',
    overflowY: 'auto',
  },
  card: {
    display: 'flex',
    flexDirection: 'column',
  },
  input: {
    width: '250px',
    paddingBottom: '24px',
  },
  gridInput: {
    display: 'inline-flex',
  },
  nameHeader: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: '24px',
  },
  breadcrumbs: {
    flexGrow: 1,
  },
  separator: {
    borderBottom: `1px solid ${theme.palette.grey[100]}`,
    margin: '0 0 24px -24px',
    paddingBottom: '24px',
    width: 'calc(100% + 48px)',
  },
  separator: {
    borderBottom: `1px solid ${theme.palette.grey[100]}`,
    margin: '0 0 24px -24px',
    paddingBottom: '24px',
    width: 'calc(100% + 48px)',
  },
  dense: {
    paddingTop: '9px',
    paddingBottom: '9px',
    height: '14px',
  },
  cancelButton: {
    marginRight: '8px',
  },
}));

const addWorkOrderCard__workOrderTypeQuery = graphql`
  query AddWorkOrderCard__workOrderTypeQuery($workOrderTypeId: ID!) {
    workOrderType: node(id: $workOrderTypeId) {
      ... on WorkOrderType {
        id
        name
        description
        propertyTypes {
          id
          name
          type
          index
          stringValue
          intValue
          booleanValue
          floatValue
          latitudeValue
          longitudeValue
          rangeFromValue
          rangeToValue
          isEditable
          isMandatory
          isInstanceProperty
          isDeleted
        }
      }
    }
  }
`;

const AddWorkOrderCard = ({workOrderTypeId}: Props) => {
  const classes = useStyles();
  const [workOrder, setWorkOrder] = useState<?WorkOrder>(null);
  const enqueueSnackbar = useEnqueueSnackbar();
  const history = useHistory();
  const match = useRouteMatch();

  const [editingCategories, dispatch] = useReducer<
    ChecklistCategoriesStateType,
    ChecklistCategoriesMutateStateActionType,
  >(reducer, []);

  const _enqueueError = useCallback(
    (message: string) => {
      enqueueSnackbar(message, {
        children: key => (
          <SnackbarItem id={key} message={message} variant="error" />
        ),
      });
    },
    [enqueueSnackbar],
  );

  const _creaetNewWorkOrder = (workOrderType: WorkOrderType): WorkOrder => {
    const initialProps = (workOrderType.propertyTypes || [])
      .filter(propertyType => !propertyType.isDeleted)
      .map(propType => getInitialPropertyFromType(propType))
      .sort(sortPropertiesByIndex);
    return {
      id: 'workOrder@tmp',
      workOrderType: workOrderType,
      workOrderTypeId: workOrderType.id,
      name: workOrderType.name,
      description: workOrderType.description,
      locationId: null,
      location: null,
      properties: initialProps,
      workOrders: [],
      ownerName: '',
      creationDate: '',
      installDate: '',
      status: 'PENDING',
      priority: 'NONE',
      equipmentToAdd: [],
      equipmentToRemove: [],
      linksToAdd: [],
      linksToRemove: [],
      files: [],
      images: [],
      assignee: '',
      projectId: null,
      checkListCategories: [],
    };
  };

  const _saveWorkOrder = () => {
    const {
      name,
      description,
      locationId,
      projectId,
      assignee,
      status,
      priority,
      properties,
      checkListCategories,
    } = nullthrows(workOrder);
    const workOrderTypeId = nullthrows(workOrder?.workOrderTypeId);
    const updatedChecklistCategories: CheckListCategoryInput[] = (
      checkListCategories || []
    ).map(category => ({
      title: category.title,
      description: category.description,
      checkList: removeTempIDs(category.checkList || []),
    }));
    const variables: AddWorkOrderMutationVariables = {
      input: {
        name,
        description,
        locationId,
        workOrderTypeId,
        assignee,
        projectId,
        status,
        priority,
        properties: toPropertyInput(properties),
        checkListCategories: updatedChecklistCategories,
      },
    };

    const callbacks: MutationCallbacks<AddWorkOrderMutationResponse> = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          _enqueueError(errors[0].message);
        } else {
          // navigate to main page
          history.push(match.url);
        }
      },
      onError: () => {
        _enqueueError('Error saving work order');
      },
    };
    ServerLogger.info(LogEvents.SAVE_PROJECT_BUTTON_CLICKED, {
      source: 'workOrder_details',
    });
    AddWorkOrderMutation(variables, callbacks);
  };

  const _setWorkOrderDetail = (
    key:
      | 'name'
      | 'description'
      | 'assignee'
      | 'projectId'
      | 'locationId'
      | 'priority'
      | 'status',
    value,
  ) => {
    setWorkOrder(prevWorkOrder => {
      if (!prevWorkOrder) {
        return;
      }
      return {...prevWorkOrder, [`${key}`]: value};
    });
  };

  const _propertyChangedHandler = index => property =>
    // eslint-disable-next-line no-warning-comments
    // $FlowFixMe - known techdebt with Property/PropertyType flow definitions
    setWorkOrder(prevWorkOrder => {
      if (!prevWorkOrder) {
        return;
      }
      return {
        ...prevWorkOrder,
        properties: [
          ...prevWorkOrder.properties.slice(0, index),
          // eslint-disable-next-line no-warning-comments
          // $FlowFixMe - known techdebt with Property/PropertyType flow definitions
          property,
          ...prevWorkOrder.properties.slice(index + 1),
        ],
      };
    });

  const _checkListCategoryChangedHandler = updatedCategories => {
    setWorkOrder(prevWorkOrder => {
      if (!prevWorkOrder) {
        return;
      }
      return {
        ...prevWorkOrder,
        checkListCategories: updatedCategories,
      };
    });
  };

  const navigateToMainPage = () => {
    ServerLogger.info(LogEvents.WORK_ORDERS_SEARCH_NAV_CLICKED, {
      source: 'work_order_details',
    });
    history.push(match.url);
  };

  return (
    <InventoryQueryRenderer
      query={addWorkOrderCard__workOrderTypeQuery}
      variables={{
        workOrderTypeId,
      }}
      render={queryData => {
        const {workOrderType} = queryData;
        if (!workOrder && workOrderType) {
          setWorkOrder(_creaetNewWorkOrder(workOrderType));
        }
        if (!workOrder) {
          return (
            <div className={classes.root}>
              <CircularProgress />
            </div>
          );
        }
        return (
          <div className={classes.root}>
            <FormValidationContextProvider>
              <div className={classes.nameHeader}>
                <Breadcrumbs
                  className={classes.breadcrumbs}
                  breadcrumbs={[
                    {
                      id: 'workOrders',
                      name: 'WorkOrders',
                      onClick: () => navigateToMainPage(),
                    },
                    {
                      id: `new_workOrder_` + Date.now(),
                      name: 'New WorkOrder',
                    },
                  ]}
                  size="large"
                />
                <FormSaveCancelPanel
                  onCancel={navigateToMainPage}
                  onSave={_saveWorkOrder}
                />
              </div>
              <div className={classes.contentRoot}>
                <div className={classes.cards}>
                  <Grid container spacing={2}>
                    <Grid item xs={8} sm={8} lg={8} xl={8}>
                      <ExpandingPanel title="Details">
                        <NameDescriptionSection
                          name={workOrder.name}
                          description={workOrder.description}
                          onNameChange={value =>
                            _setWorkOrderDetail('name', value)
                          }
                          onDescriptionChange={value =>
                            _setWorkOrderDetail('description', value)
                          }
                        />
                        <div className={classes.separator} />
                        <Grid container spacing={2}>
                          <Grid item xs={12} sm={6} lg={4} xl={4}>
                            <FormField label="Project">
                              <ProjectTypeahead
                                className={classes.gridInput}
                                margin="dense"
                                onProjectSelection={project =>
                                  _setWorkOrderDetail('projectId', project?.id)
                                }
                              />
                            </FormField>
                          </Grid>
                          {workOrder.workOrderType && (
                            <Grid item xs={12} sm={6} lg={4} xl={4}>
                              <FormField label="Type">
                                <TextField
                                  disabled
                                  variant="outlined"
                                  margin="dense"
                                  className={classes.gridInput}
                                  value={workOrder.workOrderType.name}
                                />
                              </FormField>
                            </Grid>
                          )}
                          <Grid item xs={12} sm={6} lg={4} xl={4}>
                            <FormField label="Priority">
                              <TextField
                                select
                                className={classes.gridInput}
                                variant="outlined"
                                value={workOrder.priority}
                                InputProps={{
                                  classes: {
                                    input: classes.dense,
                                  },
                                }}
                                onChange={event => {
                                  _setWorkOrderDetail(
                                    'priority',
                                    event.target.value,
                                  );
                                }}>
                                {priorityValues.map(option => (
                                  <MenuItem
                                    key={option.value}
                                    value={option.value}>
                                    {option.label}
                                  </MenuItem>
                                ))}
                              </TextField>
                            </FormField>
                          </Grid>
                          <Grid item xs={12} sm={6} lg={4} xl={4}>
                            <FormField label="Status">
                              <TextField
                                select
                                className={classes.gridInput}
                                variant="outlined"
                                value={workOrder.status}
                                InputProps={{
                                  classes: {
                                    input: classes.dense,
                                  },
                                }}
                                onChange={event => {
                                  _setWorkOrderDetail(
                                    'status',
                                    event.target.value,
                                  );
                                }}>
                                {statusValues.map(option => (
                                  <MenuItem
                                    key={option.value}
                                    value={option.value}>
                                    {option.label}
                                  </MenuItem>
                                ))}
                              </TextField>
                            </FormField>
                          </Grid>
                          <Grid item xs={12} sm={6} lg={4} xl={4}>
                            <FormField label="Location">
                              <LocationTypeahead
                                headline={null}
                                className={classes.gridInput}
                                margin="dense"
                                onLocationSelection={location =>
                                  _setWorkOrderDetail(
                                    'locationId',
                                    location?.id ?? null,
                                  )
                                }
                              />
                            </FormField>
                          </Grid>
                          {workOrder.properties
                            .filter(
                              property => !property.propertyType.isDeleted,
                            )
                            .map((property, index) => (
                              <Grid
                                key={property.id}
                                item
                                xs={12}
                                sm={6}
                                lg={4}
                                xl={4}>
                                <PropertyValueInput
                                  required={!!property.propertyType.isMandatory}
                                  disabled={
                                    !property.propertyType.isInstanceProperty
                                  }
                                  label={property.propertyType.name}
                                  className={classes.gridInput}
                                  margin="dense"
                                  inputType="Property"
                                  property={property}
                                  headlineVariant="form"
                                  fullWidth={true}
                                  onChange={_propertyChangedHandler(index)}
                                />
                              </Grid>
                            ))}
                        </Grid>
                      </ExpandingPanel>
                      <ChecklistCategoriesMutateDispatchContext.Provider
                        value={dispatch}>
                        <CheckListCategoryExpandingPanel
                          categories={editingCategories}
                          onListChanged={_checkListCategoryChangedHandler}
                        />
                      </ChecklistCategoriesMutateDispatchContext.Provider>
                    </Grid>
                    <Grid item xs={4} sm={4} lg={4} xl={4}>
                      <ExpandingPanel title="Team">
                        <UserTypeahead
                          className={classes.input}
                          headline="Assignee"
                          onUserSelection={user =>
                            _setWorkOrderDetail('assignee', user)
                          }
                          margin="dense"
                        />
                      </ExpandingPanel>
                    </Grid>
                  </Grid>
                </div>
              </div>
            </FormValidationContextProvider>
          </div>
        );
      }}
    />
  );
};

export default AddWorkOrderCard;
