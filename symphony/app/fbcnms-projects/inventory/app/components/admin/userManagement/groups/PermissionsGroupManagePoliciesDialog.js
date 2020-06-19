/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {PermissionsPolicyBase} from '../data/PermissionsPolicies';

import * as React from 'react';
import Button from '@fbcnms/ui/components/design-system/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import PermissionsPoliciesTable from '../policies/PermissionsPoliciesTable';
import Strings from '@fbcnms/strings/Strings';
import Text from '@fbcnms/ui/components/design-system/Text';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {ROW_SEPARATOR_TYPES} from '@fbcnms/ui/components/design-system/Table/TableContent';
import {TABLE_VARIANT_TYPES} from '@fbcnms/ui/components/design-system/Table/Table';
import {makeStyles} from '@material-ui/styles';
import {
  unwrapPermissionsPolicies,
  usePermissionsPolicies,
} from '../data/PermissionsPolicies';
import {useCallback, useEffect, useMemo, useState} from 'react';

const useStyles = makeStyles(() => ({
  title: {
    padding: '26px 32px 8px 32px',
  },
  header: {
    paddingBottom: '2px',
  },
  globalPolicy: {
    // opacity: 0.4,
    fill: symphony.palette.disabled,
    color: symphony.palette.disabled,
    cursor: 'not-allowed',
  },
}));

export const DIALOG_TITLE = fbt('Manage Policies', '');
const NO_CHANGES = fbt('No Change in selected policies', '');

type Props = $ReadOnly<{|
  selectedPolicies: $ReadOnlyArray<PermissionsPolicyBase>,
  open?: ?boolean,
  onClose?: ?(?$ReadOnlyArray<PermissionsPolicyBase>) => void,
|}>;

export default function PermissionsGroupManagePoliciesDialog(props: Props) {
  const {selectedPolicies, open, onClose} = props;
  const classes = useStyles();

  const policies = usePermissionsPolicies();
  const policiesData = useMemo(
    () =>
      policies.map(p => ({
        ...p,
        disabled: p.isGlobal,
        tooltip: p.isGlobal
          ? `${fbt(
              'This policy cannot be selected as it is automatically applied to all users',
              '',
            )}`
          : undefined,
        className: p.isGlobal ? classes.globalPolicy : undefined,
      })),
    [classes.globalPolicy, policies],
  );
  const [selectedIDs, setSelectedIDs] = useState([]);
  useEffect(() => {
    setSelectedIDs(selectedPolicies.map(policy => policy.id));
  }, [selectedPolicies]);

  const hasChanges = useMemo(
    () =>
      selectedIDs.length !== selectedPolicies.length ||
      selectedPolicies.findIndex(
        groupPolicy => selectedIDs.findIndex(id => id === groupPolicy.id) == -1,
      ) > -1,
    [selectedIDs, selectedPolicies],
  );

  const callClose = useCallback(
    (updateChanges: boolean = false) => {
      if (onClose == null) {
        return;
      }
      onClose(
        updateChanges
          ? unwrapPermissionsPolicies(
              policies.filter(p => selectedIDs.includes(p.id)),
            )
          : undefined,
      );
    },
    [onClose, policies, selectedIDs],
  );

  return (
    <Dialog fullWidth={true} maxWidth="md" open={open ?? true}>
      <DialogTitle disableTypography={true} className={classes.title}>
        <Text variant="h6" className={classes.header}>
          {DIALOG_TITLE}
        </Text>
        <Text variant="caption" color="gray" useEllipsis={true}>
          <fbt desc="">
            Add policies to apply them on members in this group.
          </fbt>
        </Text>
      </DialogTitle>
      <DialogContent>
        <PermissionsPoliciesTable
          policies={policiesData}
          showSelection={true}
          variant={TABLE_VARIANT_TYPES.embedded}
          dataRowsSeparator={ROW_SEPARATOR_TYPES.border}
          selectedIds={selectedIDs}
          onSelectionChanged={selectedIds => setSelectedIDs(selectedIds)}
        />
      </DialogContent>
      <DialogActions>
        <FormAction>
          <Button onClick={() => callClose()} skin="regular">
            {Strings.common.cancelButton}
          </Button>
        </FormAction>
        <FormAction
          disableOnFromError={true}
          disabled={!hasChanges}
          tooltip={hasChanges ? undefined : NO_CHANGES}>
          <Button onClick={() => callClose(true)}>
            {Strings.common.updateButton}
          </Button>
        </FormAction>
      </DialogActions>
    </Dialog>
  );
}
