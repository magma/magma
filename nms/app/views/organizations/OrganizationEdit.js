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

import type {Organization} from './Organizations';
import type {Tab} from '../../../shared/types/tabs';
import type {WithAlert} from '../../components/Alert/withAlert';

import ArrowBackIcon from '@material-ui/icons/ArrowBack';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import LoadingFiller from '../../components/LoadingFiller';
import OrganizationDialog from './OrganizationDialog';
import OrganizationSummary from './OrganizationSummary';
import OrganizationUsersTable from './OrganizationUsersTable';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
import Text from '../../../app/theme/design-system/Text';
import axios from 'axios';
import withAlert from '../../components/Alert/withAlert';

import {AltFormField} from '../../../app/components/FormField';
import {makeStyles} from '@material-ui/styles';
import {useAxios} from '../../../app/hooks';
import {useCallback, useState} from 'react';
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';
import {useParams} from 'react-router-dom';

const useStyles = makeStyles(_ => ({
  arrowBack: {
    paddingRight: '0px',
    color: 'black',
  },
  container: {
    margin: '40px 32px',
  },
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
    textTransform: 'capitalize',
  },
  titleRow: {
    margin: '16px 0',
  },
}));

export type EditUser = {
  id: string | number,
  email: string,
  role: string,
  networkIDs?: string[],
  organization?: string,
  tabs?: Array<string>,
  password?: string,
  passwordConfirmation?: string,
};

type TitleRowProps = {
  title: string,
  buttonTitle: string,
  onClick: () => void,
};

/**
 * Title and button row
 */
function TitleRow(props: TitleRowProps) {
  const classes = useStyles();
  return (
    <Grid container justifyContent="space-between" className={classes.titleRow}>
      <Text variant="h6">{props.title}</Text>
      <Button variant="text" onClick={() => props.onClick()}>
        <Text variant="body2" color="gray" weight="bold">
          {props.buttonTitle}
        </Text>
      </Button>
    </Grid>
  );
}
type Props = {
  // TODO: remove tabs concept, it should always be NMS
  getProjectTabs?: () => Array<{id: Tab, name: string}>,
  // flag to display advanced config fields in organization add/edit dialog
  hideAdvancedFields?: boolean,
};

function DialogWithConfirmationPhrase(props: DialogConfirmationProps) {
  const [confirmationPhrase, setConfirmationPhrase] = useState('');
  const {title, message, onClose, onConfirm} = props;

  return (
    <Dialog
      open={true}
      onClose={onClose}
      TransitionProps={{onExited: onClose}}
      maxWidth="sm">
      <DialogTitle>{title}</DialogTitle>
      <DialogContent>
        {message}
        <AltFormField label={'Organization Name'}>
          <OutlinedInput
            data-testid="name"
            placeholder="Organization Name"
            fullWidth={true}
            value={confirmationPhrase || ''}
            onChange={({target}) => setConfirmationPhrase(target.value)}
          />
        </AltFormField>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancel</Button>
        <Button
          variant="contained"
          color="primary"
          onClick={onConfirm}
          disabled={confirmationPhrase !== props.confirmationPhrase}>
          Confirm
        </Button>
      </DialogActions>
    </Dialog>
  );
}
/**
 * Organization detail view and Organization edit dialog
 * This component displays an Organization basic information (OrganizationSummary)
 * and its users (OrganizationUsersTable)
 */
