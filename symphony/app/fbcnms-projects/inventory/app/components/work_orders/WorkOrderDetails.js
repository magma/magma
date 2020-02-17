/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {AddImageMutationResponse} from '../../mutations/__generated__/AddImageMutation.graphql';
import type {AddImageMutationVariables} from '../../mutations/__generated__/AddImageMutation.graphql';
import type {AppContextType} from '@fbcnms/ui/context/AppContext';
import type {
  CheckListCategoryTable_list,
  CheckListTable_list,
  WorkOrderDetails_workOrder,
} from './__generated__/WorkOrderDetails_workOrder.graphql.js';
import type {ContextRouter} from 'react-router-dom';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {Property} from '../../common/Property';
import type {Theme, WithStyles} from '@material-ui/core';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithSnackbarProps} from 'notistack';

import AddHyperlinkButton from '../AddHyperlinkButton';
import AddIcon from '@fbcnms/ui/components/design-system/Icons/Actions/AddIcon';
import AddImageMutation from '../../mutations/AddImageMutation';
import AppContext from '@fbcnms/ui/context/AppContext';
import Button from '@fbcnms/ui/components/design-system/Button';
import CheckListCategoryContext, {
  CheckListCategoryContextProvider,
} from '../checklist/checkListCategory/CheckListCategoryContext';
import CheckListCategoryTable from '../checklist/checkListCategory/CheckListCategoryTable';
import CheckListTable from '../checklist/CheckListTable';
import CircularProgress from '@material-ui/core/CircularProgress';
import CloudUploadOutlinedIcon from '@material-ui/icons/CloudUploadOutlined';
import CommentsBox from '../comments/CommentsBox';
import EditToggleButton from '@fbcnms/ui/components/design-system/toggles/EditToggleButton/EditToggleButton';
import EntityDocumentsTable from '../EntityDocumentsTable';
import ExpandingPanel from '@fbcnms/ui/components/ExpandingPanel';
import FileUpload from '../FileUpload';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import FormValidationContext, {
  FormValidationContextProvider,
} from '@fbcnms/ui/components/design-system/Form/FormValidationContext';
import Grid from '@material-ui/core/Grid';
import InsertLinkIcon from '@material-ui/icons/InsertLink';
import LocationBreadcrumbsTitle from '../location/LocationBreadcrumbsTitle';
import LocationMapSnippet from '../location/LocationMapSnippet';
import LocationTypeahead from '../typeahead/LocationTypeahead';
import NameDescriptionSection from '@fbcnms/ui/components/NameDescriptionSection';
import ProjectTypeahead from '../typeahead/ProjectTypeahead';
import PropertyValueInput from '../form/PropertyValueInput';
import React from 'react';
import Select from '@fbcnms/ui/components/design-system/Select/Select';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import Strings from '../../common/CommonStrings';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import UserTypeahead from '../typeahead/UserTypeahead';
import WorkOrderDetailsPane from './WorkOrderDetailsPane';
import WorkOrderHeader from './WorkOrderHeader';
import fbt from 'fbt';
import update from 'immutability-helper';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {createFragmentContainer, graphql} from 'react-relay';
import {doneStatus, priorityValues, statusValues} from '../../common/WorkOrder';
import {formatDateForTextInput} from '@fbcnms/ui/utils/displayUtils';
import {sortPropertiesByIndex} from '../../common/Property';
import {withRouter} from 'react-router-dom';
import {withSnackbar} from 'notistack';
import {withStyles} from '@material-ui/core/styles';

type State = {
  workOrder: WorkOrderDetails_workOrder,
  checklist: CheckListTable_list,
  checkListCategories: CheckListCategoryTable_list,
  properties: Array<Property>,
  locationId: ?string,
  isLoadingDocument: boolean,
  showChecklistDesignMode: boolean,
};

type Props = {
  workOrder: WorkOrderDetails_workOrder,
  onWorkOrderExecuted: () => void,
  onDocumentUploaded: () => void,
  onWorkOrderRemoved: () => void,
  onCancelClicked: () => void,
} & WithAlert &
  WithStyles<typeof styles> &
  WithSnackbarProps &
  ContextRouter;

const FileTypeEnum = {
  IMAGE: 'IMAGE',
  FILE: 'FILE',
};

