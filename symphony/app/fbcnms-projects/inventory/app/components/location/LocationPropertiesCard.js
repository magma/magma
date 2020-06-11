/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {AppContextType} from '@fbcnms/ui/context/AppContext';
import type {Equipment} from '../../common/Equipment';
import type {LocationMenu_location} from './__generated__/LocationMenu_location.graphql';
import type {WithSnackbarProps} from 'notistack';
import type {WithStyles} from '@material-ui/core';

import AppContext from '@fbcnms/ui/context/AppContext';
import Button from '@fbcnms/ui/components/design-system/Button';
import Card from '@material-ui/core/Card';
import ErrorMessage from '@fbcnms/ui/components/ErrorMessage';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import InventoryQueryRenderer from '../../components/InventoryQueryRenderer';
import LocationBreadcrumbsTitle from './LocationBreadcrumbsTitle';
import LocationCoverageMapTab from './LocationCoverageMapTab';
import LocationDetailsTab from './LocationDetailsTab';
import LocationDocumentsCard from './LocationDocumentsCard';
import LocationFloorPlansTab from './LocationFloorPlansTab';
import LocationMenu from './LocationMenu';
import LocationNetworkMapTab from './LocationNetworkMapTab';
import LocationSiteSurveyTab from './LocationSiteSurveyTab';
import React from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import fbt from 'fbt';
import {FormContextProvider} from '../../common/FormContext';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {graphql} from 'react-relay';
import {withSnackbar} from 'notistack';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  locationId: ?string,
  selectedWorkOrderId: ?string,
  onEquipmentSelected: Equipment => void,
  onWorkOrderSelected: (workOrderId: string) => void,
  onEdit: () => void,
  onAddEquipment: () => void,
  onLocationMoved: (movedLocation: LocationMenu_location) => void,
  onLocationRemoved: (removedLocation: LocationMenu_location) => void,
} & WithStyles<typeof styles> &
  WithSnackbarProps;

type State = {
  selectedTab: 'details' | 'documents',
  isLoadingDocument: boolean,
};

const styles = theme => ({
  root: {
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
  },
  contentRoot: {
    position: 'relative',
    flexGrow: 1,
    overflow: 'auto',
  },
  tabs: {
    backgroundColor: 'white',
  },
  titleText: {
    fontWeight: 500,
  },
  section: {
    marginBottom: theme.spacing(3),
  },
  tabContainer: {
    width: 'auto',
  },
  locationNameHeader: {
    display: 'flex',
    alignItems: 'center',
    padding: '0px 8px',
    marginBottom: '16px',
    '&>:not(:last-child)': {
      marginRight: '8px',
    },
  },
  locationType: {
    fontSize: theme.typography.pxToRem(11),
    color: theme.palette.text.secondary,
    marginLeft: '6px',
    alignSelf: 'flex-end',
    lineHeight: 2.2,
    flexGrow: 1,
  },
  cardContentRoot: {
    '&:last-child': {
      paddingBottom: '0px',
    },
  },
  footer: {
    padding: '12px 16px',
    boxShadow: '0px -1px 4px rgba(0, 0, 0, 0.1)',
  },
  tabsContainer: {
    marginBottom: '16px',
  },
});

const locationsPropertiesCardQuery = graphql`
  query LocationPropertiesCardQuery($locationId: ID!) {
    location: node(id: $locationId) {
      ... on Location {
        id
        name
        latitude
        longitude
        externalId
        locationType {
          id
          name
          mapType
          mapZoomLevel
          propertyTypes {
            ...PropertyTypeFormField_propertyType
            ...DynamicPropertiesGrid_propertyTypes
          }
        }
        ...LocationBreadcrumbsTitle_locationDetails
        parentLocation {
          id
        }
        children {
          id
        }
        equipments {
          ...EquipmentTable_equipments
        }
        properties {
          ...PropertyFormField_property
          ...DynamicPropertiesGrid_properties
        }
        images {
          id
        }
        files {
          id
        }
        hyperlinks {
          id
        }
        surveys {
          id
        }
        ...LocationSiteSurveyTab_location
        ...LocationDocumentsCard_location
        ...LocationFloorPlansTab_location
        ...LocationMenu_location
      }
    }
  }
`;

class LocationPropertiesCard extends React.Component<Props, State> {
  state = {
    selectedTab: 'details',
    isLoadingDocument: false,
  };

