/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
import withAlert from '../../components/Alert/withAlert';
import type {PolicyRule, QosClassId} from '../../../generated';
import type {WithAlert} from '../../components/Alert/withAlert';

import ActionTable from '../../components/ActionTable';
import BaseNameEditDialog from './BaseNameEdit';
import EmptyState from '../../components/EmptyState';
import FormControlLabel from '@mui/material/FormControlLabel';
import Grid from '@mui/material/Grid';
import JsonEditor from '../../components/JsonEditor';
import Link from '@mui/material/Link';
import LteNetworkContext from '../../context/LteNetworkContext';
import PolicyContext from '../../context/PolicyContext';
import PolicyRuleEditDialog from './PolicyEdit';
import ProfileEditDialog from './ProfileEdit';
import RatingGroupEditDialog from './RatingGroupEdit';
import React, {useContext, useEffect, useState} from 'react';
import Text from '../../theme/design-system/Text';
import TextField from '@mui/material/TextField';
import {Checkbox} from '@mui/material';
import {MagmaTab, MagmaTabs} from '../../theme/design-system/SecondaryTabs';
import {Theme} from '@mui/material/styles';
import {colors, typography} from '../../theme/default';
import {getErrorMessage} from '../../util/ErrorUtils';
import {makeStyles} from '@mui/styles';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import {useNavigate, useParams} from 'react-router-dom';

const EMPTY_STATE_OVERVIEW =
  'A policy controls the behavior of a packet flow and determines if the flow needs to be allowed/blocked or have ' +
  'its QoS restricted. Policies may also redirect traffic, perform header enrichment, include tracking/monitoring and other advanced configurations.';
const useStyles = makeStyles<Theme>(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
  topBar: {
    backgroundColor: colors.primary.mirage,
    padding: '20px 40px 20px 40px',
    color: colors.primary.white,
  },
  tabBar: {
    backgroundColor: colors.primary.brightGray,
    padding: `0 ${theme.spacing(5)}`,
  },
  tabs: {
    color: colors.primary.white,
  },
  tab: {
    backgroundColor: colors.primary.white,
    borderRadius: '4px 4px 0 0',
    boxShadow: `inset 0 -2px 0 0 ${colors.primary.concrete}`,
    '& + &': {
      marginLeft: '4px',
    },
  },
  tabLabel: {
    padding: '16px 0 16px 0',
    display: 'flex',
    alignItems: 'center',
  },
  tabIconLabel: {
    marginRight: '8px',
  },
  appBarBtn: {
    color: colors.primary.white,
    background: colors.primary.comet,
    fontFamily: typography.button.fontFamily,
    fontWeight: typography.button.fontWeight,
    fontSize: typography.button.fontSize,
    lineHeight: typography.button.lineHeight,
    letterSpacing: typography.button.letterSpacing,

    '&:hover': {
      background: colors.primary.mirage,
    },
  },
  appBarBtnSecondary: {
    textPrimary: colors.primary.mirage,
    outlined: true,
    contained: true,
    color: colors.primary.mirage,
  },
  // TODO: remove this when we actually fill out the grid sections
  contentPlaceholder: {
    padding: '50px 0',
  },
  paper: {
    height: 100,
    padding: theme.spacing(10),
    textAlign: 'center',
  },
  formControl: {
    margin: theme.spacing(1),
    minWidth: 120,
  },
}));

type PolicyRowType = {
  id: string;
  numFlows: number;
  priority: number;
  numSubscribers: number;
  monitoringKey: string;
  rating: string;
  trackingType: string;
  networkWide: string;
};

type BaseNameRowType = {
  id: string;
  ruleNames: Array<string>;
  numSubscribers: number;
};

type ProfileRowType = {
  classID: QosClassId;
  id: string;
  uplinkBandwidth: number;
  downlinkBandwidth: number;
};

type RatingGroupRowType = {
  id: string;
  limitType: string;
};

