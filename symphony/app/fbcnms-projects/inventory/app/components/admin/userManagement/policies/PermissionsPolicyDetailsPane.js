/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {PermissionsPolicy} from '../data/PermissionsPolicies';

import * as React from 'react';
import Card from '@fbcnms/ui/components/design-system/Card/Card';
import FormFieldTextInput from '../utils/FormFieldTextInput';
import Grid from '@material-ui/core/Grid';
import ViewContainer from '@fbcnms/ui/components/design-system/View/ViewContainer';
import fbt from 'fbt';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  content: {
    paddingBottom: '16px',
  },
  nameField: {
    marginRight: '8px',
  },
  descriptionField: {
    marginTop: '16px',
  },
}));

type Props = $ReadOnly<{
  policy: PermissionsPolicy,
  onChange: PermissionsPolicy => void,
  className?: ?string,
}>;

export default function PermissionsPolicyDetailsPane(props: Props) {
  const {policy, className, onChange} = props;
  const classes = useStyles();

  return (
    <Card className={className} margins="none">
      <ViewContainer
        className={classes.content}
        header={{title: <fbt desc="">Policy Details</fbt>}}>
        <Grid container>
          <Grid item xs={12} sm={6} lg={6} xl={6}>
            <FormFieldTextInput
              className={classes.nameField}
              label={`${fbt('Policy Name', '')}`}
              validationId="name"
              value={policy.name}
              onValueChanged={name => {
                onChange({
                  ...policy,
                  name,
                });
              }}
            />
          </Grid>
          <Grid item xs={12}>
            <FormFieldTextInput
              type="multiline"
              className={classes.descriptionField}
              label={`${fbt('Policy Description', '')}`}
              value={policy.description || ''}
              onValueChanged={description => {
                onChange({
                  ...policy,
                  description,
                });
              }}
            />
          </Grid>
        </Grid>
      </ViewContainer>
    </Card>
  );
}
