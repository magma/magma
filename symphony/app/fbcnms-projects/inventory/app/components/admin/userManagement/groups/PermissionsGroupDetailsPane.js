/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {UsersGroup} from '../data/UsersGroups';
import type {UsersGroupStatus} from '../__generated__/UserManagementContextQuery.graphql';

import * as React from 'react';
import AppContext from '@fbcnms/ui/context/AppContext';
import Card from '@fbcnms/ui/components/design-system/Card/Card';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import FormFieldTextInput from '../utils/FormFieldTextInput';
import Grid from '@material-ui/core/Grid';
import Select from '@fbcnms/ui/components/design-system/Select/Select';
import ViewContainer from '@fbcnms/ui/components/design-system/View/ViewContainer';
import fbt from 'fbt';
import {GROUP_STATUSES} from '../utils/UserManagementUtils';
import {isTempId} from '../../../../common/EntUtils';
import {makeStyles} from '@material-ui/styles';
import {useContext, useMemo} from 'react';

const useStyles = makeStyles(() => ({
  viewContainer: {
    paddingBottom: '16px',
  },
  nameField: {
    marginRight: '8px',
  },
  descriptionField: {
    marginTop: '8px',
  },
}));

type Props = $ReadOnly<{
  group: UsersGroup,
  onChange: UsersGroup => void,
  className?: ?string,
}>;

export default function PermissionsGroupDetailsPane(props: Props) {
  const {group, className, onChange} = props;
  const {isFeatureEnabled} = useContext(AppContext);
  const permissionPoliciesMode = isFeatureEnabled('permission_policies');
  const classes = useStyles();
  const isNewGroup = isTempId(group.id);

  const statuses = useMemo(
    () =>
      Object.keys(GROUP_STATUSES).map((statusKey: UsersGroupStatus) => ({
        key: statusKey,
        value: statusKey,
        label: GROUP_STATUSES[statusKey].value,
      })),
    [],
  );

  return (
    <Card className={className} margins="none">
      <ViewContainer
        className={classes.viewContainer}
        header={{title: <fbt desc="">Group Details</fbt>}}>
        <Grid container>
          <Grid item xs={12} sm={6} lg={6} xl={6}>
            <FormFieldTextInput
              className={classes.nameField}
              disabled={!permissionPoliciesMode}
              label={`${fbt('Group Name', '')}`}
              validationId="name"
              value={group.name}
              onValueChanged={name => {
                onChange({
                  ...group,
                  name,
                });
              }}
            />
          </Grid>
          {permissionPoliciesMode && !isNewGroup ? (
            <Grid item xs={12} sm={6} lg={6} xl={6}>
              <FormField label={`${fbt('Status', '')}`}>
                <Select
                  options={statuses}
                  selectedValue={group.status}
                  onChange={status =>
                    onChange({
                      ...group,
                      status,
                    })
                  }
                />
              </FormField>
            </Grid>
          ) : null}
          <Grid item xs={12}>
            <FormFieldTextInput
              className={classes.descriptionField}
              label={`${fbt('Group Description', '')}`}
              value={group.description || ''}
              onValueChanged={description => {
                onChange({
                  ...group,
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
