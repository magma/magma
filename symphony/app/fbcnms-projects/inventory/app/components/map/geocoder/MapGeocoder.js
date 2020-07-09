/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @format
 * @flow
 */

'use strict';

import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import MapboxGeocoder from '@fbcnms/ui/components/map/MapboxGeocoder';
import React from 'react';
import mapboxgl from 'mapbox-gl';

import type {LngLatLike} from 'mapbox-gl/src/geo/lng_lat';
import type {
  ProjectGeoJSONFeatureCollection,
  ProjectLocation,
} from '../ProjectsMapUtils';

export type Feature = {
  title: string,
  properties: {
    name: string,
    numberOfWorkOrders: number,
    location: ProjectLocation,
  },
  id: number | string,
  center?: LngLatLike,
  bbox?: Array<LngLatLike>,
  geometry: {type: string, coordinates: LngLatLike},
};
export type Result = {feature: Feature};

type Props = {
  markers?: ?ProjectGeoJSONFeatureCollection,
  accessToken: string,
  mapRef: ?mapboxgl.Map,
  featuresType: 'Work Order' | 'Project' | '',
  headLine: string,
  onSelectFeature: Feature => void,
  // Mapbox geocoding API: https://www.mapbox.com/api-documentation/#geocoding
  apiEndpoint: string,
  // Debounce searches at this interval
  searchDebounceMs: number,
  shouldSearchPlaces?: ?(customResults: Array<Result>) => boolean,
};

class MapGeocoder extends React.Component<Props> {
  static defaultProps = {
    apiEndpoint: 'https://api.mapbox.com/geocoding/v5/mapbox.places/',
    searchDebounceMs: 200,
    headLine: '',
  };

  _getCustomResults = (originalQuery: string) => {
    const markers = this.props.markers;
    if (!markers || markers === null) {
      return {resultsType: '', results: []};
    }
    const query = originalQuery.toLowerCase();
    const matches = [];
    markers.features.forEach(feature => {
      if (String(feature.properties?.name).toLowerCase().includes(query)) {
        matches.push({feature});
      }
    });
    return {resultsType: this.props.headLine, results: matches};
  };

  _onRenderResult = (result: Result, handleClearInput: () => void) => {
    if (!result.hasOwnProperty('feature')) {
      return null;
    }
    const primaryText = <span>{result.feature.properties.name}</span>;
    const secondaryText = <span>{this.props.featuresType}</span>;
    return (
      <ListItem
        key={`${this.props.featuresType}-${result.feature.id}`}
        button
        dense
        onClick={() => {
          this.props.onSelectFeature({
            ...result.feature,
            center: result.feature.geometry.coordinates,
          });
          //  Clear the search field
          handleClearInput();
        }}>
        <ListItemText primary={primaryText} secondary={secondaryText} />
      </ListItem>
    );
  };

  render() {
    const {accessToken, mapRef, onSelectFeature} = this.props;
    return (
      <MapboxGeocoder
        accessToken={accessToken}
        mapRef={mapRef}
        onSelectFeature={onSelectFeature}
        getCustomResults={this._getCustomResults}
        shouldSearchPlaces={() => {
          return true;
        }}
        onRenderResult={this._onRenderResult}
      />
    );
  }
}

export default MapGeocoder;
