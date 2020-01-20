/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {LngLatLike} from 'mapbox-gl/src/geo/lng_lat';
import type {MapType} from '@fbcnms/magmalte/app/components/map/styles';
import type {ProjectGeoJSONFeatureCollection} from './ProjectsMapUtils';
import type {WithStyles} from '@material-ui/core';

import 'mapbox-gl/dist/mapbox-gl.css';
import * as React from 'react';
import Avatar from '@material-ui/core/Avatar';
import Chip from '@material-ui/core/Chip';
import MapButtonGroup from '@fbcnms/ui/components/map/MapButtonGroup';
import MapGeocoder from './geocoder/MapGeocoder';
import ProjectsPopover from '../projects/ProjectsPopover';
import ReactDOM from 'react-dom';
import mapboxgl from 'mapbox-gl';
import nullthrows from '@fbcnms/util/nullthrows';
import {black, blue70, white} from '@fbcnms/ui/theme/colors';
import {getMapStyleForType} from '@fbcnms/magmalte/app/components/map/styles';
import {withStyles} from '@material-ui/core/styles';

const styles = () => ({
  mapContainer: {
    height: '100%',
    width: '100%',
  },
  buttonGroupContainer: {
    position: 'absolute',
    left: 0,
    bottom: 0,
  },
  circle: {
    radius: '15.5px',
    height: '20px',
    width: '20px',
    fontSize: '14px',
    marginLeft: '5px',
    borderRadius: '50%',
    backgroundColor: blue70,
    fontFamily: 'Roboto-Bold',
  },
  chip: {
    borderRadius: '4px',
    color: black,
    opacity: 0.95,
    backgroundColor: white,
    '&:hover, &:focus': {
      backgroundColor: blue70,
      color: white,
      '& $circle': {
        backgroundColor: white,
        color: blue70,
      },
    },
  },
});

type State = {
  map: ?mapboxgl.Map,
  style: ?MapType,
  container: ?HTMLDivElement,
  projectId: ?string,
};

type Props = {
  mode: MapType,
  zoomLevel?: string,
  center?: LngLatLike,
  markers?: ?ProjectGeoJSONFeatureCollection,
  showGeocoder?: boolean,
  showMapSatelliteToggle?: boolean,
  mapButton?: React.Node,
} & WithStyles<typeof styles>;

// https://docs.mapbox.com/mapbox-gl-js/style-spec/#expressions
type MapExpression = Array<string | number | MapExpression>;

// https://docs.mapbox.com/mapbox-gl-js/style-spec/#function-type
type PaintType = 'identity' | 'exponential' | 'interval' | 'categorical';

export type ColorStop = {
  threshold: number,
  color: string,
};

export type CircleColorInterpolation = {
  property: string,
  type: PaintType,
  stops: Array<ColorStop>,
};

class ProjectsMapView extends React.Component<Props, State> {
  static defaultProps = {
    markers: null,
    center: [0, 0],
    zoomLevel: '2',
  };

  state = {
    map: null,
    style: this.props.mode,
    projectId: null,
    container: null,
  };

  mapContainer = null;

  componentDidMount() {
    this.initMap();
  }

  componentWillUnmount() {
    this.state.container &&
      ReactDOM.unmountComponentAtNode(this.state.container);
  }

  _fixCoordinates(coordinates) {
    const lng =
      Math.abs(coordinates[0]) > 180 ? coordinates[0] % 180 : coordinates[0];
    const lat =
      Math.abs(coordinates[1]) > 90 ? coordinates[1] % 90 : coordinates[1];
    return [lng, lat];
  }