const styles = (theme: Theme) => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
    height: '100%',
  },
  input: {
    paddingBottom: '24px',
  },
  gridInput: {
    display: 'inline-flex',
  },
  cards: {
    overflowY: 'auto',
    overflowX: 'hidden',
    flexGrow: 1,
    flexBasis: 0,
  },
  card: {
    display: 'flex',
    flexDirection: 'column',
  },
  separator: {
    borderBottom: `1px solid ${theme.palette.grey[100]}`,
    margin: '0 0 16px -24px',
    paddingBottom: '24px',
    width: 'calc(100% + 48px)',
  },
  uploadButtonContainer: {
    display: 'flex',
    marginRight: '8px',
    marginTop: '4px',
  },
  uploadButton: {
    cursor: 'pointer',
    fill: theme.palette.primary.main,
  },
  minimizedButton: {
    minWidth: 'unset',
  },
  dense: {
    paddingTop: '9px',
    paddingBottom: '9px',
    height: '14px',
  },
  breadcrumbs: {
    marginBottom: '16px',
  },
  propertiesGrid: {
    marginTop: '16px',
  },
  commentsBoxContainer: {
    padding: '0px',
  },
  inExpandingPanelFix: {
    paddingLeft: '24px',
    paddingRight: '24px',
  },
  commentsLog: {
    maxHeight: '400px',
  },
  map: {
    minHeight: '232px',
  },
});

class WorkOrderDetails extends React.Component<Props, State> {
  state = {
    workOrder: this.props.workOrder,
    properties: this.getEditingProperties(),
    checklist: this.props.workOrder.checkList,
    checkListCategories: this.props.workOrder.checkListCategories,
    locationId: this.props.workOrder.location?.id,
    isLoadingDocument: false,
    showChecklistDesignMode: false,
  };

  getEditingProperties(): Array<Property> {
    // eslint-disable-next-line flowtype/no-weak-types
    return ([...this.props.workOrder.properties]: any).sort(
      sortPropertiesByIndex,
    );
  }

  static contextType = AppContext;
  context: AppContextType;