export default function PolicyOverview() {
  const classes = useStyles();
  const [currTabIndex, setCurrTabIndex] = useState<number>(0);
  const [open, setOpen] = useState(false);
  const policyTabList: Array<string> = [
    'Policies',
    'Base Names',
    'Profiles',
    'Rating Groups',
  ];
  const ctx = useContext(PolicyContext);
  const cardActions = {
    buttonText: 'Add Policy',
    onClick: () => setOpen(true),
    linkText: 'Learn more about Policy',
    link: 'https://docs.magmacore.org/docs/nms/traffic#policy-configuration',
  };

  const isEmpty =
    Object.keys(ctx.state || {}).length === 0 &&
    Object.keys(ctx.ratingGroups || {}).length === 0 &&
    Object.keys(ctx.qosProfiles || {}).length === 0 &&
    Object.keys(ctx.baseNames || {}).length === 0;

  return (
    <div className={classes.dashboardRoot}>
      <PolicyRuleEditDialog
        open={open}
        onClose={() => setOpen(false)}
        rule={undefined}
      />
      {isEmpty ? (
        <Grid container justifyContent="space-between" spacing={3}>
          <EmptyState
            title={'Set up a Policy'}
            instructions={
              'Add a policy to the NMS by filling out the required fields.'
            }
            cardActions={cardActions}
            overviewTitle={'Policy Overview'}
            overviewDescription={EMPTY_STATE_OVERVIEW}
          />
        </Grid>
      ) : (
        <>
          <MagmaTabs
            value={currTabIndex}
            onChange={(_, newIndex: number) => setCurrTabIndex(newIndex)}
            variant="fullWidth">
            {policyTabList.map((k: string, idx: number) => {
              return <MagmaTab key={idx} label={k} className={classes.tab} />;
            })}
          </MagmaTabs>
          {currTabIndex === 0 && <PolicyTable />}
          {currTabIndex === 1 && <BaseNameTable />}
          {currTabIndex === 2 && <ProfileTable />}
          {currTabIndex === 3 && <RatingGroupTable />}
        </>
      )}
    </div>
  );
}

export function PolicyTableRaw(props: WithAlert) {
  const enqueueSnackbar = useEnqueueSnackbar();
  const [open, setOpen] = React.useState(false);
  const [currRow, setCurrRow] = useState<PolicyRowType>({} as PolicyRowType);
  const navigate = useNavigate();
  const ctx = useContext(PolicyContext);
  const lteNetworkCtx = useContext(LteNetworkContext);
  const lteNetwork = lteNetworkCtx.state;
  const ruleNames = new Set(
    lteNetwork?.subscriber_config?.network_wide_rule_names ?? [],
  );
  const policies = ctx.state;
  const policyRows: Array<PolicyRowType> = policies
    ? Object.keys(policies).map((policyID: string) => {
        const policyRule = policies[policyID];
        return {
          id: policyRule.id,
          numFlows: policyRule.flow_list?.length ?? 0,
          priority: policyRule.priority,
          numSubscribers: policyRule.assigned_subscribers?.length ?? 0,
          monitoringKey: policyRule.monitoring_key ?? '',
          rating: policyRule.rating_group?.toString() ?? 'Not Found',
          trackingType: policyRule.tracking_type ?? 'NO_TRACKING',
          networkWide: ruleNames.has(policyID) ? 'Enabled' : 'Disabled',
        };
      })
    : [];

  return (
    <>
      <PolicyRuleEditDialog
        open={open}
        onClose={() => setOpen(false)}
        rule={Object.keys(currRow).length ? policies[currRow.id] : undefined}
      />
      <ActionTable
        data={policyRows}
        columns={[
          {
            title: 'Policy ID',
            field: 'id',
            render: currRow => (
              <Link
                variant="body2"
                component="button"
                onClick={() => {
                  setCurrRow(currRow);
                  setOpen(true);
                }}
                underline="hover">
                {currRow.id}
              </Link>
            ),
          },
          {title: 'Flows', field: 'numFlows', type: 'numeric'},
          {title: 'Priority', field: 'priority', type: 'numeric'},
          {title: 'Subscribers', field: 'numSubscribers', type: 'numeric'},
          {
            title: 'Monitoring Key',
            field: 'monitoringKey',
            render: rowData => {
              return (
                <TextField
                  type="password"
                  variant="standard"
                  value={rowData.monitoringKey}
                  InputProps={{
                    disableUnderline: true,
                    readOnly: true,
                  }}
                />
              );
            },
          },
          {title: 'Rating', field: 'rating'},
          {title: 'Tracking Type', field: 'trackingType'},
          {title: 'Network Wide', field: 'networkWide'},
        ]}
        handleCurrRow={(row: PolicyRowType) => setCurrRow(row)}
        menuItems={[
          {
            name: 'Edit',
            handleFunc: () => {
              setOpen(true);
            },
          },
          {
            name: 'Edit JSON',
            handleFunc: () => {
              navigate(currRow.id + '/json');
            },
          },
          {name: 'Deactivate'},
          {
            name: 'Remove',
            handleFunc: () => {
              void props
                .confirm(`Are you sure you want to delete ${currRow.id}?`)
                .then(async confirmed => {
                  if (!confirmed) {
                    return;
                  }

                  try {
                    // trigger deletion
                    await ctx.setState(currRow.id);
                  } catch (e) {
                    enqueueSnackbar('failed deleting policy ' + currRow.id, {
                      variant: 'error',
                    });
                  }
                });
            },
          },
        ]}
        options={{
          actionsColumnIndex: -1,
          pageSizeOptions: [5, 10],
        }}
      />
    </>
  );
}

