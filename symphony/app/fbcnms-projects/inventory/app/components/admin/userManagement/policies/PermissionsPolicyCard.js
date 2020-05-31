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
import Button from '@fbcnms/ui/components/design-system/Button';
import DeleteIcon from '@fbcnms/ui/components/design-system/Icons/Actions/DeleteIcon';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import Grid from '@material-ui/core/Grid';
import IconButton from '@fbcnms/ui/components/design-system/IconButton';
import InventoryErrorBoundary from '../../../../common/InventoryErrorBoundary';
import PermissionsPolicyDetailsPane from './PermissionsPolicyDetailsPane';
import PermissionsPolicyGroupsPane from './PermissionsPolicyGroupsPane';
import PermissionsPolicyRulesPane from './PermissionsPolicyRulesPane';
import Strings from '@fbcnms/strings/Strings';
import ViewContainer from '@fbcnms/ui/components/design-system/View/ViewContainer';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {FormContextProvider} from '../../../../common/FormContext';
import {
  NEW_DIALOG_PARAM,
  PERMISSION_RULE_VALUES,
  POLICY_TYPES,
} from '../utils/UserManagementUtils';
import {PERMISSION_POLICIES_VIEW_NAME} from './PermissionsPoliciesView';
import {generateTempId} from '../../../../common/EntUtils';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useEffect, useMemo, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useFormAlertsContext} from '@fbcnms/ui/components/design-system/Form/FormAlertsContext';
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

const initialBasicRule = {
  isAllowed: PERMISSION_RULE_VALUES.NO,
};

const initialCUDRule = {
  create: {
    ...initialBasicRule,
  },
  update: {
    ...initialBasicRule,
  },
  delete: {
    ...initialBasicRule,
  },
};

const initialInventoryRules = {
  read: {
    isAllowed: PERMISSION_RULE_VALUES.YES,
  },
  location: {
    ...initialCUDRule,
  },
  equipment: {
    ...initialCUDRule,
  },
  equipmentType: {
    ...initialCUDRule,
  },
  locationType: {
    ...initialCUDRule,
  },
  portType: {
    ...initialCUDRule,
  },
  serviceType: {
    ...initialCUDRule,
  },
};

const initialWorkforceCUDRules = {
  ...initialCUDRule,
  assign: {
    ...initialBasicRule,
  },
  transferOwnership: {
    ...initialBasicRule,
  },
};

const initialWorkforceRules = {
  read: {
    ...initialBasicRule,
  },
  data: {
    ...initialWorkforceCUDRules,
  },
  templates: {
    ...initialCUDRule,
  },
};

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
      type === POLICY_TYPES.InventoryPolicy.key ? initialInventoryRules : null,
    workforceRules:
      type === POLICY_TYPES.WorkforcePolicy.key ? initialWorkforceRules : null,
  };
};

function PermissionsPolicyCard(props: Props) {
  const {redirectToPoliciesView, onClose} = props;
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
    const actions =
      policy?.isSystemDefault === true
        ? [
            <FormAction ignorePermissions={true} ignoreEditLocks={true}>
              <Button onClick={onClose}>{Strings.common.doneButton}</Button>
            </FormAction>,
          ]
        : [
            <FormAction ignorePermissions={true}>
              <Button skin="regular" onClick={onClose}>
                {Strings.common.cancelButton}
              </Button>
            </FormAction>,
            <FormAction disableOnFromError={true}>
              <Button
                onClick={() => {
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
              </Button>
            </FormAction>,
          ];
    if (!isOnNewPolicy && policy?.isSystemDefault !== true) {
      actions.unshift(
        <FormAction>
          <IconButton
            icon={DeleteIcon}
            skin="gray"
            onClick={() => {
              if (policy == null) {
                return;
              }
              props
                .confirm(
                  <fbt desc="">
                    Are you sure you want to delete this policy?
                  </fbt>,
                )
                .then(confirm => {
                  if (!confirm) {
                    return;
                  }
                  return deletePermissionsPolicy(policy.id).then(onClose);
                })
                .catch(handleError);
            }}
          />
        </FormAction>,
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
      <FormContextProvider permissions={{adminRightsRequired: true}}>
        <ViewContainer header={header} useBodyScrollingEffect={false}>
          <PermissionsPolicyCardBody policy={policy} onChange={setPolicy} />
        </ViewContainer>
      </FormContextProvider>
    </InventoryErrorBoundary>
  );
}

type PermissionsPolicyCardBodyProps = $ReadOnly<{|
  policy: PermissionsPolicy,
  onChange: PermissionsPolicy => void,
|}>;

function PermissionsPolicyCardBody(props: PermissionsPolicyCardBodyProps) {
  const {policy, onChange} = props;
  const classes = useStyles();

  const alerts = useFormAlertsContext();
  alerts.editLock.check({
    fieldId: 'system_default_policy',
    fieldDisplayName: 'Workforce Default Policy',
    value: policy.isSystemDefault,
    checkCallback: isSystemDefault =>
      isSystemDefault
        ? `${fbt(
            'This policy is applied to all users by default. It cannot be edited or removed.',
            '',
          )}`
        : '',
  });

  const policyDetailsPart = (
    <PermissionsPolicyDetailsPane
      policy={policy}
      onChange={onChange}
      className={classes.detailsPane}
    />
  );

  if (policy.isSystemDefault) {
    return (
      <Grid container spacing={2} className={classes.container}>
        <Grid item xs={12} className={classes.container}>
          {policyDetailsPart}
        </Grid>
      </Grid>
    );
  }

  return (
    <Grid container spacing={2} className={classes.container}>
      <Grid item xs={8} sm={8} lg={8} xl={8} className={classes.container}>
        {policyDetailsPart}
        <PermissionsPolicyRulesPane
          policy={policy}
          onChange={onChange}
          className={classes.detailsPane}
        />
      </Grid>
      <Grid item xs={4} sm={4} lg={4} xl={4} className={classes.container}>
        <PermissionsPolicyGroupsPane
          policy={policy}
          onChange={onChange}
          className={classes.detailsPane}
        />
      </Grid>
    </Grid>
  );
}

export default withAlert(PermissionsPolicyCard);