  render() {
    const {classes, onWorkOrderRemoved, onCancelClicked} = this.props;
    const {
      workOrder,
      properties,
      checklist,
      checkListCategories,
      locationId,
      showChecklistDesignMode,
    } = this.state;
    const {location} = workOrder;
    const hasCheckListCategories = checkListCategories.length > 0;
    const actionsEnabled = this.context.isFeatureEnabled('planned_equipment');
    return (
      <div className={classes.root}>
        <FormValidationContextProvider>
          <WorkOrderHeader
            workOrderName={this.props.workOrder.name}
            workOrder={workOrder}
            properties={properties}
            checklist={checklist}
            checkListCategories={checkListCategories}
            locationId={locationId}
            onWorkOrderRemoved={onWorkOrderRemoved}
            onCancelClicked={onCancelClicked}
          />
          <AppContext.Consumer>
            {({user}) => (
              <FormValidationContext.Consumer>
                {validationContext => {
                  const noOwnerError = validationContext.error.check({
                    fieldId: 'Owner',
                    fieldDisplayName: 'Owner',
                    value: workOrder.ownerName,
                    required: true,
                  });
                  validationContext.editLock.check({
                    fieldId: 'status',
                    fieldDisplayName: 'Status',
                    value: this.props.workOrder.status,
                    checkCallback: value =>
                      value === doneStatus.value
                        ? `Work order is on '${doneStatus.label}' state`
                        : '',
                  });
                  validationContext.editLock.check({
                    fieldId: 'OwnerRule',
                    fieldDisplayName: 'Owner rule',
                    value: {user, workOrder: this.props.workOrder},
                    checkCallback: checkData =>
                      checkData?.user.isSuperUser ||
                      checkData?.user.email ===
                        checkData?.workOrder.ownerName ||
                      checkData?.user.email === checkData?.workOrder.assignee
                        ? ''
                        : 'User is not allowed to edit this work order',
                  });
                  const nonOwnerAssignee = validationContext.editLock.check({
                    fieldId: 'NonOwnerAssigneeRule',
                    fieldDisplayName: 'Non Owner assignee rule',
                    value: {user, workOrder: this.props.workOrder},
                    checkCallback: checkData =>
                      checkData?.user.email !==
                        checkData?.workOrder.ownerName &&
                      checkData?.user.email === checkData?.workOrder.assignee
                        ? 'Assignee is not allowed to change owner'
                        : '',
                    notAggregated: true,
                  });
                  return (
                    <div className={classes.cards}>
                      <Grid container spacing={2}>
                        <Grid item xs={8} sm={8} lg={8} xl={8}>
                          <ExpandingPanel title="Details">
                            <NameDescriptionSection
                              name={workOrder.name}
                              description={workOrder.description}
                              onNameChange={value =>
                                this._setWorkOrderDetail('name', value)
                              }
                              onDescriptionChange={value =>
                                this._setWorkOrderDetail('description', value)
                              }
                            />
                            <Grid
                              container
                              spacing={2}
                              className={classes.propertiesGrid}>
                              <Grid item xs={12} sm={6} lg={4} xl={4}>
                                <FormField label="Project">
                                  <ProjectTypeahead
                                    className={classes.gridInput}
                                    selectedProject={
                                      workOrder.project
                                        ? {
                                            id: workOrder.project.id,
                                            name: workOrder.project.name,
                                          }
                                        : null
                                    }
                                    margin="dense"
                                    onProjectSelection={project =>
                                      this._setWorkOrderDetail(
                                        'project',
                                        project,
                                      )
                                    }
                                  />
                                </FormField>
                              </Grid>
                              <Grid item xs={12} sm={6} lg={4} xl={4}>
                                <FormField label="Priority">
                                  <Select
                                    options={priorityValues}
                                    selectedValue={workOrder.priority}
                                    onChange={value =>
                                      this._setWorkOrderDetail(
                                        'priority',
                                        value,
                                      )
                                    }
                                  />
                                </FormField>
                              </Grid>
                              <Grid item xs={12} sm={6} lg={4} xl={4}>
                                <FormField
                                  label="Status"
                                  disabled={validationContext.error.detected}>
                                  <Select
                                    options={statusValues}
                                    selectedValue={workOrder.status}
                                    onChange={value =>
                                      this.setWorkOrderStatus(value)
                                    }
                                  />
                                </FormField>
                              </Grid>
                              <Grid item xs={12} sm={6} lg={4} xl={4}>
                                <FormField label="Created On">
                                  <TextInput
                                    type="date"
                                    className={classes.gridInput}
                                    value={formatDateForTextInput(
                                      workOrder.creationDate,
                                    )}
                                  />
                                </FormField>
                              </Grid>
                              <Grid item xs={12} sm={6} lg={4} xl={4}>
                                <FormField label="Due Date">
                                  <TextInput
                                    type="date"
                                    className={classes.gridInput}
                                    value={formatDateForTextInput(
                                      workOrder.installDate,
                                    )}
                                    onChange={event => {
                                      const value =
                                        event.target.value != ''
                                          ? new Date(
                                              event.target.value,
                                            ).toISOString()
                                          : '';
                                      this._setWorkOrderDetail(
                                        'installDate',
                                        value,
                                      );
                                    }}
                                  />
                                </FormField>
                              </Grid>
                              <Grid item xs={12} sm={6} lg={4} xl={4}>
                                <FormField label="Location">
                                  <LocationTypeahead
                                    headline={null}
                                    className={classes.gridInput}
                                    margin="dense"
                                    selectedLocation={
                                      location
                                        ? {
                                            id: location.id,
                                            name: location.name,
                                          }
                                        : null
                                    }
                                    onLocationSelection={location =>
                                      this._locationChangedHandler(
                                        location?.id ?? null,
                                      )
                                    }
                                  />
                                </FormField>
                              </Grid>
                              {properties.map((property, index) => (
                                <Grid
                                  key={property.id}
                                  item
                                  xs={12}
                                  sm={6}
                                  lg={4}
                                  xl={4}>
                                  <PropertyValueInput
                                    required={
                                      !!property.propertyType.isMandatory
                                    }
                                    disabled={
                                      !property.propertyType.isInstanceProperty
                                    }
                                    label={property.propertyType.name}
                                    className={classes.gridInput}
                                    margin="dense"
                                    inputType="Property"
                                    property={property}
                                    onChange={this._propertyChangedHandler(
                                      index,
                                    )}
                                    headlineVariant="form"
                                    fullWidth={true}
                                  />
                                </Grid>
                              ))}
                            </Grid>
                            <>
                              {location && (
                                <>
                                  <div className={classes.separator} />
                                  <Text weight="regular" variant="subtitle2">
                                    Location
                                  </Text>
                                  <LocationBreadcrumbsTitle
                                    locationDetails={location}
                                    size="small"
                                  />
                                  <Grid container spacing={2}>
                                    <Grid item xs={12} md={12}>
                                      <LocationMapSnippet
                                        className={classes.map}
                                        location={{
                                          id: location.id,
                                          name: location.name,
                                          latitude: location.latitude,
                                          longitude: location.longitude,
                                          locationType: {
                                            mapType:
                                              location.locationType.mapType,
                                            mapZoomLevel: (
                                              location.locationType
                                                .mapZoomLevel || 8
                                            ).toString(),
                                          },
                                        }}
                                      />
                                    </Grid>
                                  </Grid>
                                </>
                              )}
                            </>
                          </ExpandingPanel>
                          {actionsEnabled && (
                            <ExpandingPanel title="Actions">
                              <WorkOrderDetailsPane workOrder={workOrder} />
                            </ExpandingPanel>
                          )}
                          <ExpandingPanel
                            title="Attachments"
                            rightContent={
                              <div className={classes.uploadButtonContainer}>
                                <AddHyperlinkButton
                                  className={classes.minimizedButton}
                                  skin="regular"
                                  entityType="WORK_ORDER"
                                  allowCategories={false}
                                  entityId={workOrder.id}>
                                  <InsertLinkIcon color="primary" />
                                </AddHyperlinkButton>
                                {this.state.isLoadingDocument ? (
                                  <CircularProgress size={24} />
                                ) : (
                                  <FileUpload
                                    className={classes.minimizedButton}
                                    button={
                                      <CloudUploadOutlinedIcon
                                        className={classes.uploadButton}
                                      />
                                    }
                                    onFileUploaded={this.onDocumentUploaded}
                                    onProgress={() =>
                                      this.setState({isLoadingDocument: true})
                                    }
                                  />
                                )}
                              </div>
                            }>
                            <EntityDocumentsTable
                              entityType="WORK_ORDER"
                              entityId={workOrder.id}
                              files={[
                                ...this.props.workOrder.files,
                                ...this.props.workOrder.images,
                              ]}
                              hyperlinks={this.props.workOrder.hyperlinks}
                            />
                          </ExpandingPanel>
                          <ExpandingPanel
                            title={fbt('Checklist', 'Checklist section header')}
                            rightContent={
                              <EditToggleButton
                                isOnEdit={showChecklistDesignMode}
                                onChange={newToggleValue =>
                                  this.setState({
                                    showChecklistDesignMode: newToggleValue,
                                  })
                                }
                              />
                            }>
                            <CheckListTable
                              list={checklist}
                              onChecklistChanged={this._checklistChangedHandler}
                              onDesignMode={this.state.showChecklistDesignMode}
                            />
                          </ExpandingPanel>
                          <CheckListCategoryContextProvider>
                            <CheckListCategoryContext.Consumer>
                              {categoryContext => (
                                <ExpandingPanel
                                  allowExpandCollapse={hasCheckListCategories}
                                  title={fbt(
                                    'Checklist Categories',
                                    'Checklist section header',
                                  )}
                                  rightContent={
                                    <Button
                                      variant="text"
                                      onClick={() =>
                                        categoryContext.call.addNewCategory()
                                      }>
                                      {hasCheckListCategories ? (
                                        fbt(
                                          'Add category',
                                          'Add checklist category button text',
                                        )
                                      ) : (
                                        <AddIcon color="primary" />
                                      )}
                                    </Button>
                                  }>
                                  <CheckListCategoryTable
                                    list={checkListCategories}
                                    onListChanged={
                                      this._checkListCategoryChangedHandler
                                    }
                                  />
                                </ExpandingPanel>
                              )}
                            </CheckListCategoryContext.Consumer>
                          </CheckListCategoryContextProvider>
                        </Grid>
                        <Grid item xs={4} sm={4} lg={4} xl={4}>
                          <ExpandingPanel title="Team" className={classes.card}>
                            <FormField
                              className={classes.input}
                              label="Owner"
                              required={true}
                              hasError={!!noOwnerError}
                              errorText={noOwnerError}
                              disabled={!!nonOwnerAssignee}>
                              <UserTypeahead
                                selectedUser={workOrder.ownerName}
                                onUserSelection={user =>
                                  this._setWorkOrderDetail('ownerName', user)
                                }
                                margin="dense"
                              />
                            </FormField>
                            <FormField
                              label="Assignee"
                              className={classes.input}>
                              <UserTypeahead
                                selectedUser={workOrder.assignee}
                                onUserSelection={user =>
                                  this._setWorkOrderDetail('assignee', user)
                                }
                                margin="dense"
                              />
                            </FormField>
                          </ExpandingPanel>
                          <ExpandingPanel
                            title="Comments"
                            detailsPaneClass={classes.commentsBoxContainer}
                            className={classes.card}>
                            <CommentsBox
                              boxElementsClass={classes.inExpandingPanelFix}
                              commentsLogClass={classes.commentsLog}
                              relatedEntityId={this.props.workOrder.id}
                              relatedEntityType="WORK_ORDER"
                              comments={this.props.workOrder.comments}
                            />
                          </ExpandingPanel>
                        </Grid>
                      </Grid>
                    </div>
                  );
                }}
              </FormValidationContext.Consumer>
            )}
          </AppContext.Consumer>
        </FormValidationContextProvider>
      </div>
    );
  }

