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
import type {Property} from '../../common/Property';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WorkOrderDetails_workOrder} from './__generated__/WorkOrderDetails_workOrder.graphql.js';

import AddHyperlinkButton from '../AddHyperlinkButton';
import AddImageMutation from '../../mutations/AddImageMutation';
import AppContext from '@fbcnms/ui/context/AppContext';
import CheckListCategoryExpandingPanel from '../checklist/checkListCategory/CheckListCategoryExpandingPanel';
import ChecklistCategoriesMutateDispatchContext from '../checklist/ChecklistCategoriesMutateDispatchContext';
import CircularProgress from '@material-ui/core/CircularProgress';
import CommentsBox from '../comments/CommentsBox';
import EntityDocumentsTable from '../EntityDocumentsTable';
import ExpandingPanel from '@fbcnms/ui/components/ExpandingPanel';
import FileUploadButton from '../FileUpload/FileUploadButton';
import FormContext, {FormContextProvider} from '../../common/FormContext';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import Grid from '@material-ui/core/Grid';
import IconButton from '@fbcnms/ui/components/design-system/IconButton';
import LinkIcon from '@fbcnms/ui/components/design-system/Icons/Actions/LinkIcon';
import LocationBreadcrumbsTitle from '../location/LocationBreadcrumbsTitle';
import LocationMapSnippet from '../location/LocationMapSnippet';
import LocationTypeahead from '../typeahead/LocationTypeahead';
import NameDescriptionSection from '@fbcnms/ui/components/NameDescriptionSection';
import ProjectTypeahead from '../typeahead/ProjectTypeahead';
import PropertyValueInput from '../form/PropertyValueInput';
import React, {useContext, useReducer, useState} from 'react';
import Select from '@fbcnms/ui/components/design-system/Select/Select';
import Strings from '@fbcnms/strings/Strings';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import UploadIcon from '@fbcnms/ui/components/design-system/Icons/Actions/UploadIcon';
import UserTypeahead from '../typeahead/UserTypeahead';
import WorkOrderDetailsPane from './WorkOrderDetailsPane';
import WorkOrderHeader from './WorkOrderHeader';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {NAVIGATION_OPTIONS} from '../location/LocationBreadcrumbsTitle';
import {createFragmentContainer, graphql} from 'react-relay';
import {doneStatus, priorityValues, statusValues} from '../../common/WorkOrder';
import {formatDateForTextInput} from '@fbcnms/ui/utils/displayUtils';
import {
  getInitialState,
  reducer,
} from '../checklist/ChecklistCategoriesMutateReducer';
import {makeStyles} from '@material-ui/styles';
import {sortPropertiesByIndex, toMutableProperty} from '../../common/Property';
import {useMainContext} from '../MainContext';
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
    marginRight: '4px',
    marginLeft: '8px',
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
  const [workOrder, setWorkOrder] = useState<WorkOrderDetails_workOrder>(
    propsWorkOrder,
  );
  const [properties, setProperties] = useState<Array<Property>>(
    propsWorkOrder.properties
      .filter(Boolean)
      .slice()
      .map<Property>(toMutableProperty)
      .sort(sortPropertiesByIndex),
  );
  const [locationId, setLocationId] = useState(propsWorkOrder.location?.id);
  const [isLoadingDocument, setIsLoadingDocument] = useState(false);
  const {isFeatureEnabled} = useContext(AppContext);

  const {userHasAdminPermissions, me} = useMainContext();

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
      const newNode = store.getRootField('addImage');
      const workOrderProxy = store.get(workOrderId);
      if (newNode == null || workOrderProxy == null) {
        return;
      }

      const fileType = newNode.getValue('fileType');
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
      | 'owner'
      | 'installDate'
      | 'assignedTo'
      | 'priority'
      | 'project',
    value,
  ) => {
    setWorkOrder(prevWorkOrder => ({...prevWorkOrder, [`${key}`]: value}));
  };

  const {location} = workOrder;
  const actionsEnabled = isFeatureEnabled('planned_equipment');
  const permissionsEnforcementIsOn = isFeatureEnabled(
    'permissions_ui_enforcement',
  );

  const isOwnerOrAssignee =
    me?.user?.email === workOrder?.owner?.email ||
    me?.user?.email === workOrder?.assignedTo?.email;

  return (
    <div className={classes.root}>
      <FormContextProvider
        permissions={{
          entity: 'workorder',
          action: 'update',
          ignore: isOwnerOrAssignee,
        }}>
        <WorkOrderHeader
          workOrderName={propsWorkOrder.name}
          workOrder={workOrder}
          properties={properties}
          checkListCategories={editingCategories}
          locationId={locationId}
          onWorkOrderRemoved={onWorkOrderRemoved}
          onCancelClicked={onCancelClicked}
        />
        <FormContext.Consumer>
          {form => {
            form.alerts.editLock.check({
              fieldId: 'status',
              fieldDisplayName: 'Status',
              value: propsWorkOrder.status,
              checkCallback: value =>
                value === doneStatus.value
                  ? `Work order is on '${doneStatus.label}' state`
                  : '',
            });
            const nonOwnerAssignee =
              permissionsEnforcementIsOn &&
              form.alerts.editLock.check({
                fieldId: 'NonOwnerAssigneeRule',
                fieldDisplayName: 'Non Owner assignee rule',
                value: propsWorkOrder,
                checkCallback: workOrder =>
                  !userHasAdminPermissions &&
                  me?.user?.email !== workOrder?.owner.email &&
                  me?.user?.email === workOrder?.assignedTo?.email
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
                            disabled={form.alerts.error.detected}>
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
                                  ...prevProperties.slice(index + 1),
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
                              navigateOnClick={NAVIGATION_OPTIONS.NEW_TAB}
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
                            variant="text"
                            entityType="WORK_ORDER"
                            allowCategories={false}
                            entityId={workOrder.id}>
                            <IconButton icon={LinkIcon} />
                          </AddHyperlinkButton>
                          {isLoadingDocument ? (
                            <CircularProgress size={24} />
                          ) : (
                            <FileUploadButton
                              onFileUploaded={onDocumentUploaded}
                              onProgress={() => setIsLoadingDocument(true)}>
                              {openFileUploadDialog => (
                                <IconButton
                                  className={classes.minimizedButton}
                                  onClick={openFileUploadDialog}
                                  icon={UploadIcon}
                                />
                              )}
                            </FileUploadButton>
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
                        validation={{id: 'owner', value: workOrder.owner?.id}}
                        disabled={!!nonOwnerAssignee}>
                        <UserTypeahead
                          selectedUser={workOrder.owner}
                          onUserSelection={user =>
                            _setWorkOrderDetail('owner', user)
                          }
                          margin="dense"
                        />
                      </FormField>
                      <FormField label="Assignee" className={classes.input}>
                        <UserTypeahead
                          selectedUser={workOrder.assignedTo}
                          onUserSelection={user =>
                            _setWorkOrderDetail('assignedTo', user)
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
        </FormContext.Consumer>
      </FormContextProvider>
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
          owner {
            id
            email
          }
          assignedTo {
            id
            email
          }
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
              enumSelectionMode
              selectedEnumValues
              yesNoResponse
              files {
                id
                fileName
                sizeInBytes
                modified
                uploaded
                fileType
                storeKey
                category
              }
              cellData {
                id
                networkType
                signalStrength
                timestamp
                baseStationID
                networkID
                systemID
                cellID
                locationAreaCode
                mobileCountryCode
                mobileNetworkCode
                primaryScramblingCode
                operator
                arfcn
                physicalCellID
                trackingAreaCode
                timingAdvance
                earfcn
                uarfcn
                latitude
                longitude
              }
              wifiData {
                id
                timestamp
                frequency
                channel
                bssid
                strength
                ssid
                band
                channelWidth
                capabilities
                latitude
                longitude
              }
            }
          }
        }
      `,
    }),
  ),
);
