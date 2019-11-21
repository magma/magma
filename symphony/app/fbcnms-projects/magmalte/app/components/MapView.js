/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ComponentType} from 'react';
import type {
  FilterSpecification,
  LayerSpecification,
} from 'mapbox-gl/src/style-spec/types';
import type {
  MagmaConnectionFeature,
  MagmaFeatureCollection,
} from '../common/GeoJSON';
import type {MapMarkerProps} from './MapTypes';
import type {MapMouseEvent} from 'mapbox-gl/src/ui/events';
import type {WithStyles} from '@material-ui/core';

import React from 'react';
import mapboxgl from 'mapbox-gl';
import {getDefaultMapStyle} from './map/styles';
import {isEqual} from 'lodash';
import {withStyles} from '@material-ui/core/styles';

const styles = {
  mapContainer: {
    height: '100%',
    width: '100%',
  },
};

type State = {
  map: ?mapboxgl.Map,
};

type Props = WithStyles<typeof styles> & {
  geojson: MagmaFeatureCollection,
  mapLayers?: Array<LayerSpecification>,
  MapMarker: ComponentType<MapMarkerProps>,
  onMarkerClick?: (string | number) => void,
  mapLayerFilters?: Map<string, FilterSpecification>,
  onClickFeatures?: (?Array<MagmaConnectionFeature>) => void,
  showMarkerLabels?: boolean,
  zoomHash?: string,
};

class MapView extends React.Component<Props, State> {
  state = {
    map: null,
  };

  mapContainer = null;

  componentDidUpdate(prevProps) {
    this._updateLayers(prevProps.mapLayers, prevProps.mapLayerFilters);

    // only zoom if transitioned from no markers to some markers (initial load)
    // or if zoom hash differs
    if (
      (prevProps.geojson.features.length === 0 &&
        this.props.geojson.features.length > 0) ||
      prevProps.zoomHash != this.props.zoomHash
    ) {
      this._fitBounds();
    }

    const map = this.state.map;
    if (map && prevProps.showMarkerLabels !== this.props.showMarkerLabels) {
      // we want to hide layers that render street signs and names to make the
      // map less cluttered when rendering marker labels
      map.style.stylesheet.layers.forEach(layer => {
        if (layer.type === 'symbol') {
          map.setLayoutProperty(
            layer.id,
            'visibility',
            this.props.showMarkerLabels ? 'none' : 'visible',
          );
        }
      });
    }
  }

  componentDidMount() {
    this.initMap();
  }

  initMap() {
    const map = new mapboxgl.Map({
      attributionControl: false,
      container: this.mapContainer,
      hash: false,
      style: getDefaultMapStyle(),
      zoom: 2,
      center: [0, 0],
    });

    map.addControl(
      new mapboxgl.AttributionControl({
        compact: true,
        customAttribution: mapboxgl.accessToken
          ? '' // Included by mapbox
          : '&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors',
      }),
    );
    map.addControl(new mapboxgl.NavigationControl({}));
    map.addControl(new mapboxgl.ScaleControl({position: 'bottom-left'}));

    this.setState({map});
  }

  _mapLoad = () => {
    this._registerListeners();
    this._fitBounds();
  };

  _updateLayers(
    prevMapLayers?: Array<LayerSpecification>,
    prevMapLayerFilters?: Map<string, FilterSpecification>,
  ): boolean {
    // returns true if displayed layers were updated

    const {mapLayers} = this.props;
    const {map} = this.state;

    if (!map) {
      return false;
    }

    if (!isEqual(prevMapLayers, mapLayers)) {
      if (prevMapLayers) {
        prevMapLayers.map(layer =>
          map.removeLayer(layer.id).removeSource(layer.id),
        );
      }
      if (mapLayers) {
        mapLayers.map(layer => {
          map.addLayer(layer);
        });
        // if we redraw layers, then the filters must be updated too
        this._updateLayerFilters();
      }
      return true;
    } else if (this._updateLayerFilters(prevMapLayerFilters)) {
      return true;
    }
    return false;
  }

  onMapClick = (event: MapMouseEvent) => {
    const {onClickFeatures, mapLayers} = this.props;
    const {map} = this.state;

    if (!map || !mapLayers || !onClickFeatures) {
      return;
    }

    const layers = mapLayers.map(l => l.id);
    if (layers.length > 0) {
      onClickFeatures(
        map.queryRenderedFeatures(event.point, {
          layers,
        }),
      );
    } else {
      onClickFeatures(null);
    }
  };

  _registerListeners() {
    const {map} = this.state;

    if (!map) {
      return;
    }

    this.props.onClickFeatures && map.on('click', this.onMapClick);

    map.on('mousemove', event => {
      if (this.props.mapLayers) {
        const features = map.queryRenderedFeatures(event.point, {
          layers: this.props.mapLayers.map(l => l.id),
        });
        // change cursor to pointer if hovered over a rendered feature/layer
        map.getCanvas().style.cursor = features.length ? 'pointer' : '';
      }
    });
  }

  _updateLayerFilters(
    prevLayerFilters: ?Map<string, FilterSpecification>,
  ): boolean {
    // assumes that mapLayerFilters.keys stays constant
    // returns true if a filter has been updated.
    const {map} = this.state;
    const {mapLayerFilters} = this.props;

    if (!map || !mapLayerFilters) {
      return false;
    }

    let hasUpdated = false;
    mapLayerFilters.forEach((filterSpec, layerName) => {
      if (
        map.getLayer(layerName) &&
        (!prevLayerFilters ||
          !isEqual(filterSpec, prevLayerFilters.get(layerName)))
      ) {
        map.setFilter(layerName, filterSpec);
        hasUpdated = true;
      }
    });
    return hasUpdated;
  }

  _fixCoordinates(coordinates) {
    const lng =
      Math.abs(coordinates[0]) > 180 ? coordinates[0] % 180 : coordinates[0];
    const lat =
      Math.abs(coordinates[1]) > 90 ? coordinates[1] % 90 : coordinates[1];
    return [lng, lat];
  }

  _fitBounds = () => {
    const {geojson} = this.props;
    const {map} = this.state;

    if (!map || geojson.features.length == 0) {
      return;
    }
    const bounds = new mapboxgl.LngLatBounds();

    geojson.features.map(feature => {
      const coords = mapboxgl.LngLat.convert(
        this._fixCoordinates(feature.geometry.coordinates),
      );
      bounds.extend(coords);
    });

    map.fitBounds(bounds, {
      padding: {top: 50, bottom: 50, left: 50, right: 50},
      easing: t => t * (2 - t),
      duration: 1000,
      maxZoom: 19, // 19 = ~city block
    });
  };

  render() {
    const {map} = this.state;
    const {classes, geojson, MapMarker} = this.props;

    let markers = [];
    if (map) {
      map.on('load', _e => this._mapLoad());
      markers = geojson.features.map(feature => (
        <MapMarker
          key={feature.properties.id}
          map={map}
          feature={feature}
          onClick={this.props.onMarkerClick}
          showLabel={this.props.showMarkerLabels}
        />
      ));
    }

    return (
      <>
        <div
          ref={e => (this.mapContainer = e)}
          className={classes.mapContainer}
        />
        {markers}
      </>
    );
  }
}

export default withStyles(styles)(MapView);
