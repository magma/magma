/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {MapMarkerProps} from './map/MapTypes';

import GatewayMapMarkerPopup from './GatewayMapMarkerPopup';
import PlaceIcon from '@material-ui/icons/Place';
import React from 'react';
import ReactDOM from 'react-dom';
import mapboxgl from 'mapbox-gl';

type State = {
  marker: ?mapboxgl.Marker,
};

class GatewayMapMarker extends React.Component<MapMarkerProps, State> {
  state = {
    marker: null,
  };

  componentDidMount() {
    this.addMarker();
  }

  shouldComponentUpdate() {
    // Let mapboxgl control this component.  Once it's
    // added the lifecycle belongs to mapbox
    return false;
  }

  _featureLngLat() {
    let lng = this.props.feature.geometry.coordinates[0];
    let lat = this.props.feature.geometry.coordinates[1];
    lng = Math.abs(lng) > 180 ? lng % 180 : lng;
    lat = Math.abs(lat) > 90 ? lat % 90 : lat;
    return [lng, lat];
  }

  addMarker() {
    const {feature, map} = this.props;

    const markerEl = document.createElement('div');
    ReactDOM.render(this.renderMarker(), markerEl);

    const popupEl = document.createElement('div');
    ReactDOM.render(
      <GatewayMapMarkerPopup gateway={feature.properties.gateway} />,
      popupEl,
    );

    const popup = new mapboxgl.Popup({offset: 50}).setDOMContent(popupEl);

    const marker = new mapboxgl.Marker({
      element: markerEl,
      offset: [0, -30],
    })
      .setLngLat(this._featureLngLat())
      .setPopup(popup)
      .addTo(map);

    this.setState({marker});
  }

  renderMarker() {
    const {feature} = this.props;
    const {id} = feature.properties;

    return (
      <div position={this._featureLngLat()} key={id}>
        <PlaceIcon color="primary" style={{fontSize: 48}} />
      </div>
    );
  }

  render() {
    // render nothing here, since it's handled in mapboxgl
    return null;
  }
}

export default GatewayMapMarker;
