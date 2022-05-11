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
import BusinessIcon from '@material-ui/icons/Business';
import Button from '@material-ui/core/Button';
import CloseIcon from '@material-ui/icons/Close';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import ExitToAppIcon from '@material-ui/icons/ExitToApp';
import Grid from '@material-ui/core/Grid';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import Paper from '@material-ui/core/Paper';
import PersonIcon from '@material-ui/icons/Person';
import Popper from '@material-ui/core/Popper';
import React from 'react';
import Text from '../../../app/theme/design-system/Text';

import {colors} from '../../../app/theme/default';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_ => ({
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
    backgroundColor: colors.primary.comet,
    color: colors.primary.white,
  },
  paper: {
    backgroundColor: '#FFFFFF',
    minHeight: '56px',
    display: 'flex',
    alignItems: 'center',
    minWidth: '350px',
    padding: '16px',
  },
  popperHelperText: {
    fontSize: '18px',
  },
  popperHelperSubtitle: {
    maxWidth: '400px',
    padding: '10px 0',
  },
  popper: {
    zIndex: 1,
    marginTop: '24px',
    '&[x-placement*="bottom-start"] $arrow': {
      '-webkit-filter': 'drop-shadow(0 -1px 1px rgba(0,0,0,0.1))',
      top: 0,
      left: 0,
      marginTop: '-0.9em',
      width: '3em',
      height: '1em',
      '&::before': {
        borderWidth: '0 1em 1em 1em',
        borderColor: `transparent transparent ${colors.primary.white} transparent`,
      },
    },
    '&[x-placement*="left-end"] $arrow': {
      top: '16px!important',
      right: 0,
      marginRight: '-0.9em',
      height: '3em',
      width: '1em',
      '&::before': {
        borderWidth: '1em 0 1em 1em',
        borderColor: `transparent transparent transparent ${colors.primary.white}`,
      },
    },
  },
  arrow: {
    position: 'absolute',
    fontSize: 12,
    width: '3em',
    height: '3em',
    '&::before': {
      content: '""',
      margin: 'auto',
      display: 'block',
      width: 0,
      height: 0,
      borderStyle: 'solid',
    },
  },
}));

type OnboardingDialogType = {
  open: boolean,
  setOpen: boolean => void,
};
export function OnboardingDialog(props: OnboardingDialogType) {
  const classes = useStyles();
  return (
    <Dialog
      classes={{paper: classes.onBoardingDialog}}
      maxWidth={'sm'}
      fullWidth={true}
      open={props.open}
      keepMounted
      onClose={() => props.setOpen(false)}>
      <DialogTitle classes={{root: classes.onBoardingDialogTitle}}>
        Welcome to Magma Host Portal
      </DialogTitle>
      <DialogContent classes={{root: classes.onBoardingDialogContent}}>
        <Text variant="subtitle1">
          In this portal, you can add and edit organizations and its user.
          Follow these steps to get started:
        </Text>
        <List dense={true}>
          <ListItem disableGutters>
            <ListItemIcon>
              <BusinessIcon />
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
              <ExitToAppIcon />
            </ListItemIcon>
            <Text variant="subtitle1">
              Log in to the organization portal with the user account you
              created
            </Text>
          </ListItem>
        </List>
      </DialogContent>
      <DialogActions classes={{root: classes.onBoardingDialogActions}}>
        <Button
          className={classes.onBoardingDialogButton}
          color="primary"
          onClick={() => props.setOpen(false)}>
          Get Started
        </Button>
      </DialogActions>
    </Dialog>
  );
}

export function OnboardingLinkHelper(props: {
  open: boolean,
  linkRef: React.ElementRef<ElementType>,
  onClose: () => void,
}) {
  const [arrowRef, setArrowRef] = React.useState(null);
  const classes = useStyles();
  return (
    <Popper
      className={classes.popper}
      open={props.open && props.linkRef !== null}
      anchorEl={props.linkRef}
      placement={'bottom-start'}
      modifiers={{
        arrow: {
          enabled: true,
          element: arrowRef,
        },
      }}>
      <span className={classes.arrow} ref={setArrowRef} />
      <Paper elevation={2} className={classes.paper}>
        <Grid
          container
          alignContent="center"
          justify="space-around"
          direction="row">
          <Grid>
            <div>
              <Grid container alignContent="center" direction="column">
                <span className={classes.popperHelperText}>
                  Log into the Organization Portal
                </span>
                <span className={classes.popperHelperSubtitle}>
                  Add and manage the Network, Access Gateway, Subscribers, and
                  Policies for the organization.
                </span>
              </Grid>
            </div>
          </Grid>
          <Grid>
            <CloseIcon
              onClick={() => {
                props.onClose();
              }}
            />
          </Grid>
        </Grid>
      </Paper>
    </Popper>
  );
}

export function OnboardingAddButtonHelper(props: {
  open: boolean,
  buttonRef: React.ElementRef<ElementType>,
  onClose: () => void,
}) {
  const [arrowRef, setArrowRef] = React.useState(null);
  const classes = useStyles();
  return (
    <Popper
      className={classes.popper}
      placement={'left-end'}
      open={props.open}
      anchorEl={props.buttonRef}
      modifiers={{
        arrow: {
          enabled: true,
          element: arrowRef,
        },
      }}>
      <span className={classes.arrow} ref={setArrowRef} />
      <Paper elevation={2} className={classes.paper}>
        <Grid container alignContent="center" justifyContent="space-around">
          <span className={classes.popperHelperText}>
            Start by adding an organization
          </span>
          <CloseIcon onClick={() => props.onClose()} />
        </Grid>
      </Paper>
    </Popper>
  );
}