  static contextType = AppContext;
  context: AppContextType;

  render() {
    const {
      classes,
      locationId,
      onLocationMoved,
      onLocationRemoved,
      onAddEquipment,
    } = this.props;
    if (!locationId) {
      return null;
    }

    const networkTopologyEnabled = this.context.isFeatureEnabled(
      'network_topology',
    );
    const siteSurveyEnabled = this.context.isFeatureEnabled('site_survey');
    const floorPlansEnabled = this.context.isFeatureEnabled('floor_plans');
    const coverageMapEnabled = this.context.isFeatureEnabled('coverage_maps');

    return (
      <InventoryQueryRenderer
        query={locationsPropertiesCardQuery}
        variables={{
          locationId: locationId,
        }}
        render={props => {
          const location = props.location;
          if (!location) {
            return (
              <Card className={classes.root}>
                <ErrorMessage message="It appears this location does not exist" />
              </Card>
            );
          }

          return (
            <FormContextProvider
              permissions={{
                entity: 'location',
                action: 'update',
              }}>
              <div className={classes.root}>
                <div className={classes.cardHeader}>
                  <div className={classes.locationNameHeader}>
                    <LocationBreadcrumbsTitle
                      locationDetails={location}
                      hideTypes={false}
                    />
                    <FormAction>
                      <LocationMenu
                        location={location}
                        popoverMenuClassName={classes.popoverMenu}
                        onLocationMoved={onLocationMoved}
                        onLocationRemoved={onLocationRemoved}
                      />
                    </FormAction>
                    <FormAction>
                      <Button onClick={this.props.onEdit}>
                        <fbt desc="">Edit Location</fbt>
                      </Button>
                    </FormAction>
                  </div>
                  <div className={classes.tabsContainer}>
                    <Tabs
                      className={classes.tabs}
                      value={this.state.selectedTab}
                      onChange={(_e, selectedTab) => {
                        ServerLogger.info(LogEvents.LOCATION_CARD_TAB_CLICKED, {
                          tab: selectedTab,
                        });
                        this.setState({selectedTab});
                      }}
                      indicatorColor="primary"
                      textColor="primary">
                      <Tab
                        classes={{root: classes.tabContainer}}
                        label="Details"
                        value="details"
                      />
                      <Tab
                        classes={{root: classes.tabContainer}}
                        label="Documents"
                        value="documents"
                      />
                      {networkTopologyEnabled && (
                        <Tab
                          classes={{root: classes.tabContainer}}
                          label="Network Map"
                          value="network_map"
                        />
                      )}
                      {siteSurveyEnabled && (
                        <Tab
                          classes={{root: classes.tabContainer}}
                          label="Site Surveys"
                          value="site_survey"
                        />
                      )}
                      {coverageMapEnabled && (
                        <Tab
                          classes={{root: classes.tabContainer}}
                          label="Coverage Maps"
                          value="coverage_map"
                        />
                      )}
                      {floorPlansEnabled && (
                        <Tab
                          classes={{root: classes.tabContainer}}
                          label="Floor Plans"
                          value="floor_plans"
                        />
                      )}
                    </Tabs>
                  </div>
                </div>
                <div className={classes.contentRoot}>
                  {this.state.selectedTab === 'details' ? (
                    <LocationDetailsTab
                      location={location}
                      selectedWorkOrderId={this.props.selectedWorkOrderId}
                      onEquipmentSelected={this.props.onEquipmentSelected}
                      onWorkOrderSelected={this.props.onWorkOrderSelected}
                      onAddEquipment={onAddEquipment}
                    />
                  ) : null}
                  {this.state.selectedTab === 'documents' ? (
                    <LocationDocumentsCard location={location} />
                  ) : null}
                  {this.state.selectedTab === 'network_map' ? (
                    <LocationNetworkMapTab locationId={location.id} />
                  ) : null}
                  {this.state.selectedTab === 'site_survey' ? (
                    <LocationSiteSurveyTab location={location} />
                  ) : null}
                  {this.state.selectedTab === 'coverage_map' ? (
                    <LocationCoverageMapTab location={location} />
                  ) : null}
                  {this.state.selectedTab === 'floor_plans' && (
                    <LocationFloorPlansTab location={location} />
                  )}
                </div>
              </div>
            </FormContextProvider>
          );
        }}
      />
    );
  }
}

export default withStyles(styles)(withSnackbar(LocationPropertiesCard));
