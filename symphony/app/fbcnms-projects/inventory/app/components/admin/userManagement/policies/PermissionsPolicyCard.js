/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {PermissionsPolicy} from '../utils/UserManagementUtils';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';

import * as React from 'react';
import Breadcrumbs from '@fbcnms/ui/components/Breadcrumbs';
import DeleteIcon from '@fbcnms/ui/components/design-system/Icons/Actions/DeleteIcon';
import Grid from '@material-ui/core/Grid';
import InventoryErrorBoundary from '../../../../common/InventoryErrorBoundary';
import PermissionsPolicyDetailsPane from './PermissionsPolicyDetailsPane';
import PermissionsPolicyGroupsPane from './PermissionsPolicyGroupsPane';
import PermissionsPolicyRulesPane from './PermissionsPolicyRulesPane';
import Strings from '@fbcnms/strings/Strings';
import ViewContainer from '@fbcnms/ui/components/design-system/View/ViewContainer';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {
  ButtonAction,
  IconAction,
} from '@fbcnms/ui/components/design-system/View/ViewHeaderActions';
import {
  NEW_DIALOG_PARAM,
  POLICY_TYPES,
  initInventoryRulesInput,
  initWorkforceRulesInput,
} from '../utils/UserManagementUtils';
import {PERMISSION_POLICIES_VIEW_NAME} from './PermissionsPoliciesView';
import {generateTempId} from '../../../../common/EntUtils';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useEffect, useMemo, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useLocation} from 'react-router-dom';
import {useParams} from 'react-router';
import {useUserManagement} from '../UserManagementContext';

const useStyles = makeStyles(() => ({
  detailsPane: {
    display: 'flex',
    flexDirection: 'column',
    borderRadius: '4px',
    boxShadow: symphony.shadows.DP1,
    '&:not(:first-child)': {
      marginTop: '16px',
    },
  },
  container: {
    maxHeight: '100%',
  },
}));

type Props = $ReadOnly<{|
  redirectToPoliciesView: () => void,
  onClose: () => void,
  ...WithAlert,
|}>;

const getInitialNewPolicy: (policyType: ?string) => PermissionsPolicy = (
  policyType: ?string,
) => {
  let type = POLICY_TYPES.InventoryPolicy.key;
  if (policyType == POLICY_TYPES.WorkforcePolicy.key) {
    type = POLICY_TYPES.WorkforcePolicy.key;
  }

  return {
    id: generateTempId(),
    name: '',
    description: '',
    type,
    isGlobal: false,
    groups: [],
    inventoryRules:
      type === POLICY_TYPES.InventoryPolicy.key
        ? initInventoryRulesInput()
        : null,
    workforceRules:
      type === POLICY_TYPES.WorkforcePolicy.key
        ? initWorkforceRulesInput()
        : null,
  };
};

function PermissionsPolicyCard(props: Props) {
  const {redirectToPoliciesView, onClose} = props;
  const classes = useStyles();
  const location = useLocation();
  const {
    policies,
    addPermissionsPolicy,
    editPermissionsPolicy,
    deletePermissionsPolicy,
  } = useUserManagement();
  const {id: policyId} = useParams();
  const isOnNewPolicy = policyId?.startsWith(NEW_DIALOG_PARAM) || false;
  const queryParams = new URLSearchParams(location.search);
  const [policy, setPolicy] = useState<?PermissionsPolicy>(
    isOnNewPolicy ? getInitialNewPolicy(queryParams.get('type')) : null,
  );

  useEffect(() => {
    if (isOnNewPolicy) {
      return;
    }
    const requestedPolicy =
      policyId == null ? null : policies.find(policy => policy.id === policyId);
    if (requestedPolicy == null) {
      redirectToPoliciesView();
    }
    setPolicy(requestedPolicy);
  }, [policyId, isOnNewPolicy, redirectToPoliciesView, policies]);

  const enqueueSnackbar = useEnqueueSnackbar();
  const handleError = useCallback(
    (error: string) => {
      enqueueSnackbar(error, {variant: 'error'});
    },
    [enqueueSnackbar],
  );

  const header = useMemo(() => {
    const breadcrumbs = [
      {
        id: 'policies',
        name: `${PERMISSION_POLICIES_VIEW_NAME}`,
        onClick: redirectToPoliciesView,
      },
      {
        id: 'policyName',
        name: isOnNewPolicy ? `${fbt('New Policy', '')}` : policy?.name || '',
      },
    ];
    const actions = [
      <ButtonAction skin="regular" action={onClose}>
        {Strings.common.cancelButton}
      </ButtonAction>,
      <ButtonAction
        disableOnFromError={true}
        action={() => {
          if (policy == null) {
            return;
          }

          const saveAction = isOnNewPolicy
            ? addPermissionsPolicy
            : editPermissionsPolicy;
          saveAction(policy)
            .then(onClose)
            .catch(handleError);
        }}>
        {Strings.common.saveButton}
      </ButtonAction>,
    ];
    if (!isOnNewPolicy) {
      actions.unshift(
        <IconAction
          icon={DeleteIcon}
          skin="gray"
          action={() => {
            if (policy == null) {
              return;
            }
            props
              .confirm(
                <fbt desc="">Are you sure you want to delete this policy?</fbt>,
              )
              .then(confirm => {
                if (!confirm) {
                  return;
                }
                return deletePermissionsPolicy(policy.id).then(onClose);
              })
              .catch(handleError);
          }}
        />,
      );
    }
    return {
      title: <Breadcrumbs breadcrumbs={breadcrumbs} />,
      subtitle: fbt('Edit this policy and apply it to groups.', ''),
      actionButtons: actions,
    };
  }, [
    redirectToPoliciesView,
    isOnNewPolicy,
    policy,
    onClose,
    addPermissionsPolicy,
    editPermissionsPolicy,
    handleError,
    props,
    deletePermissionsPolicy,
  ]);

  if (policy == null) {
    return null;
  }
  return (
    <InventoryErrorBoundary>
      <ViewContainer header={header} useBodyScrollingEffect={false}>
        <Grid container spacing={2} className={classes.container}>
          <Grid item xs={8} sm={8} lg={8} xl={8} className={classes.container}>
            <PermissionsPolicyDetailsPane
              policy={policy}
              onChange={setPolicy}
              className={classes.detailsPane}
            />
            <PermissionsPolicyRulesPane
              policy={policy}
              onChange={setPolicy}
              className={classes.detailsPane}
            />
          </Grid>
          <Grid item xs={4} sm={4} lg={4} xl={4} className={classes.container}>
            <PermissionsPolicyGroupsPane
              policy={policy}
              className={classes.detailsPane}
            />
          </Grid>
        </Grid>
      </ViewContainer>
    </InventoryErrorBoundary>
  );
}

export default withAlert(PermissionsPolicyCard);