export function BaseNameTableRaw(props: WithAlert) {
  const enqueueSnackbar = useEnqueueSnackbar();
  const [open, setOpen] = React.useState(false);
  const [currRow, setCurrRow] = useState<BaseNameRowType>(
    {} as BaseNameRowType,
  );
  const ctx = useContext(PolicyContext);
  const baseNames = ctx.baseNames;
  const baseNameRows: Array<BaseNameRowType> = baseNames
    ? Object.keys(baseNames).map((baseNameID: string) => {
        const baseNameRecord = baseNames[baseNameID];
        return {
          id: baseNameID,
          ruleNames: baseNameRecord.rule_names,
          numSubscribers: baseNameRecord?.assigned_subscribers?.length || 0,
        };
      })
    : [];

  return (
    <>
      <BaseNameEditDialog
        open={open}
        onClose={() => setOpen(false)}
        baseNameId={currRow?.id}
      />
      <ActionTable
        data={baseNameRows}
        columns={[
          {
            title: 'Base Name ID',
            field: 'id',
            render: currRow => (
              <Link
                variant="body2"
                component="button"
                onClick={() => {
                  setCurrRow(currRow);
                  setOpen(true);
                }}
                underline="hover">
                {currRow.id}
              </Link>
            ),
          },
          {
            title: 'Rule Names',
            field: 'ruleNames',
            render: rowData =>
              rowData.ruleNames ? rowData.ruleNames.join(', ') : '',
          },
          {
            title: '# of Assigned Subscribers',
            field: 'numSubscribers',
            type: 'numeric',
          },
        ]}
        handleCurrRow={(row: BaseNameRowType) => setCurrRow(row)}
        menuItems={[
          {
            name: 'Edit',
            handleFunc: () => {
              setOpen(true);
            },
          },
          {
            name: 'Remove',
            handleFunc: () => {
              void props
                .confirm(`Are you sure you want to delete ${currRow.id}?`)
                .then(async confirmed => {
                  if (!confirmed) {
                    return;
                  }

                  try {
                    // trigger deletion
                    await ctx.setBaseNames(currRow.id);
                  } catch (e) {
                    enqueueSnackbar('failed deleting base name ' + currRow.id, {
                      variant: 'error',
                    });
                  }
                });
            },
          },
        ]}
        options={{
          actionsColumnIndex: -1,
          pageSizeOptions: [5, 10],
        }}
      />
    </>
  );
}

export function ProfileTableRaw(props: WithAlert) {
  const enqueueSnackbar = useEnqueueSnackbar();
  const [open, setOpen] = React.useState(false);
  const [currRow, setCurrRow] = useState<ProfileRowType>({} as ProfileRowType);
  const ctx = useContext(PolicyContext);
  const profiles = ctx.qosProfiles;
  const profileRows: Array<ProfileRowType> = profiles
    ? Object.keys(profiles).map((profileID: string) => {
        const profile = profiles[profileID];
        return {
          id: profile.id,
          classID: profile.class_id,
          uplinkBandwidth: profile.max_req_bw_ul,
          downlinkBandwidth: profile.max_req_bw_dl,
        };
      })
    : [];

  return (
    <>
      <ProfileEditDialog
        open={open}
        onClose={() => setOpen(false)}
        profile={Object.keys(currRow).length ? profiles[currRow.id] : undefined}
      />
      <ActionTable
        data={profileRows}
        columns={[
          {
            title: 'Profile ID',
            field: 'id',
            render: currRow => (
              <Link
                variant="body2"
                component="button"
                onClick={() => {
                  setCurrRow(currRow);
                  setOpen(true);
                }}
                underline="hover">
                {currRow.id}
              </Link>
            ),
          },
          {title: 'Class ID', field: 'classID', type: 'numeric'},
          {
            title: 'Uplink Bandwidth',
            field: 'uplinkBandwidth',
            type: 'numeric',
          },
          {
            title: 'Downlink Bandwidth',
            field: 'downlinkBandwidth',
            type: 'numeric',
          },
        ]}
        handleCurrRow={(row: ProfileRowType) => setCurrRow(row)}
        menuItems={[
          {
            name: 'Edit',
            handleFunc: () => {
              setOpen(true);
            },
          },
          {
            name: 'Remove',
            handleFunc: () => {
              void props
                .confirm(`Are you sure you want to delete ${currRow.id}?`)
                .then(async confirmed => {
                  if (!confirmed) {
                    return;
                  }

                  try {
                    // trigger deletion
                    await ctx.setQosProfiles(currRow.id);
                  } catch (e) {
                    enqueueSnackbar('failed deleting profile ' + currRow.id, {
                      variant: 'error',
                    });
                  }
                });
            },
          },
        ]}
        options={{
          actionsColumnIndex: -1,
          pageSizeOptions: [5, 10],
        }}
      />
    </>
  );
}

