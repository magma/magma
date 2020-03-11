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
import type {ChecklistCategoriesMutateStateActionType} from '../checklist/ChecklistCategoriesMutateAction';
import type {ChecklistCategoriesStateType} from '../checklist/ChecklistCategoriesMutateState';
import type {ContextRouter} from 'react-router-dom';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WorkOrderDetails_workOrder} from './__generated__/WorkOrderDetails_workOrder.graphql.js';

import AddHyperlinkButton from '../AddHyperlinkButton';
import AddImageMutation from '../../mutations/AddImageMutation';
import AppContext from '@fbcnms/ui/context/AppContext';
import CheckListCategoryExpandingPanel from '../checklist/checkListCategory/CheckListCategoryExpandingPanel';
import ChecklistCategoriesMutateDispatchContext from '../checklist/ChecklistCategoriesMutateDispatchContext';
import CircularProgress from '@material-ui/core/CircularProgress';
import CloudUploadOutlinedIcon from '@material-ui/icons/CloudUploadOutlined';
import CommentsBox from '../comments/CommentsBox';
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
import React, {useContext, useReducer, useState} from 'react';
import Select from '@fbcnms/ui/components/design-system/Select/Select';
import Strings from '../../common/CommonStrings';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import UserTypeahead from '../typeahead/UserTypeahead';
import WorkOrderDetailsPane from './WorkOrderDetailsPane';
import WorkOrderHeader from './WorkOrderHeader';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {createFragmentContainer, graphql} from 'react-relay';
import {doneStatus, priorityValues, statusValues} from '../../common/WorkOrder';
import {formatDateForTextInput} from '@fbcnms/ui/utils/displayUtils';
import {
  getInitialState,
  reducer,
} from '../checklist/ChecklistCategoriesMutateReducer';
import {makeStyles} from '@material-ui/styles';
import {sortPropertiesByIndex} from '../../common/Property';
import {withRouter} from 'react-router-dom';

type Props = {
  workOrder: WorkOrderDetails_workOrder,
  onWorkOrderRemoved: () => void,
  onCancelClicked: () => void,
  ...WithAlert,
  ...ContextRouter,
};

const FileTypeEnum = {
  IMAGE: 'IMAGE',
  FILE: 'FILE',
};

const useStyles = makeStyles(() => ({
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
    borderBottom: `1px solid ${symphony.palette.D50}`,
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
    fill: symphony.palette.primary,
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
}));

