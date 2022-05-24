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
 *
 * @flow strict-local
 * @format
 */
import type {WithAlert} from '../../components/Alert/withAlert';
import type {
  policy_rule,
  qos_class_id,
} from '../../../generated/MagmaAPIBindings';

import ActionTable from '../../components/ActionTable';
import BaseNameEditDialog from './BaseNameEdit';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Grid from '@material-ui/core/Grid';
import JsonEditor from '../../components/JsonEditor';
import Link from '@material-ui/core/Link';
import LteNetworkContext from '../../components/context/LteNetworkContext';
import PolicyContext from '../../components/context/PolicyContext';
import PolicyRuleEditDialog from './PolicyEdit';
import ProfileEditDialog from './ProfileEdit';
import RatingGroupEditDialog from './RatingGroupEdit';
import React from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Text from '../../theme/design-system/Text';
import TextField from '@material-ui/core/TextField';
import withAlert from '../../components/Alert/withAlert';

import {Checkbox} from '@material-ui/core';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext, useEffect, useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';
import {useNavigate, useParams} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const useStyles = makeStyles(theme => ({
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
    padding: `0 ${theme.spacing(5)}px`,
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
  policyID: string,
  numFlows: number,
  priority: number,
  numSubscribers: number,
  monitoringKey: string,
  rating: string,
  trackingType: string,
  networkWide: string,
};

type BaseNameRowType = {
  baseNameID: string,
  ruleNames: Array<string>,
  numSubscribers: number,
};

type ProfileRowType = {
  classID: qos_class_id,
  profileID: string,
  uplinkBandwidth: number,
  downlinkBandwidth: number,
};

type RatingGroupRowType = {
  ratingGroupID: string,
  limitType: string,
};

const MagmaTabs = withStyles({
  indicator: {
    backgroundColor: colors.secondary.dodgerBlue,
  },
})(Tabs);

const MagmaTab = withStyles({
  root: {
    fontFamily: typography.body1.fontFamily,
    fontWeight: typography.body1.fontWeight,
    fontSize: typography.body1.fontSize,
    lineHeight: typography.body1.lineHeight,
    letterSpacing: typography.body1.letterSpacing,
    color: colors.primary.brightGray,
    textTransform: 'none',
  },
})(Tab);

export default function PolicyOverview() {
  const classes = useStyles();
  const [currTabIndex, setCurrTabIndex] = useState<number>(0);
  const policyTabList: Array<string> = [
    'Policies',
    'Base Names',
    'Profiles',
    'Rating Groups',
  ];

  return (
    <div className={classes.dashboardRoot}>
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
    </div>
  );
}

export function PolicyTableRaw(props: WithAlert) {
  const enqueueSnackbar = useEnqueueSnackbar();
  const [open, setOpen] = React.useState(false);
  const [currRow, setCurrRow] = useState<PolicyRowType>({});
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
          policyID: policyRule.id,
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
        rule={
          Object.keys(currRow).length ? policies[currRow.policyID] : undefined
        }
      />
      <ActionTable
        data={policyRows}
        columns={[
          {
            title: 'Policy ID',
            field: 'policyID',
            render: currRow => (
              <Link
                variant="body2"
                component="button"
                onClick={() => {
                  setCurrRow(currRow);
                  setOpen(true);
                }}>
                {currRow.policyID}
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
              navigate(currRow.policyID + '/json');
            },
          },
          {name: 'Deactivate'},
          {
            name: 'Remove',
            handleFunc: () => {
              props
                .confirm(`Are you sure you want to delete ${currRow.policyID}?`)
                .then(async confirmed => {
                  if (!confirmed) {
                    return;
                  }

                  try {
                    // trigger deletion
                    ctx.setState(currRow.policyID);
                  } catch (e) {
                    enqueueSnackbar(
                      'failed deleting policy ' + currRow.policyID,
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

export function BaseNameTableRaw(props: WithAlert) {
  const enqueueSnackbar = useEnqueueSnackbar();
  const [open, setOpen] = React.useState(false);
  const [currRow, setCurrRow] = useState<BaseNameRowType>({});
  const ctx = useContext(PolicyContext);
  const baseNames = ctx.baseNames;
  const baseNameRows: Array<BaseNameRowType> = baseNames
    ? Object.keys(baseNames).map((baseNameID: string) => {
        const baseNameRecord = baseNames[baseNameID];
        return {
          baseNameID: baseNameID,
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
        baseNameId={currRow?.baseNameID}
      />
      <ActionTable
        data={baseNameRows}
        columns={[
          {
            title: 'Base Name ID',
            field: 'baseNameID',
            render: currRow => (
              <Link
                variant="body2"
                component="button"
                onClick={() => {
                  setCurrRow(currRow);
                  setOpen(true);
                }}>
                {currRow.baseNameID}
              </Link>
            ),
          },
          {
            title: 'Rule Names',
            field: 'ruleNames',
            render: rowData => {
              return (
                <TextField
                  value={rowData.ruleNames.join(', ')}
                  InputProps={{
                    disableUnderline: true,
                    readOnly: true,
                  }}
                />
              );
            },
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
              props
                .confirm(
                  `Are you sure you want to delete ${currRow.baseNameID}?`,
                )
                .then(async confirmed => {
                  if (!confirmed) {
                    return;
                  }

                  try {
                    // trigger deletion
                    ctx.setBaseNames(currRow.baseNameID);
                  } catch (e) {
                    enqueueSnackbar(
                      'failed deleting base name ' + currRow.baseNameID,
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

export function ProfileTableRaw(props: WithAlert) {
  const enqueueSnackbar = useEnqueueSnackbar();
  const [open, setOpen] = React.useState(false);
  const [currRow, setCurrRow] = useState<ProfileRowType>({});
  const ctx = useContext(PolicyContext);
  const profiles = ctx.qosProfiles;
  const profileRows: Array<ProfileRowType> = profiles
    ? Object.keys(profiles).map((profileID: string) => {
        const profile = profiles[profileID];
        return {
          profileID: profile.id,
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
        profile={
          Object.keys(currRow).length ? profiles[currRow.profileID] : undefined
        }
      />
      <ActionTable
        data={profileRows}
        columns={[
          {
            title: 'Profile ID',
            field: 'profileID',
            render: currRow => (
              <Link
                variant="body2"
                component="button"
                onClick={() => {
                  setCurrRow(currRow);
                  setOpen(true);
                }}>
                {currRow.profileID}
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
              props
                .confirm(
                  `Are you sure you want to delete ${currRow.profileID}?`,
                )
                .then(async confirmed => {
                  if (!confirmed) {
                    return;
                  }

                  try {
                    // trigger deletion
                    ctx.setQosProfiles(currRow.profileID);
                  } catch (e) {
                    enqueueSnackbar(
                      'failed deleting profile ' + currRow.profileID,
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

export function RatingGroupTableRaw(props: WithAlert) {
  const enqueueSnackbar = useEnqueueSnackbar();
  const [open, setOpen] = React.useState(false);
  const [currRow, setCurrRow] = useState<RatingGroupRowType>({});
  const ctx = useContext(PolicyContext);
  const ratingGroups = ctx.ratingGroups;
  const ratingGroupRow: Array<RatingGroupRowType> = ratingGroups
    ? Object.keys(ratingGroups).map((ratingGroupID: string) => {
        const ratingGroup = ratingGroups[ratingGroupID];
        return {
          ratingGroupID: ratingGroup.id.toString(),
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
          Object.keys(currRow).length
            ? ratingGroups[currRow.ratingGroupID]
            : undefined
        }
      />
      <ActionTable
        data={ratingGroupRow}
        columns={[
          {
            title: 'Rating Group ID',
            field: 'ratingGroupID',
            render: currRow => (
              <Link
                variant="body2"
                component="button"
                onClick={() => {
                  setCurrRow(currRow);
                  setOpen(true);
                }}>
                {currRow.ratingGroupID}
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
              props
                .confirm(
                  `Are you sure you want to delete Rating Group ${currRow.ratingGroupID}?`,
                )
                .then(async confirmed => {
                  if (!confirmed) {
                    return;
                  }

                  try {
                    ctx.setRatingGroups(currRow.ratingGroupID.toString());
                  } catch (e) {
                    enqueueSnackbar(
                      'failed deleting rating group ' + currRow.ratingGroupID,
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
  const policyID: string = params.policyId;
  const enqueueSnackbar = useEnqueueSnackbar();
  const ctx = useContext(PolicyContext);
  const lteNetworkCtx = useContext(LteNetworkContext);
  const policies = ctx.state;
  const policy: policy_rule = policies[policyID] || DEFAULT_POLICY_CONFIG;
  const lteNetwork = lteNetworkCtx.state;
  const [isNetworkWide, setIsNetworkWide] = useState(false);

  useEffect(() => {
    if (policyID) {
      setIsNetworkWide(
        lteNetwork?.subscriber_config?.network_wide_rule_names?.includes(
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
          setError(e.response?.data?.message ?? e.message);
        }
      }}
    />
  );
}

const PolicyTable = withAlert(PolicyTableRaw);
const BaseNameTable = withAlert(BaseNameTableRaw);
const ProfileTable = withAlert(ProfileTableRaw);
const RatingGroupTable = withAlert(RatingGroupTableRaw);
