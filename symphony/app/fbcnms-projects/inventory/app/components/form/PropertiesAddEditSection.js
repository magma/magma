/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {Property} from '../../common/Property';
import type {PropertyType} from '../../common/PropertyType';
import type {WithStyles} from '@material-ui/core';

import Grid from '@material-ui/core/Grid';
import PropertyValueInput from './PropertyValueInput';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import {withStyles} from '@material-ui/core/styles';

type Props = WithStyles<typeof styles> & {
  properties: Array<Property>,
  onChange: (propertyIndex: number) => (Property | PropertyType) => void,
};

const styles = theme => ({
  subheader: {
    marginTop: theme.spacing(2),
    marginBottom: theme.spacing(2),
  },
  input: {
    display: 'inline-flex',
    width: '100%',
  },
  titleText: {
    fontSize: '20px',
    lineHeight: '24px',
    fontWeight: 500,
  },
});

class PropertiesAddEditSection extends React.Component<Props> {
  render() {
    const {classes, properties} = this.props;
    return (
      <>
        <div className={classes.subheader}>
          <Text variant="subtitle1" className={classes.titleText}>
            Properties
          </Text>
        </div>
        <Grid container spacing={2}>
          {properties.map((property, index) => (
            <Grid key={property.id} item xs={12} sm={12} lg={6} xl={4}>
              <PropertyValueInput
                hasSpacer={true}
                required={!!property.propertyType.isMandatory}
                disabled={!property.propertyType.isInstanceProperty}
                label={property.propertyType.name}
                className={classes.input}
                margin="none"
                inputType="Property"
                property={property}
                onChange={this.props.onChange(index)}
                headlineVariant="form"
              />
            </Grid>
          ))}
        </Grid>
      </>
    );
  }
}

export default withStyles(styles)(PropertiesAddEditSection);
