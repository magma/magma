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
  AddProjectMutationResponse,
  AddProjectMutationVariables,
} from '../../mutations/__generated__/AddProjectMutation.graphql';
import type {ContextRouter} from 'react-router-dom';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {Project} from '../../common/Project';
import type {WithSnackbarProps} from 'notistack';
import type {WithStyles} from '@material-ui/core';

import AddProjectMutation from '../../mutations/AddProjectMutation';
import Breadcrumbs from '@fbcnms/ui/components/Breadcrumbs';
import CircularProgress from '@material-ui/core/CircularProgress';
import ExpandingPanel from '@fbcnms/ui/components/ExpandingPanel';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import FormSaveCancelPanel from '@fbcnms/ui/components/design-system/Form/FormSaveCancelPanel';
import Grid from '@material-ui/core/Grid';
import LocationTypeahead from '../typeahead/LocationTypeahead';
import NameDescriptionSection from '@fbcnms/ui/components/NameDescriptionSection';
import PropertyValueInput from '../form/PropertyValueInput';
import React from 'react';
import RelayEnvironment from '../../common/RelayEnvironment.js';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import TextField from '@material-ui/core/TextField';
import UserTypeahead from '../typeahead/UserTypeahead';
import nullthrows from '@fbcnms/util/nullthrows';
import update from 'immutability-helper';
import {FormValidationContextProvider} from '@fbcnms/ui/components/design-system/Form/FormValidationContext';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {fetchQuery, graphql} from 'relay-runtime';
import {getInitialPropertyFromType} from '../../common/PropertyType';
import {sortPropertiesByIndex, toPropertyInput} from '../../common/Property';
import {withRouter} from 'react-router-dom';
import {withSnackbar} from 'notistack';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  projectTypeId: ?string,
} & WithStyles<typeof styles> &
  ContextRouter &
  WithSnackbarProps;

const styles = theme => ({
  root: {
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
    padding: '40px 32px',
  },
  contentRoot: {
    position: 'relative',
    flexGrow: 1,
    overflowX: 'auto',
    display: 'flex',
  },
  cards: {
    marginTop: '24px',
    overflowY: 'auto',
    overflowX: 'hidden',
    flexGrow: 1,
  },
  card: {
    display: 'flex',
    flexDirection: 'column',
  },
  input: {
    paddingBottom: '24px',
  },
  gridInput: {
    display: 'inline-flex',
  },
  nameHeader: {
    display: 'flex',
    alignItems: 'center',
  },
  breadcrumbs: {
    flexGrow: 1,
    marginBottom: '16px',
  },
  separator: {
    borderBottom: `1px solid ${theme.palette.grey[100]}`,
    margin: '0 0 24px -24px',
    paddingBottom: '24px',
    width: 'calc(100% + 48px)',
  },
  cancelButton: {
    marginRight: '8px',
  },
});

type State = {
  project: ?Project,
  locationId: ?string,
};

const addProjectCard__projectTypeQuery = graphql`
  query AddProjectCard__projectTypeQuery($projectTypeId: ID!) {
    projectType: node(id: $projectTypeId) {
      ... on ProjectType {
        id
        name
        description
        properties {
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
          isInstanceProperty
          isDeleted
          isMandatory
        }
      }
    }
  }
`;

