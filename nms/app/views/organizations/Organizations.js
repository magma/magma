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

import type {EnqueueSnackbarOptions} from 'notistack';
import type {OrganizationPlainAttributes} from '../../../shared/sequelize_models/models/organization';
import type {UserType} from '../../../shared/sequelize_models/models/user.js';
import type {WithAlert} from '../../components/Alert/withAlert';

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import ActionTable from '../../components/ActionTable';
import BusinessIcon from '@material-ui/icons/Business';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import ExitToAppIcon from '@material-ui/icons/ExitToApp';
import Grid from '@material-ui/core/Grid';
import Link from '@material-ui/core/Link';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
// $FlowFixMe migrated to typescript
import LoadingFiller from '../../components/LoadingFiller';
import OrganizationDialog from './OrganizationDialog';
import PersonAdd from '@material-ui/icons/PersonAdd';
import PersonIcon from '@material-ui/icons/Person';
import React from 'react';
import Text from '../../theme/design-system/Text';
import axios from 'axios';
import withAlert from '../../components/Alert/withAlert';

import {colors} from '../../../app/theme/default';
import {makeStyles} from '@material-ui/styles';
import {useAxios} from '../../../app/hooks';
import {useCallback, useEffect, useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';
import {useNavigate} from 'react-router-dom';

export type Organization = OrganizationPlainAttributes;

const ORGANIZATION_DESCRIPTION =
  'Multiple organizations can be independently managed, each with access to their own networks. ' +
  'As a host user, you can create and manage organizations here. You can also create users for these organizations.';

const useStyles = makeStyles(_ => ({
  addButton: {
    minWidth: '150px',
  },
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
  onBoardingDialog: {
    padding: '24px 0',
  },
  onBoardingDialogTitle: {
    padding: '0 24px',
    fontSize: '24px',
    color: colors.primary.comet,
    backgroundColor: colors.primary.concrete,
  },
  onBoardingDialogContent: {
    minHeight: '200px',
    padding: '16px 24px',
  },
  onBoardingDialogActions: {
    padding: '0 24px',
    backgroundColor: colors.primary.concrete,
    boxShadow: 'none',
  },
  onBoardingDialogButton: {
    minWidth: '120px',
  },
  subtitle: {
    margin: '16px 0',
  },
  index: {
    color: colors.primary.gullGray,
  },
}));

type Props = {...WithAlert};

function OnboardingDialog(props: {setClosed: () => void}) {
  const classes = useStyles();
  return (
    <Dialog
      classes={{paper: classes.onBoardingDialog}}
      maxWidth={'sm'}
      fullWidth={true}
      open={true}
      onClose={() => props.setClosed()}
      aria-describedby="alert-dialog-slide-description">
      <DialogTitle
        data-testid="onboardingDialog"
        classes={{root: classes.onBoardingDialogTitle}}>
        {'Welcome to Magma Host Portal'}
      </DialogTitle>
      <DialogContent classes={{root: classes.onBoardingDialogContent}}>
        <DialogContentText id="alert-dialog-slide-description">
          <Text variant="subtitle1">
            In this portal, you can add and edit organizations and its user.
            Follow these steps to get started:
          </Text>
          <List dense={true}>
            <ListItem disablegutters="true">
              <ListItemIcon>
                <BusinessIcon />
              </ListItemIcon>
              <Text variant="subtitle1">Add an organization</Text>
            </ListItem>
            <ListItem disablegutters="true">
              <ListItemIcon>
                <PersonIcon />
              </ListItemIcon>
              <Text variant="subtitle1">Add a user for the organization</Text>
            </ListItem>
            <ListItem disablegutters="true">
              <ListItemIcon>
                <ExitToAppIcon />
              </ListItemIcon>
              <Text variant="subtitle1">
                Log in to the organization portal with the user account you
                created
              </Text>
            </ListItem>
          </List>
        </DialogContentText>
      </DialogContent>
      <DialogActions classes={{root: classes.onBoardingDialogActions}}>
        <Button
          className={classes.onBoardingDialogButton}
          variant="contained"
          color="primary"
          onClick={() => props.setClosed()}>
          Get Started
        </Button>
      </DialogActions>
    </Dialog>
  );
}

async function getUsers(
  organizations: Organization[],
  setUsers: (Array<?UserType>) => void,
  enqueueSnackbar: (
    msg: string,
    cfg: EnqueueSnackbarOptions,
  ) => ?(string | number),
) {
  const requests = organizations.map(async organization => {
    try {
      const response = await axios.get(
        `/host/organization/async/${organization.name}/users`,
      );
      return response.data;
    } catch (error) {
      enqueueSnackbar(error.message, {
        variant: 'error',
      });
    }
  });
  const organizationUsers = await Promise.all(requests);
  if (organizationUsers) {
    setUsers([...organizationUsers.flat()]);
  }
}

function Organizations(props: Props) {
  const classes = useStyles();
  const navigate = useNavigate();
  const [organizations, setOrganizations] = useState<?(Organization[])>(null);
  const [addingUserFor, setAddingUserFor] = useState<?Organization>(null);
  const [currRow, setCurrRow] = useState<OrganizationRowType>({});
  const [users, setUsers] = useState<Array<?UserType>>([]);
  const [showOnboardingDialog, setShowOnboardingDialog] = useState(false);
  const [showOrganizationDialog, setShowOrganizationDialog] = useState(false);
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
      getUsers(organizations, setUsers, enqueueSnackbar);
    }
  }, [organizations, addingUserFor, enqueueSnackbar]);

  const renderNetworksColumn = useCallback(rowData => {
    // only display 3 networks if more
    if (rowData.networks.length > 3) {
      return `${rowData.networks.slice(0, 3).join(', ')} + ${
        rowData.networks.length - 3
      } more`;
    }
    return rowData.networks.length ? rowData.networks.join(', ') : '-';
  }, []);

  const renderIndexColumn = useCallback(
    rowData => {
      return (
        <Text className={classes.index} variant="caption">
          {rowData.tableData?.index + 1 || ''}
        </Text>
      );
    },
    [classes.index],
  );

  const renderLinkColumn = useCallback(rowData => {
    return (
      <Link href={rowData.portalLink}>
        {`Visit ${rowData.name} Organization Portal`}
      </Link>
    );
  }, []);

  if (error || isLoading || !organizations) {
    return <LoadingFiller />;
  }

  const onDelete = org => {
    props
      .confirm('Are you sure you want to delete this organization?')
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

  const orgName = window.CONFIG.appData.user.tenant;
  const organizationRows: Array<OrganizationRowType> = organizations.map(
    row => {
      return {
        name: row.name,
        networks: Array.from(row.networkIDs || {}),
        portalLink: `${window.location.origin.replace(orgName, row.name)}`,
        userNumber: users?.filter(user => user?.organization === row.name)
          .length,
        id: row.id,
      };
    },
  );

  const menuItems = [
    {
      name: 'View',
      handleFunc: () => {
        navigate(`detail/${currRow.name}`);
      },
    },
  ];
  return (
    <div className={classes.paper}>
      <Grid container>
        <Grid container justifyContent="space-between">
          <Text data-testid="organizationTitle" variant="h3">
            Organizations
          </Text>
          <Button
            className={classes.addButton}
            color="primary"
            variant="contained"
            onClick={() => setShowOrganizationDialog(true)}>
            Add Organization
          </Button>
        </Grid>
        <Grid item xs={12} className={classes.description}>
          <Text variant="body2">{ORGANIZATION_DESCRIPTION}</Text>
        </Grid>
        <>
          {showOnboardingDialog && (
            <OnboardingDialog
              setClosed={() => setShowOnboardingDialog(false)}
            />
          )}
        </>
        <Grid item xs={12}>
          <ActionTable
            data={organizationRows}
            columns={[
              {
                title: '',
                field: '',
                width: '40px',
                editable: 'never',
                render: renderIndexColumn,
              },
              {title: 'Name', field: 'name'},
              {
                title: 'Accessible Networks',
                field: 'networks',
                render: renderNetworksColumn,
              },
              {
                title: 'Link to Organization Portal',
                field: 'portalLink',
                render: renderLinkColumn,
              },
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
                  setShowOrganizationDialog(true);
                },
              },
            ]}
            menuItems={
              currRow.name === 'host'
                ? menuItems
                : [
                    ...menuItems,
                    {
                      name: 'Delete',
                      handleFunc: () => {
                        onDelete(currRow);
                      },
                    },
                  ]
            }
            options={{
              actionsColumnIndex: -1,
              pageSizeOptions: [5, 10],
              toolbar: false,
            }}
          />
        </Grid>
        <OrganizationDialog
          edit={false}
          hideAdvancedFields={false}
          organization={null}
          user={null}
          open={showOrganizationDialog}
          addingUserFor={addingUserFor}
          onClose={() => {
            setShowOrganizationDialog(false);
            setAddingUserFor(null);
          }}
          onCreateOrg={org => {
            axios
              .post('/host/organization/async', org)
              .then(response => {
                enqueueSnackbar('Organization added successfully', {
                  variant: 'success',
                });
                const newOrganization = response.data.organization;
                setOrganizations([...organizations, newOrganization]);
                setAddingUserFor(newOrganization);
              })
              .catch(error => {
                enqueueSnackbar(error.message, {
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
                setShowOrganizationDialog(false);
              })
              .catch(error => {
                enqueueSnackbar(error.message, {
                  variant: 'error',
                });
              });
          }}
        />
      </Grid>
    </div>
  );
}

export default withAlert(Organizations);
