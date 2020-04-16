/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import ProjectMoreActionsButton from './ProjectMoreActionsButton';
import type {ContextRouter} from 'react-router-dom';
import type {
  EditProjectMutationResponse,
  EditProjectMutationVariables,
} from '../../mutations/__generated__/EditProjectMutation.graphql';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {ProjectDetails_project} from './__generated__/ProjectDetails_project.graphql.js';
import type {Property} from '../../common/Property';
import type {Theme, WithStyles} from '@material-ui/core';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithSnackbarProps} from 'notistack';

import Breadcrumbs from '@fbcnms/ui/components/Breadcrumbs';
import CommentsBox from '../comments/CommentsBox';
import EditProjectMutation from '../../mutations/EditProjectMutation';
import ExpandingPanel from '@fbcnms/ui/components/ExpandingPanel';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import FormSaveCancelPanel from '@fbcnms/ui/components/design-system/Form/FormSaveCancelPanel';
import Grid from '@material-ui/core/Grid';
import LocationBreadcrumbsTitle from '../location/LocationBreadcrumbsTitle';
import LocationMapSnippet from '../location/LocationMapSnippet';
import LocationTypeahead from '../typeahead/LocationTypeahead';
import NameDescriptionSection from '@fbcnms/ui/components/NameDescriptionSection';
import ProjectWorkOrdersList from './ProjectWorkOrdersList';
import PropertyValueInput from '../form/PropertyValueInput';
import React from 'react';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import UserTypeahead from '../typeahead/UserTypeahead';
import update from 'immutability-helper';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {FormContextProvider} from '../../common/FormContext';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {createFragmentContainer, graphql} from 'react-relay';
import {getGraphError} from '../../common/EntUtils';
import {sortPropertiesByIndex, toPropertyInput} from '../../common/Property';
import {withRouter} from 'react-router-dom';
import {withSnackbar} from 'notistack';
import {withStyles} from '@material-ui/core/styles';

type State = {
  editedProject: ProjectDetails_project,
  locationId: ?string,
  properties: Array<Property>,
};

type Props = {
  project: ProjectDetails_project,
  onProjectRemoved: () => void,
} & WithAlert &
  WithStyles<typeof styles> &
  ContextRouter &
  WithSnackbarProps;

