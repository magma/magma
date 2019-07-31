/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @format
 */

'use strict';

import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import Paper from '@material-ui/core/Paper';
import PropTypes from 'prop-types';
import React from 'react';
import SearchBar from './SearchBar';
import axios from 'axios';
import {withStyles} from '@material-ui/core/styles';

const styles = {
  root: {
    position: 'absolute',
    width: 240,
    top: 10,
    left: 10,
    zIndex: 10,
  },
  resultsPaper: {
    marginTop: 4,
  },
};

class MapboxGeocoder extends React.Component {
  state = {
    value: '',
    isLoading: false,
    results: [], // {feature: obj}, or custom structures via getCustomResults()
  };

  getResults = query => {
    // Fetch results for the given query
    const {getCustomResults, shouldSearchPlaces} = this.props;

    // Fetch any custom results first
    let results = [];
    if (getCustomResults) {
      results = getCustomResults(query);
      if (shouldSearchPlaces && !shouldSearchPlaces(results)) {
        // Don't search for default place results?
        this.setState({results});
        return;
      }
    }

    // Fetch default place results (if needed)
    this.setState({results, isLoading: true}, () =>
      this.mapboxPlacesSearch(query),
    );
  };

  mapboxPlacesSearch = query => {
    // Send an API request for the given query
    const {apiEndpoint, accessToken} = this.props;

    // Construct GET request
    // See: https://www.mapbox.com/api-documentation/#search-for-places
    const params = {
      access_token: accessToken,
      ...this.getProximity(),
    };
    const encodedParams = Object.entries(params)
      .map(kv => kv.map(encodeURIComponent).join('='))
      .join('&');
    const uri =
      apiEndpoint + encodeURIComponent(query) + '.json?' + encodedParams;

    // Send request
    axios
      .get(uri)
      .then(response => {
        // Store the results
        const {features} = response.data;
        if (features) {
          this.setState({
            results: [
              ...this.state.results,
              ...features.map(feature => ({feature})),
            ],
            isLoading: false,
          });
        }
      })
      .catch(_err => {
        // TODO handle this better
        this.setState({results: [], isLoading: false});
      });
  };

  getProximity() {
    // Return proximity arguments based on the current map center and zoom level
    // (or none if not applicable)
    const {mapRef} = this.props;
    if (mapRef && mapRef.getZoom() > 9) {
      const center = mapRef.getCenter().wrap();
      return {proximity: [center.lng, center.lat].join(',')};
    }
    return {};
  }

  handleInput = e => {
    // Handle an input value update
    this.setState({value: e.target.value});
  };

  handleClearInput = () => {
    // Clear the current input and results
    this.setState({value: '', results: [], isLoading: false});
  };

  renderResult = result => {
    // Render a single result
    const {onSelectFeature, onRenderResult} = this.props;

    // Use custom renderer (if applicable)
    if (onRenderResult) {
      const listItem = onRenderResult(result, this.handleClearInput.bind(this));
      if (listItem !== null) {
        return listItem;
      }
    }

    // Render feature
    if (!result.hasOwnProperty('feature')) {
      return null; // shouldn't happen (unhandled result in getCustomResults)
    }
    const {feature} = result;
    const primaryText = feature.text;
    let secondaryText =
      (feature.properties && feature.properties.address) || feature.place_name;
    if (secondaryText === primaryText) {
      secondaryText = undefined; // don't show duplicate text
    }

    return (
      <ListItem
        key={'feature-' + feature.id}
        button
        dense
        onClick={() => {
          // Selected a map feature
          onSelectFeature(result.feature);

          // Clear the search field
          this.handleClearInput();
        }}>
        <ListItemText primary={primaryText} secondary={secondaryText} />
      </ListItem>
    );
  };

  render() {
    const {classes, searchDebounceMs} = this.props;
    const {value, isLoading, results} = this.state;

    return (
      <div className={classes.root}>
        <SearchBar
          value={value}
          onChange={this.handleInput}
          onClearInput={this.handleClearInput}
          onSearch={this.getResults}
          isLoading={isLoading}
          debounceMs={searchDebounceMs}
        />

        {results.length > 0 ? (
          <Paper className={classes.resultsPaper} elevation={2}>
            <List component="nav">
              {results.map(result => this.renderResult(result))}
            </List>
          </Paper>
        ) : null}
      </div>
    );
  }
}

MapboxGeocoder.propTypes = {
  classes: PropTypes.object.isRequired,
  accessToken: PropTypes.string.isRequired,
  mapRef: PropTypes.object,
  onSelectFeature: PropTypes.func.isRequired,

  // Mapbox geocoding API: https://www.mapbox.com/api-documentation/#geocoding
  apiEndpoint: PropTypes.string,

  // Debounce searches at this interval
  searchDebounceMs: PropTypes.number,

  // (query : str) => results : arr of obj
  getCustomResults: PropTypes.func,

  // If getCustomResults is defined, should we search for default places?
  // (customResults : arr of obj) => bool
  shouldSearchPlaces: PropTypes.func,

  // (result : obj, handleClearInput : func) => <ListItem> or null
  onRenderResult: PropTypes.func,
};

MapboxGeocoder.defaultProps = {
  apiEndpoint: 'https://api.mapbox.com/geocoding/v5/mapbox.places/',
  searchDebounceMs: 200,
};

export default withStyles(styles)(MapboxGeocoder);
