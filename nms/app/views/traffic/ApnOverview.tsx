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
import ActionTable from '../../components/ActionTable';
import ApnContext from '../../components/context/ApnContext';
import ApnEditDialog from './ApnEdit';
import CardTitleRow from '../../components/layout/CardTitleRow';
import EmptyState from '../../components/EmptyState';
import Grid from '@material-ui/core/Grid';
import JsonEditor from '../../components/JsonEditor';
import Link from '@material-ui/core/Link';
import React from 'react';
import RssFeedIcon from '@material-ui/icons/RssFeed';
import withAlert from '../../components/Alert/withAlert';
import {Apn} from '../../../generated';
import {Theme} from '@material-ui/core/styles';
import {colors, typography} from '../../theme/default';
import {getErrorMessage} from '../../util/ErrorUtils';
import {makeStyles} from '@material-ui/styles';
import {useContext, useState} from 'react';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import {useNavigate, useParams} from 'react-router-dom';
import type {WithAlert} from '../../components/Alert/withAlert';

const DEFAULT_APN_CONFIG = {
  apn_configuration: {
    ambr: {
      max_bandwidth_dl: 1000000,
      max_bandwidth_ul: 1000000,
    },
    qos_profile: {
      class_id: 9,
      preemption_capability: false,
      preemption_vulnerability: false,
      priority_level: 15,
    },
  },
  apn_name: '',
};
const EMPTY_STATE_OVERVIEW =
  'APN is an access point name. APN is used to identify the packet data network(PDN), the UE wants to be connected to.' +
  ' From Magma’s perspective, APN configuration consists of two main entities: The APN id and the QoS profile being applied to it.';

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
    padding: `0 ${theme.spacing(5)}px`,
  },
  tabs: {
    color: colors.primary.white,
  },
  tab: {
    fontSize: '18px',
    textTransform: 'none',
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
    color: colors.primary.white,
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

type ApnRowType = {
  apnID: string;
  classID: number;
  arpPriorityLevel: number;
  maxReqdULBw: number;
  maxReqDLBw: number;
  arpPreEmptionCapability: boolean;
  arpPreEmptionVulnerability: boolean;
};

const APN_TITLE = 'APNs';
function ApnOverview(props: WithAlert) {
  const classes = useStyles();
  const enqueueSnackbar = useEnqueueSnackbar();
  const navigate = useNavigate();
  const [currRow, setCurrRow] = useState<ApnRowType>({} as ApnRowType);
  const [open, setOpen] = React.useState(false);
  const ctx = useContext(ApnContext);
  const apns = ctx.state;
  const apnRows: Array<ApnRowType> = apns
    ? Object.keys(apns).map((apn: string) => {
        const cfg = apns[apn].apn_configuration;
        return {
          apnID: apn,
          classID: cfg?.qos_profile.class_id ?? 0,
          arpPriorityLevel: cfg?.qos_profile.priority_level ?? 0,
          maxReqdULBw: cfg?.ambr.max_bandwidth_ul ?? 0,
          maxReqDLBw: cfg?.ambr.max_bandwidth_dl ?? 0,
          arpPreEmptionCapability:
            cfg?.qos_profile.preemption_capability ?? false,
          arpPreEmptionVulnerability:
            cfg?.qos_profile.preemption_vulnerability ?? false,
        };
      })
    : [];
  const cardActions = {
    buttonText: 'Add APN',
    onClick: () => setOpen(true),
    linkText: 'Learn more about APN',
    link:
      'https://docs.magmacore.org/docs/lte/deploy_config_apn#define-apn-configurations',
  };
  return (
    <div className={classes.dashboardRoot}>
      <>
        <ApnEditDialog
          open={open}
          onClose={() => setOpen(false)}
          apn={Object.keys(currRow).length ? apns[currRow.apnID] : undefined}
        />
        {Object.keys(ctx.state).length > 0 ? (
          <>
            <CardTitleRow key="title" icon={RssFeedIcon} label={APN_TITLE} />
            <ActionTable
              data={apnRows}
              columns={[
                {
                  title: 'Apn ID',
                  field: 'apnID',
                  render: currRow => (
                    <Link
                      variant="body2"
                      component="button"
                      onClick={() => {
                        setCurrRow(currRow);
                        setOpen(true);
                      }}>
                      {currRow.apnID}
                    </Link>
                  ),
                },
                {title: 'Class ID', field: 'classID', type: 'numeric'},
                {
                  title: 'Priority Level',
                  field: 'arpPriorityLevel',
                  type: 'numeric',
                },
                {
                  title: 'Max Reqd UL Bw',
                  field: 'maxReqdULBw',
                  type: 'numeric',
                },
                {title: 'Max Reqd DL Bw', field: 'maxReqDLBw', type: 'numeric'},
                {
                  title: 'Pre-emption Capability',
                  field: 'arpPreEmptionCapability',
                  type: 'numeric',
                },
                {
                  title: 'Pre-emption Vulnerability',
                  field: 'arpPreEmptionVulnerability',
                  type: 'numeric',
                },
              ]}
              handleCurrRow={(row: ApnRowType) => setCurrRow(row)}
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
                    navigate(currRow.apnID + '/json');
                  },
                },
                {name: 'Deactivate'},
                {
                  name: 'Remove',
                  handleFunc: () => {
                    void props
                      .confirm(
                        `Are you sure you want to delete ${currRow.apnID}?`,
                      )
                      .then(async confirmed => {
                        if (!confirmed) {
                          return;
                        }

                        try {
                          // trigger deletion
                          await ctx.setState(currRow.apnID);
                        } catch (e) {
                          enqueueSnackbar(
                            'failed deleting APN ' + currRow.apnID,
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
        ) : (
          <Grid container justify="space-between" spacing={3}>
            <EmptyState
              title={'Set up an APN'}
              instructions={
                'Add an APN to the NMS. The APNs can then be assigned to subscriber profiles.'
              }
              cardActions={cardActions}
              overviewTitle={'APN Overview'}
              overviewDescription={EMPTY_STATE_OVERVIEW}
            />
          </Grid>
        )}
      </>
    </div>
  );
}

export function ApnJsonConfig() {
  const navigate = useNavigate();
  const params = useParams();
  const [error, setError] = useState('');
  const apnName: string = params.apnId!;
  const enqueueSnackbar = useEnqueueSnackbar();
  const ctx = useContext(ApnContext);
  const apns = ctx.state;
  const apn: Apn = apns[apnName] || DEFAULT_APN_CONFIG;
  return (
    <JsonEditor
      content={apn}
      error={error}
      onSave={async apn => {
        try {
          if (apn.apn_name === '') {
            throw Error('Invalid Name');
          }
          await ctx.setState(apn.apn_name, apn);
          enqueueSnackbar('APN saved successfully', {
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

export default withAlert(ApnOverview);
