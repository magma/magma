/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {LocationType} from '../common/LocationType';
import type {WithStyles} from '@material-ui/core';

import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import FormField from '@fbcnms/ui/components/FormField';
import React from 'react';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  locationType: ?LocationType,
} & WithStyles<typeof styles>;

const styles = {
  root: {
    alignItems: 'flex-start',
    display: 'flex',
    minHeight: '250px',
    justifyContent: 'space-around',
  },
  field: {
    marginLeft: 5,
  },
};

class LocationTypePropertiesCard extends React.Component<Props> {
  render() {
    const {locationType} = this.props;
    return (
      <Card className={this.props.classes.root}>
        <CardContent>
          <FormField
            className={this.props.classes.field}
            label="Name"
            value={locationType?.name}
          />
          <FormField
            className={this.props.classes.field}
            label="Inherited Properties:"
          />
          <FormField
            className={this.props.classes.field}
            label="Add New Property:"
          />
        </CardContent>
      </Card>
    );
  }
}

export default withStyles(styles)(LocationTypePropertiesCard);
