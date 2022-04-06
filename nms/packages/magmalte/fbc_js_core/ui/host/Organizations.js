/**
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

import type {OrganizationPlainAttributes} from '../../../fbc_js_core/sequelize_models/models/organization';
import type {UserType} from '../../../fbc_js_core/sequelize_models/models/user.js';
import type {WithAlert} from '../../../fbc_js_core/ui/components/Alert/withAlert';

import ActionTable from '../components/ActionTable';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import Grid from '@material-ui/core/Grid';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import LoadingFiller from '../../../fbc_js_core/ui/components/LoadingFiller';
import NestedRouteLink from '../../../fbc_js_core/ui/components/NestedRouteLink';
import OrganizationDialog from './OrganizationDialog';
import PersonAdd from '@material-ui/icons/PersonAdd';
import PersonIcon from '@material-ui/icons/Person';
import React from 'react';
import Text from '../components/design-system/Text';
import axios from 'axios';
import withAlert from '../../../fbc_js_core/ui/components/Alert/withAlert';

import {Route} from 'react-router-dom';
import {comet, concrete} from '../../../fbc_js_core/ui/theme/colors';
import {makeStyles} from '@material-ui/styles';
import {useAxios, useRouter} from '../../../fbc_js_core/ui/hooks';
import {useCallback, useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '../../../fbc_js_core/ui/hooks/useSnackbar';
import {useRelativePath, useRelativeUrl} from '../../../fbc_js_core/ui/hooks/useRouter';

export type Organization = OrganizationPlainAttributes;

const ORGANIZATION_DESCRIPTION =
  'Multiple organizations can be independently managed, each with access to their own networks. ' +
  'As a host user, you can create and manage organizations here. You can also create users for these organizations.';

const useStyles = makeStyles(_ => ({
  description: {
    margin: '20px 0',
  },
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  paper: {
    margin: '40px 32px',
  },
  dialogTitle: {
    fontSize: '24px',
    color: comet,
    backgroundColor: concrete,
  },
  dialog: {
    minHeight: '200px',
    padding: '24px',
  },
  dialogActions: {
    backgroundColor: concrete,
    boxShadow: 'none',
  },
  dialogButton: {
    backgroundColor: comet,
    color: concrete,
    '&:hover': {
      backgroundColor: concrete,
      color: comet,
    },
  },

  subtitle: {
    margin: '16px 0',
  },
}));

type Props = {...WithAlert};

function OnboardingDialog() {
  const classes = useStyles();
  const [open, setOpen] = useState(true);
  return (
    <Dialog
      maxWidth={'sm'}
      fullWidth={true}
      open={open}
      keepMounted
      onClose={() => setOpen(false)}
      aria-describedby="alert-dialog-slide-description">
      <DialogTitle classes={{root: classes.dialogTitle}}>
        {'Welcome to Magma Host Portal'}
      </DialogTitle>
      <DialogContent classes={{root: classes.dialog}}>
        <DialogContentText id="alert-dialog-slide-description">
          <Text variant="subtitle1">
            In this portal, you can add and edit organizations and its user.
            Follow these steps to get started:
          </Text>
          <List dense={true}>
            <ListItem disableGutters>
              <ListItemIcon>
                <PersonIcon />
              </ListItemIcon>
              <Text variant="subtitle1">Add an organization</Text>
            </ListItem>
            <ListItem disableGutters>
              <ListItemIcon>
                <PersonIcon />
              </ListItemIcon>
              <Text variant="subtitle1">Add a user for the organization</Text>
            </ListItem>
            <ListItem disableGutters>
              <ListItemIcon>
                <PersonIcon />
              </ListItemIcon>
              <Text variant="subtitle1">
                Log in to the organization portal with the user account you
                created
              </Text>
            </ListItem>
          </List>
        </DialogContentText>
      </DialogContent>
      <DialogActions classes={{root: classes.dialogActions}}>
        <Button
          classes={{root: classes.dialogButton}}
          onClick={() => setOpen(false)}>
          Get Started
        </Button>
      </DialogActions>
    </Dialog>
  );
}

async function getUsers(
  organizations: Organization[],
  setUsers: (Array<?UserType>) => void,
) {
  const requests = organizations.map(async organization => {
    try {
      const response = await axios.get(
        `/host/organization/async/${organization.name}/users`,
      );
      return response.data;
    } catch (error) {}
  });
  const organizationUsers = await Promise.all(requests);
  if (organizationUsers) {
    setUsers([...organizationUsers.flat()]);
  }
}

function Organizations(props: Props) {
  const classes = useStyles();
  const relativeUrl = useRelativeUrl();
  const relativePath = useRelativePath();
  const {history} = useRouter();
  const [organizations, setOrganizations] = useState<?(Organization[])>(null);
  const [addingUserFor, setAddingUserFor] = useState<?Organization>(null);
  const [currRow, setCurrRow] = useState<OrganizationRowType>({});
  const [users, setUsers] = useState<Array<?UserType>>([]);
  const [showOnboardingDialog, setShowOnboardingDialog] = useState(false);
  const [addUser, setAddUser] = useState(false);
  const enqueueSnackbar = useEnqueueSnackbar();
  const {error, isLoading} = useAxios({
    url: '/host/organization/async',
    onResponse: useCallback(res => {
      setOrganizations(res.data.organizations);
      if (res.data.organizations.length < 3) {
        setShowOnboardingDialog(true);
      }
    }, []),
  });
  useEffect(() => {
    if (organizations?.length) {
      getUsers(organizations, setUsers);
    }
  }, [organizations, addingUserFor]);

  if (error || isLoading || !organizations) {
    return <LoadingFiller />;
  }

  const onDelete = org => {
    props
      .confirm('Are you sure you want to delete this org?')
      .then(async confirm => {
        if (!confirm) return;
        await axios.delete(`/host/organization/async/${org.id}`);
        const newOrganizations = organizations.filter(
          organization => organization.id !== org.id,
        );
        setOrganizations([...newOrganizations]);
      });
  };

  type OrganizationRowType = {
    name: string,
    networks: Array<string>,
    portalLink: string,
    userNumber: number,
    id: number,
  };

  const organizationRows: Array<OrganizationRowType> = organizations.map(
    row => {
      return {
        name: row.name,
        networks: row.networkIDs,
        portalLink: `${row.name}`,
        userNumber: users?.filter(user => user?.organization === row.name)
          .length,
        id: row.id,
      };
    },
  );
  return (
    <div className={classes.paper}>
      <Grid container>
        <Grid container justify="space-between">
          <Text variant="h3">Organizations</Text>
          <NestedRouteLink to="/new">
            <Button color="primary" variant="contained">
              Add Organization
            </Button>
          </NestedRouteLink>
        </Grid>
        <Grid xs={12} className={classes.description}>
          <Text variant="body2">{ORGANIZATION_DESCRIPTION}</Text>
        </Grid>
        <>{showOnboardingDialog && <OnboardingDialog />}</>
        <Grid xs={12}>
          <ActionTable
            data={organizationRows}
            columns={[
              {
                title: '',
                field: '',
                width: '40px',
                editable: 'never',
                render: rowData => (
                  <Text variant="subtitle3">
                    {rowData.tableData?.id + 1 || ''}
                  </Text>
                ),
              },
              {title: 'Name', field: 'name'},
              {title: 'Accessible Networks', field: 'networks'},
              {title: 'Link to Organization Portal', field: 'portalLink'},
              {title: 'Number of Users', field: 'userNumber'},
            ]}
            handleCurrRow={(row: OrganizationRowType) => {
              setCurrRow(row);
            }}
            actions={[
              {
                icon: () => <PersonAdd />,
                tooltip: 'Add User',
                onClick: (event, row) => {
                  setAddingUserFor(row);
                },
              },
            ]}
            menuItems={[
              {
                name: 'View',
                handleFunc: () => {
                  history.push(relativeUrl(`/detail/${currRow.name}`));
                },
              },
              {
                name: 'Delete',
                handleFunc: () => {
                  onDelete(currRow);
                },
              },
            ]}
            options={{
              actionsColumnIndex: -1,
              pageSizeOptions: [5, 10],
              toolbar: false,
            }}
          />
        </Grid>
        <Route
          path={relativePath('/new')}
          render={() => (
            <OrganizationDialog
              addUser={addUser}
              setAddUser={() => setAddUser(true)}
              onClose={() => {
                setAddUser(false);
                history.push(relativeUrl(''));
              }}
              onCreateOrg={org => {
                let newOrg = null;
                axios
                  .post('/host/organization/async', org)
                  .then(() => {
                    enqueueSnackbar('Organization added successfully', {
                      variant: 'success',
                    });
                    axios
                      .get(`/host/organization/async/${org.name}`)
                      .then(resp => {
                        newOrg = resp.data.organization;
                        if (newOrg) {
                          setOrganizations([...organizations, newOrg]);
                          setAddingUserFor(newOrg);
                        }
                      });
                  })
                  .catch(error => {
                    setAddUser(false);
                    history.push(relativeUrl(''));
                    enqueueSnackbar(error?.response?.data?.error || error, {
                      variant: 'error',
                    });
                  });
              }}
              onCreateUser={user => {
                axios
                  .post(
                    `/host/organization/async/${
                      addingUserFor?.name || ''
                    }/add_user`,
                    user,
                  )
                  .then(() => {
                    enqueueSnackbar('User added successfully', {
                      variant: 'success',
                    });
                    setAddingUserFor(null);
                    history.push(relativeUrl(''));
                  })
                  .catch(error => {
                    enqueueSnackbar(error?.response?.data?.error || error, {
                      variant: 'error',
                    });
                  });
              }}
            />
          )}
        />
      </Grid>
    </div>
  );
}

export default withAlert(Organizations);