export function RatingGroupTableRaw(props: WithAlert) {
  const enqueueSnackbar = useEnqueueSnackbar();
  const [open, setOpen] = React.useState(false);
  const [currRow, setCurrRow] = useState<RatingGroupRowType>(
    {} as RatingGroupRowType,
  );
  const ctx = useContext(PolicyContext);
  const ratingGroups = ctx.ratingGroups;
  const ratingGroupRow: Array<RatingGroupRowType> = ratingGroups
    ? Object.keys(ratingGroups).map((ratingGroupID: string) => {
        const ratingGroup = ratingGroups[ratingGroupID];
        return {
          id: ratingGroup.id.toString(),
          limitType: ratingGroup.limit_type,
        };
      })
    : [];
  return (
    <>
      <RatingGroupEditDialog
        open={open}
        onClose={() => setOpen(false)}
        ratingGroup={
          Object.keys(currRow).length ? ratingGroups[currRow.id] : undefined
        }
      />
      <ActionTable
        data={ratingGroupRow}
        columns={[
          {
            title: 'Rating Group ID',
            field: 'id',
            render: currRow => (
              <Link
                variant="body2"
                component="button"
                onClick={() => {
                  setCurrRow(currRow);
                  setOpen(true);
                }}
                underline="hover">
                {currRow.id}
              </Link>
            ),
          },
          {title: 'Limit type', field: 'limitType'},
        ]}
        handleCurrRow={(row: RatingGroupRowType) => setCurrRow(row)}
        menuItems={[
          {
            name: 'Edit',
            handleFunc: () => {
              setOpen(true);
            },
          },
          {
            name: 'Remove',
            handleFunc: () => {
              void props
                .confirm(
                  `Are you sure you want to delete Rating Group ${currRow.id}?`,
                )
                .then(async confirmed => {
                  if (!confirmed) {
                    return;
                  }

                  try {
                    await ctx.setRatingGroups(currRow.id.toString());
                  } catch (e) {
                    enqueueSnackbar(
                      'failed deleting rating group ' + currRow.id,
                      {
                        variant: 'error',
                      },
                    );
                  }
                });
            },
          },
        ]}
        options={{
          actionsColumnIndex: -1,
          pageSizeOptions: [5, 10],
        }}
      />
    </>
  );
}

// trigger deletion

const DEFAULT_POLICY_CONFIG = {
  flow_list: [],
  id: '',
  monitoring_key: '',
  priority: 1,
};

export function PolicyJsonConfig() {
  const navigate = useNavigate();
  const params = useParams();
  const [error, setError] = useState('');
  const policyID: string = params.policyId!;
  const enqueueSnackbar = useEnqueueSnackbar();
  const ctx = useContext(PolicyContext);
  const lteNetworkCtx = useContext(LteNetworkContext);
  const policies = ctx.state;
  const policy: PolicyRule = policies[policyID] || DEFAULT_POLICY_CONFIG;
  const lteNetwork = lteNetworkCtx.state;
  const [isNetworkWide, setIsNetworkWide] = useState(false);

  useEffect(() => {
    if (policyID) {
      setIsNetworkWide(
        !!lteNetwork?.subscriber_config?.network_wide_rule_names?.includes(
          policyID,
        ),
      );
    }
  }, [policyID, lteNetwork]);
  return (
    <JsonEditor
      content={policy}
      error={error}
      customFilter={
        <Grid item>
          <FormControlLabel
            control={
              <Checkbox
                checked={isNetworkWide}
                onChange={() => setIsNetworkWide(!isNetworkWide)}
                color="primary"
              />
            }
            label={
              <Text weight="medium" variant="body2">
                Network Wide
              </Text>
            }
          />
        </Grid>
      }
      onSave={async policy => {
        try {
          await ctx.setState(policy.id, policy, isNetworkWide);
          enqueueSnackbar('Policy saved successfully', {
            variant: 'success',
          });
          setError('');
          navigate(-1);
        } catch (e) {
          setError(getErrorMessage(e));
        }
      }}
    />
  );
}

const PolicyTable = withAlert(PolicyTableRaw);
const BaseNameTable = withAlert(BaseNameTableRaw);
const ProfileTable = withAlert(ProfileTableRaw);
const RatingGroupTable = withAlert(RatingGroupTableRaw);