const WorkOrderDetails = ({
  workOrder: propsWorkOrder,
  onWorkOrderRemoved,
  onCancelClicked,
  confirm,
}: Props) => {
  const classes = useStyles();
  const [workOrder, setWorkOrder] = useState(propsWorkOrder);
  const [properties, setProperties] = useState(
    // eslint-disable-next-line flowtype/no-weak-types
    ([...propsWorkOrder.properties]: any).sort(sortPropertiesByIndex),
  );
  const [locationId, setLocationId] = useState(propsWorkOrder.location?.id);
  const [isLoadingDocument, setIsLoadingDocument] = useState(false);
  const {user, isFeatureEnabled} = useContext(AppContext);

  const [editingCategories, dispatch] = useReducer<
    ChecklistCategoriesStateType,
    ChecklistCategoriesMutateStateActionType,
    $ElementType<WorkOrderDetails_workOrder, 'checkListCategories'>,
  >(reducer, propsWorkOrder.checkListCategories, getInitialState);

  const onDocumentUploaded = (file, key) => {
    const workOrderId = propsWorkOrder.id;
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
      // $FlowFixMe (T62907961) Relay flow types
      const newNode = store.getRootField('addImage');
      const fileType = newNode.getValue('fileType');

      // $FlowFixMe (T62907961) Relay flow types
      const workOrderProxy = store.get(workOrderId);
      if (fileType === FileTypeEnum.IMAGE) {
        // $FlowFixMe (T62907961) Relay flow types
        const imageNodes = workOrderProxy.getLinkedRecords('images') || [];
        // $FlowFixMe (T62907961) Relay flow types
        workOrderProxy.setLinkedRecords([...imageNodes, newNode], 'images');
      } else {
        // $FlowFixMe (T62907961) Relay flow types
        const fileNodes = workOrderProxy.getLinkedRecords('files') || [];
        // $FlowFixMe (T62907961) Relay flow types
        workOrderProxy.setLinkedRecords([...fileNodes, newNode], 'files');
      }
    };

    const callbacks: MutationCallbacks<AddImageMutationResponse> = {
      onCompleted: () => {
        setIsLoadingDocument(false);
      },
      onError: () => {},
    };

    AddImageMutation(variables, callbacks, updater);
  };

  const setWorkOrderStatus = value => {
    if (!value || value == workOrder.status) {
      return;
    }

    const verification = new Promise((resolve, reject) => {
      if (value != doneStatus.value) {
        resolve();
      } else {
        confirm({
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
        }).then(confirmed => {
          if (confirmed) {
            resolve();
          } else {
            reject();
          }
        });
      }
    });

    verification.then(() => {
      setWorkOrder({...workOrder, status: value});
    });
  };

  const _setWorkOrderDetail = (
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
    setWorkOrder(prevWorkOrder => ({...prevWorkOrder, [`${key}`]: value}));
  };

  const {location} = workOrder;
  const actionsEnabled = isFeatureEnabled('planned_equipment');
  return (
    <div className={classes.root}>
      <FormValidationContextProvider>
        <WorkOrderHeader
          workOrderName={propsWorkOrder.name}
          workOrder={workOrder}
          properties={properties}
          checkListCategories={editingCategories}
          locationId={locationId}
          onWorkOrderRemoved={onWorkOrderRemoved}
          onCancelClicked={onCancelClicked}
        />
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
              value: propsWorkOrder.status,
              checkCallback: value =>
                value === doneStatus.value
                  ? `Work order is on '${doneStatus.label}' state`
                  : '',
            });
            validationContext.editLock.check({
              fieldId: 'OwnerRule',
              fieldDisplayName: 'Owner rule',
              value: {user, workOrder: propsWorkOrder},
              checkCallback: checkData =>
                checkData?.user.isSuperUser ||
                checkData?.user.email === checkData?.workOrder.ownerName ||
                checkData?.user.email === checkData?.workOrder.assignee
                  ? ''
                  : 'User is not allowed to edit this work order',
            });
            const nonOwnerAssignee = validationContext.editLock.check({
              fieldId: 'NonOwnerAssigneeRule',
              fieldDisplayName: 'Non Owner assignee rule',
              value: {user, workOrder: propsWorkOrder},
              checkCallback: checkData =>
                checkData?.user.email !== checkData?.workOrder.ownerName &&
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
                          _setWorkOrderDetail('name', value)
                        }
                        onDescriptionChange={value =>
                          _setWorkOrderDetail('description', value)
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
                                _setWorkOrderDetail('project', project)
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
                                _setWorkOrderDetail('priority', value)
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
                              onChange={value => setWorkOrderStatus(value)}
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
                                    ? new Date(event.target.value).toISOString()
                                    : '';
                                _setWorkOrderDetail('installDate', value);
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
                                setLocationId(location?.id ?? null)
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
                              required={!!property.propertyType.isMandatory}
                              disabled={
                                !property.propertyType.isInstanceProperty
                              }
                              label={property.propertyType.name}
                              className={classes.gridInput}
                              margin="dense"
                              inputType="Property"
                              property={property}
                              onChange={property =>
                                setProperties(prevProperties => [
                                  ...prevProperties.slice(0, index),
                                  property,
                                  prevProperties.slice(index + 1),
                                ])
                              }
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
                                      mapType: location.locationType.mapType,
                                      mapZoomLevel: (
                                        location.locationType.mapZoomLevel || 8
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
                          {isLoadingDocument ? (
                            <CircularProgress size={24} />
                          ) : (
                            <FileUpload
                              className={classes.minimizedButton}
                              button={
                                <CloudUploadOutlinedIcon
                                  className={classes.uploadButton}
                                />
                              }
                              onFileUploaded={onDocumentUploaded}
                              onProgress={() => setIsLoadingDocument(true)}
                            />
                          )}
                        </div>
                      }>
                      <EntityDocumentsTable
                        entityType="WORK_ORDER"
                        entityId={workOrder.id}
                        files={[
                          ...propsWorkOrder.files,
                          ...propsWorkOrder.images,
                        ]}
                        hyperlinks={propsWorkOrder.hyperlinks}
                      />
                    </ExpandingPanel>
                    <ChecklistCategoriesMutateDispatchContext.Provider
                      value={dispatch}>
                      <CheckListCategoryExpandingPanel
                        categories={editingCategories}
                      />
                    </ChecklistCategoriesMutateDispatchContext.Provider>
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
                            _setWorkOrderDetail('ownerName', user)
                          }
                          margin="dense"
                        />
                      </FormField>
                      <FormField label="Assignee" className={classes.input}>
                        <UserTypeahead
                          selectedUser={workOrder.assignee}
                          onUserSelection={user =>
                            _setWorkOrderDetail('assignee', user)
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
                        relatedEntityId={propsWorkOrder.id}
                        relatedEntityType="WORK_ORDER"
                        comments={propsWorkOrder.comments}
                      />
                    </ExpandingPanel>
                  </Grid>
                </Grid>
              </div>
            );
          }}
        </FormValidationContext.Consumer>
      </FormValidationContextProvider>
    </div>
  );
};

export default withRouter(
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
          checkListCategories {
            id
            title
            description
            checkList {
              id
              index
              type
              title
              helpText
              checked
              enumValues
              stringValue
            }
          }
        }
      `,
    }),
  ),
);
