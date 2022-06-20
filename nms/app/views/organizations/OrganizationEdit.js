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
import type {OrganizationPlainAttributes} from '../../../shared/sequelize_models/models/organization';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import ArrowBackIcon from '@material-ui/icons/ArrowBack';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {WithAlert} from '../../components/Alert/withAlert';
// $FlowFixMe migrated to typescript
import LoadingFiller from '../../components/LoadingFiller';
import OrganizationDialog from './OrganizationDialog';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import OrganizationSummary from './OrganizationSummary';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import OrganizationUsersTable from './OrganizationUsersTable';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../../theme/design-system/Text';
import axios from 'axios';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import withAlert from '../../components/Alert/withAlert';
// $FlowFixMe migrated to typescript
import {AltFormField} from '../../components/FormField';
import {makeStyles} from '@material-ui/styles';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useAxios} from '../../hooks';
import {useCallback, useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import {useNavigate, useParams} from 'react-router-dom';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {OrganizationUser} from './types';

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
  // flag to display advanced config fields in organization add/edit dialog
  hideAdvancedFields?: boolean,
};

type DialogConfirmationProps = {
  title: string,
  message: string,
  confirmationPhrase: string,
  onClose: () => void,
  onConfirm: () => void | Promise<void>,
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
  const navigate = useNavigate();
  const [addingUserFor, setAddingUserFor] = useState<?Organization>(null);
  const classes = useStyles();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [dialog, setDialog] = useState(false);
  const [createError, setCreateError] = useState('');
  const [user, setUser] = useState<?OrganizationUser>(null);
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

  if (orgRequest.isLoading || networksRequest.isLoading) {
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
        {organizationToDelete !== null && (
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
              navigate('/host/organizations');
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
                    onClick={() => navigate(-1)}
                    className={classes.arrowBack}
                    color="primary">
                    <ArrowBackIcon />
                  </IconButton>
                </Grid>
                <Grid>
                  <Text
                    className={classes.header}
                    data-testid="organizationEditTitle"
                    variant="h4">
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
                  if (organization) {
                    setOrganizationToDelete(organization.name);
                  }
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
                  editUser={(newUser: ?OrganizationUser) => {
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
