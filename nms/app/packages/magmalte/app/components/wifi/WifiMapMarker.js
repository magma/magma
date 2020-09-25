/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow
 * @format
 */

import type {MapMarkerProps} from '@fbcnms/ui/insights/map/MapTypes';
import type {WithStyles} from '@material-ui/core';

import React from 'react';
import ReactDOM from 'react-dom';
import WifiMapMarkerIcon from './WifiMapMarkerIcon';
import WifiMapMarkerPopup from './WifiMapMarkerPopup';
import mapboxgl from 'mapbox-gl';
import nullthrows from '@fbcnms/util/nullthrows';
import {withStyles} from '@material-ui/core/styles';

const styles = {
  detailList: {
    margin: '0px',
  },
  markerContainer: {
    cursor: 'pointer',
  },
  deviceLabel: {
    position: 'absolute',
    marginTop: '-10px',
    whiteSpace: 'nowrap',
  },
  popupContent: {
    padding: '10px 10px 15px',
  },
  // we want to remove the padding of the MapBox popup container because it
  // intefers with the `mouseleave` event listener we install on the popup
  // element rendered inside the container
  '@global': {
    'div.mapboxgl-popup-content': {padding: 0},
  },
};

type State = {
  // $FlowFixMe[value-as-type] cannot import directly
  marker: mapboxgl.Marker,
  popupEl: HTMLDivElement,
  // $FlowFixMe[value-as-type] cannot import directly
  popup: mapboxgl.Popup,
};

class WifiMapMarker extends React.Component<
  MapMarkerProps & WithStyles<typeof styles>,
  State,
> {
  // create markerEl and marker so that state does not have to be updated
  // in componentDidMount and componentWillUnmount and componentDidUpdate
  state = {
    marker: new mapboxgl.Marker({
      element: document.createElement('div'),
      offset: [0, 0],
    }),
    popupEl: document.createElement('div'),
    popup: new mapboxgl.Popup({
      closeButton: false,
      closeOnClick: false,
      anchor: 'top-left',
    }),
  };

  componentDidMount() {
    this.updateMarker();
  }

  componentWillUnmount() {
    this.removeMarker();
  }

  shouldComponentUpdate(nextProps: MapMarkerProps) {
    const lastDevice = this.props.feature.properties.device;
    const nextDevice = nextProps.feature.properties.device;
    // if both are undefined
    if (!lastDevice && !nextDevice) {
      return false;
    }
    // if either is undefined
    if (!lastDevice || !nextDevice) {
      return true;
    }

    if (nextProps.showLabel !== this.props.showLabel) {
      return true;
    }

    if (
      this.props.feature.properties.useLargeIcon !==
      nextProps.feature.properties.useLargeIcon
    ) {
      return true;
    }

    // if checkinTime or readTime has changed
    return (
      lastDevice.checkinTime != nextDevice.checkinTime ||
      lastDevice.readTime != nextDevice.readTime
    );
  }

  componentDidUpdate() {
    this.removeMarker();
    this.updateMarker();
  }

  updateMarker() {
    const {feature, map, showLabel} = this.props;
    const {marker, popup, popupEl} = this.state;

    const {useLargeIcon} = feature.properties;

    const markerEl = marker.getElement();

    ReactDOM.render(
      <WifiMapMarkerIcon
        feature={feature}
        showLabel={showLabel}
        useLargeIcon={useLargeIcon}
      />,
      markerEl,
    );
    this.props.onClick && map.on('click', this.onClick);

    ReactDOM.render(<WifiMapMarkerPopup feature={feature} />, popupEl);
    popupEl.className = this.props.classes.popupContent;

    markerEl.addEventListener('mouseenter', () => {
      popup
        .setLngLat(feature.geometry.coordinates)
        .setDOMContent(popupEl)
        .addTo(map);
    });

    markerEl.addEventListener('mouseleave', (event: MouseEvent) => {
      const {relatedTarget} = event;
      if (
        relatedTarget !== popupEl &&
        (!relatedTarget || !(relatedTarget: any).contains(popupEl))
      ) {
        popup.remove();
      }
    });

    popupEl.addEventListener('mouseleave', (event: MouseEvent) => {
      if (event.relatedTarget !== markerEl) {
        popup.remove();
      }
    });

    this.state.marker.setLngLat(feature.geometry.coordinates).addTo(map);
  }

  onClick = event => {
    const targetElement = event.originalEvent.target;
    const markerEl = this.state.marker.getElement();
    if (targetElement === markerEl || markerEl.contains(targetElement)) {
      nullthrows(this.props.onClick)(this.props.feature.properties.id);
    }
  };

  removeMarker() {
    this.state.marker.remove();
    this.props.onClick && this.props.map.off('click', this.onClick);
  }

  render() {
    // render nothing here, since it's handled in mapboxgl
    return null;
  }
}

export default withStyles(styles)(WifiMapMarker);
