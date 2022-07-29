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
 */

import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import Grid from '@mui/material/Grid';
import IconButton from '@mui/material/IconButton';
import LoadingFiller from '../../components/LoadingFiller';
import OrganizationDialog from './OrganizationDialog';
import OrganizationSummary from './OrganizationSummary';
import OrganizationUsersTable from './OrganizationUsersTable';
import OutlinedInput from '@mui/material/OutlinedInput';
import React from 'react';
import Text from '../../theme/design-system/Text';
import axios, {AxiosResponse} from 'axios';
import withAlert from '../../components/Alert/withAlert';
import {AltFormField} from '../../components/FormField';
import {OrganizationUser} from './types';
import {TableRef} from '../../components/ActionTable';
import {getErrorMessage} from '../../util/ErrorUtils';
import {makeStyles} from '@mui/styles';
import {useAxios} from '../../hooks';
import {useCallback, useState} from 'react';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import {useNavigate, useParams} from 'react-router-dom';
import type {Organization} from './Organizations';
import type {OrganizationRawType} from '../../../shared/sequelize_models/models/organization';
import type {WithAlert} from '../../components/Alert/withAlert';

const useStyles = makeStyles({
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
});

type TitleRowProps = {
  title: string;
  buttonTitle: string;
  onClick: () => void;
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
        <Text variant="body2" weight="bold">
          {props.buttonTitle}
        </Text>
      </Button>
    </Grid>
  );
}
type Props = {
  // flag to display advanced config fields in organization add/edit dialog
  hideAdvancedFields?: boolean;
};

type DialogConfirmationProps = {
  title: string;
  message: string;
  confirmationPhrase: string;
  onClose: () => void;
  onConfirm: () => void | Promise<void>;
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
          onClick={() => void onConfirm()}
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
  const [addingUserFor, setAddingUserFor] = useState<Organization | null>(null);
  const classes = useStyles();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [dialog, setDialog] = useState(false);
  const [user, setUser] = useState<OrganizationUser | null>(null);
  const [organization, setOrganization] = useState<Organization | null>(null);
  const tableRef: TableRef = React.useRef();
  const [organizationToDelete, setOrganizationToDelete] = useState<
    string | null
  >(null);
  const orgRequest = useAxios<{organization: Organization}>({
    method: 'get',
    url: '/host/organization/async/' + params.name!,
    onResponse: useCallback(
      (res: AxiosResponse<{organization: Organization}>) => {
        setOrganization(res.data.organization);
      },
      [],
    ),
  });

  const networksRequest = useAxios({
    method: 'get',
    url: '/host/networks/async',
  });

  if (orgRequest.isLoading || networksRequest.isLoading) {
    return <LoadingFiller />;
  }

  const onSave = (org: Partial<OrganizationRawType>) => {
    axios
      .put('/host/organization/async/' + params.name!, org)
      .then(() => {
        setOrganization(org as Organization);
        enqueueSnackbar('Updated organization successfully', {
          variant: 'success',
        });
      })
      .catch(error => {
        const message = getErrorMessage(error);
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
          organization={organization}
          open={dialog}
          addingUserFor={addingUserFor}
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
                  enqueueSnackbar(getErrorMessage(error), {
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
                  enqueueSnackbar(getErrorMessage(error), {
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
                    color="primary"
                    size="large">
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
                  editUser={(newUser: OrganizationUser | null) => {
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
