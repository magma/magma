/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {InventoryPolicy} from '../utils/UserManagementUtils';

import * as React from 'react';
import AppContext from '@fbcnms/ui/context/AppContext';
import PermissionsPolicyRulesSection from './PermissionsPolicyRulesSection';

import Switch from '@fbcnms/ui/components/design-system/switch/Switch';
import fbt from 'fbt';
import {
  bool2PermissionRuleValue,
  permissionRuleValue2Bool,
} from '../utils/UserManagementUtils';
import {makeStyles} from '@material-ui/styles';
import {useContext} from 'react';

const useStyles = makeStyles(() => ({
  root: {
    marginLeft: '4px',
    maxHeight: '100%',
    display: 'flex',
    flexDirection: 'column',
  },
  readRule: {
    marginLeft: '4px',
  },
  section: {
    marginTop: '32px',
  },
}));

type Props = $ReadOnly<{|
  policy: ?InventoryPolicy,
  onChange: InventoryPolicy => void,
|}>;

export default function PermissionsPolicyInventoryDataRulesTab(props: Props) {
  const {policy, onChange} = props;
  const classes = useStyles();
  const {isFeatureEnabled} = useContext(AppContext);
  const userManagementDevMode = isFeatureEnabled('user_management_dev');

  if (policy == null) {
    return null;
  }

  const readAllowed = permissionRuleValue2Bool(policy.read.isAllowed);

  return (
    <div className={classes.root}>
      {userManagementDevMode ? (
        <Switch
          className={classes.readRule}
          title={fbt('View inventory data', '')}
          checked={readAllowed}
          onChange={checked =>
            onChange({
              ...policy,
              read: {
                isAllowed: bool2PermissionRuleValue(checked),
              },
            })
          }
        />
      ) : null}
      <PermissionsPolicyRulesSection
        title={fbt('Locations', '')}
        subtitle={fbt(
          'Location data includes location details, properties, floor plans and coverage maps.',
          '',
        )}
        disabled={!readAllowed}
        rule={policy.location}
        className={classes.section}
        onChange={location =>
          onChange({
            ...policy,
            location,
          })
        }
      />
      <PermissionsPolicyRulesSection
        title={fbt('Equipment', '')}
        subtitle={fbt(
          'Equipment data includes equipment items, ports, links, services and network maps.',
          '',
        )}
        className={classes.section}
        disabled={!readAllowed}
        rule={policy.equipment}
        onChange={equipment =>
          onChange({
            ...policy,
            equipment,
          })
        }
      />
    </div>
  );
}