  onDocumentUploaded = (file, key) => {
    const workOrderId = this.props.workOrder.id;
    const variables: AddImageMutationVariables = {
      input: {
        entityType: 'WORK_ORDER',
        entityId: workOrderId,
        imgKey: key,
        fileName: file.name,
        fileSize: file.size,
        modified: new Date(file.lastModified).toISOString(),
        contentType: file.type,
      },
    };

    const updater = store => {
      const newNode = store.getRootField('addImage');
      const fileType = newNode.getValue('fileType');

      const workOrderProxy = store.get(workOrderId);
      if (fileType === FileTypeEnum.IMAGE) {
        const imageNodes = workOrderProxy.getLinkedRecords('images') || [];
        workOrderProxy.setLinkedRecords([...imageNodes, newNode], 'images');
      } else {
        const fileNodes = workOrderProxy.getLinkedRecords('files') || [];
        workOrderProxy.setLinkedRecords([...fileNodes, newNode], 'files');
      }
    };

    const callbacks: MutationCallbacks<AddImageMutationResponse> = {
      onCompleted: () => {
        this.setState({
          isLoadingDocument: false,
        });
      },
      onError: () => {},
    };

    AddImageMutation(variables, callbacks, updater);
  };

  setWorkOrderStatus = value => {
    if (!value || value == this.state.workOrder.status) {
      return;
    }

    const verification = new Promise((resolve, reject) => {
      if (value != doneStatus.value) {
        resolve();
      } else {
        this.props
          .confirm({
            title: fbt(
              // eslint-disable-next-line prettier/prettier
              "Are you sure you want to mark this work order as 'Done'?",
              'Verification message title',
            ),
            message: fbt(
              // eslint-disable-next-line prettier/prettier
              "Once saved with 'Done' status, the work order will be locked for editing.",
              'Verification message details',
            ),
            confirmLabel: Strings.common.okButton,
          })
          .then(confirmed => {
            if (confirmed) {
              resolve();
            } else {
              reject();
            }
          });
      }
    });

    verification.then(() => {
      this.setState({
        workOrder: update(this.state.workOrder, {status: {$set: value}}),
      });
    });
  };