const styles = (theme: Theme) => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
    flexGrow: 1,
  },
  labelName: {
    fontSize: '20px',
    fontWeight: 500,
    lineHeight: '28px',
    textAlign: 'left',
    paddingBottom: '24px',
    color: theme.palette.blueGrayDark,
  },
  description: {
    margin: '10px',
  },
  input: {
    paddingBottom: '24px',
  },
  gridInput: {
    display: 'inline-flex',
  },
  cards: {
    flexGrow: 1,
    overflowY: 'auto',
    overflowX: 'hidden',
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
  breadcrumbs: {
    flexGrow: 1,
  },
  propertiesGrid: {
    marginTop: '16px',
  },
  button: {
    marginRight: '8px',
  },
  nameHeader: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: '24px',
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

class ProjectDetails extends React.Component<Props, State> {
  state = {
    editedProject: this.props.project,
    properties: this.getEditingProperties(),
    locationId: this.props.project.location?.id,
  };

  getEditingProperties(): Array<Property> {
    // eslint-disable-next-line flowtype/no-weak-types
    return ([...this.props.project.properties]: any).sort(
      sortPropertiesByIndex,
    );
  }

  _setProjectDetail = (key: 'name' | 'description' | 'createdBy', value) => {
    this.setState(prevState => {
      return {
        // $FlowFixMe Set state for each field
        editedProject: update(prevState.editedProject, {[key]: {$set: value}}),
      };
    });
  };

  _propertyChangedHandler = index => property => {
    this.setState(prevState => {
      return {
        properties: update(prevState.properties, {[index]: {$set: property}}),
      };
    });
  };

  _locationChangedHandler = (locationId: ?string) =>
    this.setState({locationId});

  saveProject = () => {
    const {id, name, description, createdBy, type} = this.state.editedProject;
    const variables: EditProjectMutationVariables = {
      input: {
        id,
        name,
        description,
        creatorId: createdBy?.id,
        type: type.id,
        properties: toPropertyInput(this.state.properties),
        location: this.state.locationId,
      },
    };
    const callbacks: MutationCallbacks<EditProjectMutationResponse> = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          const msg = errors[0].message;
          this.props.enqueueSnackbar(msg, {
            children: key => (
              <SnackbarItem id={key} message={msg} variant="error" />
            ),
          });
        } else {
          // navigate to main page
          this.props.history.push(this.props.match.url);
        }
      },
      onError: (error: Error) => {
        const msg = getGraphError(error);
        this.props.enqueueSnackbar(msg, {
          children: key => (
            <SnackbarItem id={key} message={msg} variant="error" />
          ),
        });
      },
    };
    ServerLogger.info(LogEvents.SAVE_PROJECT_BUTTON_CLICKED, {
      source: 'project_details',
    });
    EditProjectMutation(variables, callbacks);
  };

  render() {
    const {classes, onProjectRemoved} = this.props;
    const project = this.state.editedProject;
    const {location} = project;
    const {properties} = this.state;
    return (
      <div className={classes.root}>
        <FormContextProvider>
          <div className={classes.nameHeader}>
            <div className={classes.breadcrumbs}>
              <Breadcrumbs
                breadcrumbs={[
                  {
                    id: 'projects',
                    name: 'Projects',
                    onClick: () => this.navigateToMainPage(),
                  },
                  {
                    id: project.id,
                    name: this.props.project.name,
                    subtext: `ID: ${project.id}`,
                  },
                ]}
                size="large"
              />
            </div>
            <ProjectMoreActionsButton
              className={classes.button}
              project={project}
              onProjectRemoved={onProjectRemoved}
            />
            <FormSaveCancelPanel
              onCancel={() => this.props.history.push(this.props.match.url)}
              onSave={this.saveProject}
            />
          </div>
          <div className={classes.cards}>
            <Grid container spacing={2}>
              <Grid item xs={8} sm={8} lg={8} xl={8}>
                <ExpandingPanel title="Details">
                  <NameDescriptionSection
                    name={project.name}
                    description={project.description}
                    onNameChange={value =>
                      this._setProjectDetail('name', value)
                    }
                    onDescriptionChange={value =>
                      this._setProjectDetail('description', value)
                    }
                  />
                  <Grid
                    container
                    spacing={2}
                    className={classes.propertiesGrid}>
                    {project.type && (
                      <Grid item xs={12} sm={6} lg={4} xl={4}>
                        <FormField label="Type">
                          <TextInput
                            disabled={true}
                            variant="outlined"
                            type="string"
                            value={project.type.name}
                          />
                        </FormField>
                      </Grid>
                    )}
                    <Grid item xs={12} sm={6} lg={4} xl={4}>
                      <FormField label="Location">
                        <LocationTypeahead
                          headline={null}
                          className={classes.gridInput}
                          margin="dense"
                          selectedLocation={
                            location
                              ? {id: location.id, name: location.name}
                              : null
                          }
                          onLocationSelection={location =>
                            this._locationChangedHandler(location?.id ?? null)
                          }
                        />
                      </FormField>
                    </Grid>
                    {properties.map((property, index) => (
                      <Grid key={property.id} item xs={12} sm={6} lg={4} xl={4}>
                        <PropertyValueInput
                          required={!!property.propertyType.isMandatory}
                          disabled={!property.propertyType.isInstanceProperty}
                          headlineVariant="form"
                          fullWidth={true}
                          label={property.propertyType.name}
                          className={classes.gridInput}
                          margin="dense"
                          inputType="Property"
                          property={property}
                          onChange={this._propertyChangedHandler(index)}
                        />
                      </Grid>
                    ))}
                  </Grid>
                  <>
                    {location && (
                      <>
                        <div className={classes.separator} />

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
                <ExpandingPanel title="Work Orders">
                  <ProjectWorkOrdersList
                    workOrders={project.workOrders}
                    onNavigateToWorkOrder={this.navigateToWorkOrder}
                  />
                </ExpandingPanel>
              </Grid>
              <Grid item xs={4} sm={4} lg={4} xl={4}>
                <ExpandingPanel title="Team">
                  <FormField>
                    <UserTypeahead
                      className={classes.input}
                      selectedUser={project.createdBy}
                      headline="Owner"
                      onUserSelection={user =>
                        this._setProjectDetail('createdBy', user)
                      }
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
                    relatedEntityId={project.id}
                    relatedEntityType="PROJECT"
                    comments={this.props.project.comments}
                  />
                </ExpandingPanel>
              </Grid>
            </Grid>
          </div>
        </FormContextProvider>
      </div>
    );
  }

  navigateToMainPage = () => {
    ServerLogger.info(LogEvents.PROJECTS_SEARCH_NAV_CLICKED, {
      source: 'project_details',
    });
    const {match} = this.props;
    this.props.history.push(match.url);
  };

  navigateToWorkOrder = (WorkOrderId: ?string) => {
    const {history} = this.props;
    if (WorkOrderId) {
      ServerLogger.info(LogEvents.WORK_ORDER_DETAILS_NAV_CLICKED, {
        source: 'project_details',
      });
      history.push(`/workorders/search?workorder=${WorkOrderId}`);
    }
  };
}

export default withRouter(
  withSnackbar(
    withStyles(styles)(
      withAlert(
        createFragmentContainer(ProjectDetails, {
          project: graphql`
            fragment ProjectDetails_project on Project {
              id
              name
              description
              createdBy {
                id
                email
              }
              type {
                name
                id
              }
              location {
                id
                name
                latitude
                longitude
                locationType {
                  mapType
                  mapZoomLevel
                }
                ...LocationBreadcrumbsTitle_locationDetails
              }
              properties {
                id
                stringValue
                intValue
                floatValue
                booleanValue
                latitudeValue
                longitudeValue
                rangeFromValue
                rangeToValue
                nodeValue {
                  id
                  name
                }
                propertyType {
                  id
                  name
                  type
                  nodeType
                  isEditable
                  isMandatory
                  isInstanceProperty
                  stringValue
                  intValue
                  floatValue
                  booleanValue
                  latitudeValue
                  longitudeValue
                  rangeFromValue
                  rangeToValue
                }
              }
              workOrders {
                ...ProjectWorkOrdersList_workOrders
              }
              comments {
                ...CommentsBox_comments
              }
            }
          `,
        }),
      ),
    ),
  ),
);
