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
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';

import ActionTable from '../../components/ActionTable';
import ApnContext from '../../components/context/ApnContext';
import ApnEditDialog from './ApnEdit';
import CardTitleRow from '../../components/layout/CardTitleRow';
import JsonEditor from '../../components/JsonEditor';
import Link from '@material-ui/core/Link';
import React from 'react';
import RssFeedIcon from '@material-ui/icons/RssFeed';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';

import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

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
  apnID: string,
  description: string,
  qosProfile: number,
  added: Date,
};

const APN_TITLE = 'APNs';
function ApnOverview(props: WithAlert) {
  const classes = useStyles();
  const enqueueSnackbar = useEnqueueSnackbar();
  const {history, relativeUrl} = useRouter();
  const [currRow, setCurrRow] = useState<ApnRowType>({});
  const [open, setOpen] = React.useState(false);
  const ctx = useContext(ApnContext);
  const apns = ctx.state;
  const apnRows: Array<ApnRowType> = apns
    ? Object.keys(apns).map((apn: string) => {
        return {
          apnID: apn,
          description: 'Test APN description',
          qosProfile: 1,
          added: new Date(0),
        };
      })
    : [];
  return (
    <div className={classes.dashboardRoot}>
      <>
        <CardTitleRow key="title" icon={RssFeedIcon} label={APN_TITLE} />
        <ApnEditDialog
          open={open}
          onClose={() => setOpen(false)}
          apn={Object.keys(currRow).length ? apns[currRow.apnID] : undefined}
        />
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
            {title: 'Description', field: 'description'},
            {title: 'Qos Profile', field: 'qosProfile', type: 'numeric'},
            {title: 'Added', field: 'added', type: 'datetime'},
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
                history.push(relativeUrl('/' + currRow.apnID + '/json'));
              },
            },
            {name: 'Deactivate'},
            {
              name: 'Remove',
              handleFunc: () => {
                props
                  .confirm(`Are you sure you want to delete ${currRow.apnID}?`)
                  .then(async confirmed => {
                    if (!confirmed) {
                      return;
                    }

                    try {
                      // trigger deletion
                      ctx.setState(currRow.apnID);
                    } catch (e) {
                      enqueueSnackbar('failed deleting APN ' + currRow.apnID, {
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
    </div>
  );
}

export function ApnJsonConfig() {
  const {match, history} = useRouter();
  const [error, setError] = useState('');
  const apnName: string = match.params.apnId;
  const enqueueSnackbar = useEnqueueSnackbar();
  const ctx = useContext(ApnContext);
  const apns = ctx.state;
  const apn: apn = apns[apnName] || DEFAULT_APN_CONFIG;
  return (
    <JsonEditor
      content={apn}
      error={error}
      onSave={async apn => {
        try {
          if (apn.apn_name === '') {
            throw Error('Invalid Name');
          }
          ctx.setState(apn.apn_name, apn);
          enqueueSnackbar('APN saved successfully', {
            variant: 'success',
          });
          setError('');
          history.goBack();
        } catch (e) {
          setError(e.response?.data?.message ?? e.message);
        }
      }}
    />
  );
}

export default withAlert(ApnOverview);