  _checkListCategoryChangedHandler = updatedCheckListCategories => {
    this.setState(() => {
      return {
        checkListCategories: updatedCheckListCategories,
      };
    });
  };

  _checklistChangedHandler = updatedChecklist => {
    this.setState(() => ({
      checklist: updatedChecklist,
    }));
  };

  _setWorkOrderDetail = (
    key:
      | 'name'
      | 'description'
      | 'ownerName'
      | 'installDate'
      | 'assignee'
      | 'priority'
      | 'project',
    value,
  ) => {
    this.setState(prevState => {
      return {
        // $FlowFixMe Set state for each field
        workOrder: update(prevState.workOrder, {[key]: {$set: value}}),
      };
    });
  };

  _locationChangedHandler = (locationId: ?string) =>
    this.setState({locationId});

  _propertyChangedHandler = index => property => {
    this.setState(prevState => {
      return {
        properties: update(prevState.properties, {[index]: {$set: property}}),
      };
    });
  };
  _enqueueError = (message: string) => {
    this.props.enqueueSnackbar(message, {
      children: key => (
        <SnackbarItem id={key} message={message} variant="error" />
      ),
    });
  };
}

export default withRouter(
  withSnackbar(
    withStyles(styles)(
      withAlert(
        createFragmentContainer(WorkOrderDetails, {
          workOrder: graphql`
            fragment WorkOrderDetails_workOrder on WorkOrder {
              id
              name
              description
              workOrderType {
                name
                id
              }
              location {
                name
                id
                latitude
                longitude
                locationType {
                  mapType
                  mapZoomLevel
                }
                ...LocationBreadcrumbsTitle_locationDetails
              }
              ownerName
              assignee
              creationDate
              installDate
              status
              priority
              ...WorkOrderDetailsPane_workOrder
              properties {
                ...PropertyFormField_property @relay(mask: false)
              }
              images {
                ...EntityDocumentsTable_files
              }
              files {
                ...EntityDocumentsTable_files
              }
              hyperlinks {
                ...EntityDocumentsTable_hyperlinks
              }
              comments {
                ...CommentsBox_comments
              }
              project {
                name
                id
                type {
                  id
                  name
                }
              }
              checkList {
                ...CheckListTable_list @relay(mask: false)
              }
              checkListCategories {
                ...CheckListCategoryTable_list
                id
                title
                description
                checkList {
                  ...CheckListCategoryItemsDialog_items
                }
              }
            }
          `,
        }),
      ),
    ),
  ),
);
