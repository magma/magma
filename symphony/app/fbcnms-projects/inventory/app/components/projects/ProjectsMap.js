/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {ProjectMapMarkerData} from '../map/ProjectsMapUtils';

import * as React from 'react';
import ProjectsMapView from './../map/ProjectsMapView';
import {createFragmentContainer, graphql} from 'react-relay';
import {projectToGeoJson} from './../map/ProjectsMapUtils';
import {withRouter} from 'react-router-dom';

type Props = {
  projects: Array<ProjectMapMarkerData>,
};

const ProjectsMap = (props: Props) => {
  const {projects} = props;

  return (
    <ProjectsMapView
      mode="streets"
      showMapSatelliteToggle={true}
      showGeocoder={true}
      markers={projectToGeoJson(projects.filter(w => w.location != null))}
    />
  );
};

export default withRouter(
  createFragmentContainer(ProjectsMap, {
    projects: graphql`
      fragment ProjectsMap_projects on Project @relay(plural: true) {
        id
        name
        location {
          id
          name
          latitude
          longitude
        }
        numberOfWorkOrders
      }
    `,
  }),
);