  _fitBounds = () => {
    const {map} = this.state;
    const markers = this.props.markers;

    if (!markers) {
      return;
    }

    if (!map || !markers || markers.features.length == 0) {
      return;
    }
    const bounds = new mapboxgl.LngLatBounds();

    markers.features.map(feature => {
      const geometry = nullthrows(feature.geometry);
      if (geometry.type !== 'Point') {
        return;
      }
      const coords = mapboxgl.LngLat.convert(
        this._fixCoordinates(geometry.coordinates),
      );
      bounds.extend(coords);
    });

    if (!bounds.isEmpty()) {
      map.fitBounds(bounds, {
        padding: {top: 50, bottom: 50, left: 50, right: 50},
        easing: t => t * (2 - t),
        duration: 0,
        maxZoom: 19, // 19 = ~city block
      });
    }
  };

  onClickOutside = () => {
    this.setState({
      projectId: null,
    });
  };

  _onProjectMarkerClick = selectedFeatureId => {
    this.setState({
      projectId: selectedFeatureId,
    });
  };

  initMap() {
    const map = new mapboxgl.Map({
      attributionControl: true,
      container: this.mapContainer,
      hash: false,
      style: getMapStyleForType(this.props.mode),
      zoom: this.props.zoomLevel,
      center: this.props.center,
    });

    map.on('style.load', () => {
      this._addMarkers();
    });

    map.on('click', () => {
      this.onClickOutside();
    });

    map.addControl(new mapboxgl.NavigationControl({}));
    this.setState({map}, this._fitBounds);
  }

  render() {
    const {classes, mapButton} = this.props;
    const {map} = this.state;

    return (
      <div
        ref={e => {
          this.mapContainer = e;
          map && map.resize();
        }}
        className={classes.mapContainer}>
        {map && mapboxgl.accessToken && (
          <>
            {this.props.showGeocoder && (
              <MapGeocoder
                accessToken={mapboxgl.accessToken}
                mapRef={map}
                onSelectFeature={this._onGeocoderEvent}
                markers={this.props.markers}
                featuresType={'Project'}
                headLine={'Projects'}
              />
            )}
            <div>
              {this.state.projectId !== null && (
                <ProjectsPopover projectId={this.state.projectId} />
              )}
            </div>
            <>
              {this.props.showMapSatelliteToggle && (
                <div className={classes.buttonGroupContainer}>
                  <MapButtonGroup
                    initiallySelectedButton={0}
                    onIconClicked={id => {
                      (id === 'satellite' || id === 'streets') &&
                        this._onIconButtonEvent(id);
                    }}
                    buttons={[
                      {item: 'Map', id: 'streets'},
                      {item: 'Satellite', id: 'satellite'},
                    ]}
                  />
                  {mapButton}
                </div>
              )}
            </>
          </>
        )}
      </div>
    );
  }

  _onGeocoderEvent = feature => {
    // Move to a location returned by the geocoder
    const {map} = this.state;
    if (map) {
      const {center} = feature;
      map.flyTo({center, zoom: 19});
    }
  };

  _onIconButtonEvent = (id: MapType) => {
    const {map} = this.state;
    if (map && this.state.style != id) {
      map.setStyle(getMapStyleForType(id));
      this.setState({style: id === 'streets' ? 'streets' : 'satellite'});
    }
  };

  _addMarkers = () => {
    const {classes} = this.props;
    const markers = this.props.markers;
    if (!markers) {
      return;
    }
    const map = nullthrows(this.state.map);
    markers.features.forEach(feature => {
      const geometry = nullthrows(feature.geometry);
      const selectedFeatureId = feature.properties?.id;
      if (geometry.type === 'Point') {
        const marker = new mapboxgl.Marker(<div />)
          .setLngLat(geometry.coordinates)
          .addTo(map);
        ReactDOM.render(
          <Chip
            avatar={
              <Avatar className={classes.circle}>
                {String(feature.properties?.numberOfWorkOrders)}
              </Avatar>
            }
            label={feature.properties?.name}
            id={selectedFeatureId}
            clickable
            className={classes.chip}
            color="primary"
            onClick={e => {
              this._onProjectMarkerClick(selectedFeatureId);
              e.stopPropagation();
            }}
          />,
          marker.getElement(),
        );
      }
    });
  };
}

export default withStyles(styles)(ProjectsMapView);