class AddProjectCard extends React.Component<Props, State> {
  state = {
    locationId: null,
    project: null,
  };
  componentDidMount() {
    this.getProject().then(project => {
      this.setState({
        project,
      });
    });
  }
  render() {
    const {classes} = this.props;
    const {project} = this.state;
    if (!project) {
      return (
        <div className={classes.root}>
          <CircularProgress />
        </div>
      );
    }
    const {properties} = project;
    return (
      <div className={classes.root}>
        <FormValidationContextProvider>
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
                    id: `new_project_` + Date.now(),
                    name: 'New Project',
                  },
                ]}
                size="large"
              />
            </div>
            <FormSaveCancelPanel
              onCancel={() => this.props.history.push(this.props.match.url)}
              onSave={this._saveProject}
            />
          </div>
          <div className={classes.contentRoot}>
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
                    <div className={classes.separator} />
                    <Grid container spacing={2}>
                      {project.type && (
                        <Grid item xs={12} sm={6} lg={4} xl={4}>
                          <FormField label="Type">
                            <TextField
                              disabled
                              variant="outlined"
                              margin="dense"
                              className={classes.gridInput}
                              value={project.type.name}
                            />
                          </FormField>
                        </Grid>
                      )}
                      <Grid item xs={12} sm={6} lg={4} xl={4}>
                        <FormField label="Location">
                          <LocationTypeahead
                            className={classes.gridInput}
                            headline={null}
                            margin="dense"
                            onLocationSelection={location =>
                              this._locationChangeHandler(location?.id ?? null)
                            }
                          />
                        </FormField>
                      </Grid>
                      {properties &&
                        properties
                          .filter(property => !property.propertyType.isDeleted)
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
                  </ExpandingPanel>
                </Grid>
                <Grid item xs={4} sm={4} lg={4} xl={4}>
                  <ExpandingPanel title="Team">
                    <UserTypeahead
                      className={classes.input}
                      headline="Owner"
                      onUserSelection={user =>
                        this._setProjectDetail('creatorId', user?.id)
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
  }

  async getProject(): Promise<Project> {
    const response = await fetchQuery(
      RelayEnvironment,
      addProjectCard__projectTypeQuery,
      {
        projectTypeId: this.props.projectTypeId,
      },
    );
    const projectType = response.projectType;

    let initialProps = [];
    if (projectType.properties) {
      initialProps = projectType.properties
        .filter(propertyType => !propertyType.isDeleted)
        .map(propType => getInitialPropertyFromType(propType));
      initialProps = initialProps.sort(sortPropertiesByIndex);
    }

    return {
      id: `project@tmp-${Date.now()}`,
      type: projectType,
      name: projectType.name,
      description: projectType.description,
      creatorId: null,
      location: null,
      properties: initialProps,
      workOrders: [],
      numberOfWorkOrders: 0,
    };
  }

  _saveProject = () => {
    const {name, description, creatorId, properties, type} = nullthrows(
      this.state.project,
    );
    const variables: AddProjectMutationVariables = {
      input: {
        name,
        description,
        creatorId: creatorId,
        type: type?.id ?? '',
        properties: toPropertyInput(properties),
        location: this.state.locationId,
      },
    };
    const callbacks: MutationCallbacks<AddProjectMutationResponse> = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          this._enqueueError(errors[0].message);
        } else {
          // navigate to main page
          this.props.history.push(this.props.match.url);
        }
      },
      onError: () => {
        this._enqueueError('Error saving work order');
      },
    };
    ServerLogger.info(LogEvents.SAVE_PROJECT_BUTTON_CLICKED, {
      source: 'project_details',
    });
    AddProjectMutation(variables, callbacks);
  };

  _enqueueError = (message: string) => {
    this.props.enqueueSnackbar(message, {
      children: key => (
        <SnackbarItem id={key} message={message} variant="error" />
      ),
    });
  };

  _setProjectDetail = (key: 'name' | 'description' | 'creatorId', value) => {
    this.setState(prevState => {
      return {
        // $FlowFixMe Set state for each field
        project: update(prevState.project, {[key]: {$set: value}}),
      };
    });
  };

  _propertyChangedHandler = index => property =>
    this.setState(prevState => {
      return {
        project: update(prevState.project, {
          properties: {[index]: {$set: property}},
        }),
      };
    });

  _locationChangeHandler = locationId => this.setState({locationId});

  navigateToMainPage = () => {
    ServerLogger.info(LogEvents.WORK_ORDERS_SEARCH_NAV_CLICKED, {
      source: 'work_order_details',
    });
    const {match} = this.props;
    this.props.history.push(match.url);
  };
}

export default withSnackbar(withRouter(withStyles(styles)(AddProjectCard)));
