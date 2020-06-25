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

import * as React from 'react';
import Button from '@fbcnms/ui/components/design-system/Button';
import Card from '@fbcnms/ui/components/design-system/Card/Card';
import PermissionsGroupManagePoliciesDialog, {
  DIALOG_TITLE,
} from './PermissionsGroupManagePoliciesDialog';
import PermissionsPoliciesTable from '../policies/PermissionsPoliciesTable';
import PermissionsPolicyRulesDisplay from '../policies/PermissionsPolicyRulesDisplay';
import Text from '@fbcnms/ui/components/design-system/Text';
import ViewContainer from '@fbcnms/ui/components/design-system/View/ViewContainer';
import fbt from 'fbt';
import {NewTabIcon} from '@fbcnms/ui/components/design-system/Icons';
import {POSITION} from '@fbcnms/ui/components/design-system/Dialog/DialogFrame';
import {ROW_SEPARATOR_TYPES} from '@fbcnms/ui/components/design-system/Table/TableContent';
import {TABLE_VARIANT_TYPES} from '@fbcnms/ui/components/design-system/Table/Table';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useMemo, useState} from 'react';
import {useDialogShowingContext} from '@fbcnms/ui/components/design-system/Dialog/DialogShowingContext';
import {wrapRawPermissionsPolicies} from '../data/PermissionsPolicies';

const useStyles = makeStyles(() => ({
  policyCardTitle: {
    display: 'flex',
    flexDirection: 'column',
    paddingLeft: '8px',
  },
  policyName: {
    marginBottom: '8px',
  },
  actionsBar: {
    marginTop: '24px',
  },
}));

type Props = $ReadOnly<{|
  group: UsersGroup,
  className?: ?string,
  onChange: UsersGroup => void,
|}>;

export default function PermissionsGroupPoliciesPane(props: Props) {
  const {group, className, onChange} = props;

  const classes = useStyles();

  const policies = useMemo(() => wrapRawPermissionsPolicies(group.policies), [
    group.policies,
  ]);
  const [showManagePoliciesDialog, setShowManagePoliciesDialog] = useState(
    false,
  );

  const dialogShowingContext = useDialogShowingContext();

  const showPolicyDetails = useCallback(
    policyId => {
      if (policyId == null) {
        return;
      }

      const policy = policies.find(p => p.id === policyId);
      if (policy == null) {
        return;
      }

      const title = (
        <div className={classes.policyCardTitle}>
          <Text className={classes.policyName} variant="h6">
            {policy.name}
          </Text>
          <Text variant="body2" color="gray">
            {policy.description}
          </Text>
          <div className={classes.actionsBar}>
            <Button
              rightIcon={NewTabIcon}
              onClick={() =>
                window.open(`/admin/user_management/policy/${policy.id}`)
              }>
              <fbt desc="">Edit Policy</fbt>
            </Button>
          </div>
        </div>
      );

      dialogShowingContext.showDialog(
        {
          title,
          children: <PermissionsPolicyRulesDisplay policy={policy} />,
          onClose: dialogShowingContext.hideDialog,
          position: POSITION.right,
        },
        true,
      );
    },
    [
      policies,
      classes.policyCardTitle,
      classes.policyName,
      classes.actionsBar,
      dialogShowingContext,
    ],
  );

  return (
    <Card className={className} margins="none">
      <ViewContainer
        header={{
          title: <fbt desc="">Policies</fbt>,
          subtitle: (
            <fbt desc="">
              Add policies to apply them on members in this group.
            </fbt>
          ),
          actionButtons: [
            <Button onClick={() => setShowManagePoliciesDialog(true)}>
              {DIALOG_TITLE}
            </Button>,
          ],
        }}>
        {policies.length > 0 ? (
          <PermissionsPoliciesTable
            policies={policies}
            showGroupsColumn={false}
            variant={TABLE_VARIANT_TYPES.embedded}
            dataRowsSeparator={ROW_SEPARATOR_TYPES.border}
            onPolicySelected={showPolicyDetails}
          />
        ) : null}
      </ViewContainer>
      <PermissionsGroupManagePoliciesDialog
        selectedPolicies={group.policies}
        onClose={policies => {
          if (policies != null) {
            onChange({
              ...group,
              policies,
            });
          }
          setShowManagePoliciesDialog(false);
        }}
        open={showManagePoliciesDialog}
      />
    </Card>
  );
}