function OrganizationEdit(props: WithAlert & Props) {
  const params = useParams();
  const [addingUserFor, setAddingUserFor] = useState<?Organization>(null);
  const classes = useStyles();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [dialog, setDialog] = useState(false);
  const [createError, setCreateError] = useState('');
  const [user, setUser] = useState<?EditUser>(null);
  const [organization, setOrganization] = useState<?Organization>(null);
  const tableRef = React.createRef();
  const [organizationToDelete, setOrganizationToDelete] = useState(null);
  const orgRequest = useAxios<null, {organization: Organization}>({
    method: 'get',
    url: '/host/organization/async/' + params.name,
    onResponse: useCallback(res => {
      setOrganization(res.data.organization);
    }, []),
  });

  const networksRequest = useAxios({
    method: 'get',
    url: '/host/networks/async',
  });

  if (
    orgRequest.isLoading ||
    networksRequest.isLoading ||
    !orgRequest.response
  ) {
    return <LoadingFiller />;
  }

  const onSave = (org: $Shape<OrganizationPlainAttributes>) => {
    axios
      .put('/host/organization/async/' + params.name, org)
      .then(_res => {
        setOrganization(org);
        enqueueSnackbar('Updated organization successfully', {
          variant: 'success',
        });
      })
      .catch(error => {
        const message = error.response?.data?.error || error;
        enqueueSnackbar(`Unable to save organization: ${message}`, {
          variant: 'error',
        });
      });
  };

  return (
    <>
      <div className={classes.container}>
        <OrganizationDialog
          user={user}
          hideAdvancedFields={props.hideAdvancedFields ?? false}
          edit={true}
          organization={organization}
          open={dialog}
          addingUserFor={addingUserFor}
          createError={createError}
          onClose={() => {
            setAddingUserFor(null);
            setDialog(false);
          }}
          onCreateOrg={org => {
            onSave(org);
            setDialog(false);
          }}
          onCreateUser={newUser => {
            if (!user?.id) {
              axios
                .post(
                  `/host/organization/async/${params.name || ''}/add_user`,
                  newUser,
                )
                .then(() => {
                  enqueueSnackbar('User added successfully', {
                    variant: 'success',
                  });
                  // refresh table data
                  tableRef.current?.onQueryChange();
                  setDialog(false);
                })
                .catch(error => {
                  setCreateError(error);
                  enqueueSnackbar(error.message, {
                    variant: 'error',
                  });
                });
            } else {
              axios
                .put(`/user/async/${user.id}`, newUser)
                .then(() => {
                  enqueueSnackbar('User updated successfully', {
                    variant: 'success',
                  });
                  // refresh table data
                  tableRef.current?.onQueryChange();
                  setDialog(false);
                })
                .catch(error => {
                  setCreateError(error);
                  enqueueSnackbar(error?.response?.data?.error || newUser, {
                    variant: 'error',
                  });
                });
            }
          }}
        />
        {(organizationToDelete ?? false) && (
          <DialogWithConfirmationPhrase
            title="Warning"
            message={`Please type the Organization name below to remove it.`}
            label="Organization name"
            confirmationPhrase={organizationToDelete}
            onClose={() => setOrganizationToDelete(null)}
            onConfirm={async () => {
              // delete organization
              await axios.delete(
                `/host/organization/async/${organization?.id || ''}`,
              );
              history.push('/host/organizations');
              setOrganizationToDelete(null);
            }}
          />
        )}
        <Grid container spacing={4}>
          <Grid container justifyContent="space-between">
            <Grid item>
              <Grid container alignItems="center">
                <Grid>
                  <IconButton
                    onClick={() => history.goBack()}
                    className={classes.arrowBack}
                    color="primary">
                    <ArrowBackIcon />
                  </IconButton>
                </Grid>
                <Grid>
                  <Text className={classes.header} variant="h4">
                    {organization?.name}
                  </Text>
                </Grid>
              </Grid>
            </Grid>
            <Grid item>
              <Button
                disabled={organization?.name === 'host'}
                variant="contained"
                color="primary"
                onClick={() => {
                  setOrganizationToDelete(organization?.name);
                }}>
                Remove Organization
              </Button>
            </Grid>
          </Grid>
          <Grid item xs={12}>
            <Grid container spacing={4}>
              <Grid item xs={12} md={6}>
                <TitleRow
                  title={'Organizations'}
                  buttonTitle={'Edit'}
                  onClick={() => {
                    setAddingUserFor(null);
                    setDialog(true);
                  }}
                />
                <OrganizationSummary
                  name={organization?.name}
                  networkIds={organization?.networkIDs}
                />
              </Grid>

              <Grid item xs={12} md={6}>
                <TitleRow
                  title={'Users'}
                  buttonTitle={'Add User'}
                  onClick={() => {
                    setUser(null);
                    setAddingUserFor(organization);
                    setDialog(true);
                  }}
                />
                <OrganizationUsersTable
                  editUser={(newUser: ?EditUser) => {
                    setUser(newUser);
                    setAddingUserFor(organization);
                    setDialog(true);
                  }}
                  tableRef={tableRef}
                />
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </div>
    </>
  );
}

export default withAlert(OrganizationEdit);
