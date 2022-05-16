/**
 * Copyright 2022 The Magma Authors.
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
 * @flow
 * @format
 */

import ApnContext from './context/ApnContext';
import Button from '@material-ui/core/Button';
import GatewayContext from './context/GatewayContext';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import LockIcon from '@material-ui/icons/Lock';
import NetworkContext from './context/NetworkContext';
import PlaylistAddCheckIcon from '@material-ui/icons/PlaylistAddCheck';
import Popout from './Popout';
import React, {useCallback, useContext, useMemo} from 'react';
import SubscriberContext from './context/SubscriberContext';
import Text from '../theme/design-system/Text';
import classNames from 'classnames';
import {colors} from '../theme/default';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  button: {
    display: 'flex',
    alignItems: 'center',
    width: '100%',
    padding: '15px 28px',
    cursor: 'pointer',
    outline: 'none',
    '&:hover $icon, &:hover $label, &$selected $icon, &$selected $label': {
      color: colors.primary.white,
    },
  },
  label: {
    '&&': {
      color: colors.primary.gullGray,
      whiteSpace: 'nowrap',
      paddingLeft: '16px',
    },
  },
  selected: {
    backgroundColor: colors.secondary.dodgerBlue,

    '& $icon': {
      color: colors.primary.white,
    },
  },
  icon: {
    color: colors.primary.gullGray,
    display: 'flex',
    justifyContent: 'center',
  },
  itemGutters: {
    '&&': {
      minWidth: '200px',
      padding: '6px 17px',
    },
  },
  divider: {
    margin: '6px 17px',
  },
  profileList: {
    marginTop: '6px',
    backgroundColor: colors.primary.white,
    '&&': {
      padding: '10px 0',
    },
  },
  profileItemText: {
    fontSize: '14px',
    lineHeight: '20px',
  },
  test: {
    backgroundColor: colors.primary.mercury,
    padding: '20px',
    minWidth: '420px',
  },
  onBoardingAction: {
    padding: '0',
  },
  icon: {
    color: colors.primary.nobel,
  },
  completedAction: {
    backgroundColor: colors.primary.selago,
    color: colors.state.positiveAlt,
    cursor: 'default',
  },
}));

type Props = {
  isMenuOpen: boolean,
  setMenuOpen: (isOpen: boolean) => void,
  expanded: boolean,
};

const OnBoardingChecklist = (props: Props) => {
  const classes = useStyles();
  const gatewayCtx = useContext(GatewayContext);
  const networkCtx = useContext(NetworkContext);
  const apnCtx = useContext(ApnContext);
  const subscriberCtx = useContext(SubscriberContext);
  const onboardingActions = useMemo(() => {
    return {
      requiredOnboardingActions: [
        'Add Network',
        'Set Up Access Gateway',
        'Add APN',
        'Add Subscribers',
      ],
      recommendedOnboardingActions: [
        'Enable Alerts',
        'Add eNodeBs',
        'Add Policies',
      ],
    };
  }, []);

  const isActionCompleted = useCallback(
    (actionTitle: string) => {
      switch (actionTitle) {
        case 'Add Network':
          return !!(networkCtx.networkId ?? false);
        case 'Set Up Access Gateway':
          return Object.keys(gatewayCtx?.state || {})?.length > 0;
        case 'Add APN':
          return Object.keys(apnCtx?.state || {})?.length > 0;
        case 'Add Subscribers':
          return Object.keys(subscriberCtx?.state || {})?.length > 0;
        default:
          return false;
      }
    },
    [
      gatewayCtx?.state,
      networkCtx.networkId,
      apnCtx?.state,
      subscriberCtx?.state,
    ],
  );

  const getCurrentAction = useCallback(
    (actions: Array<string>) => {
      let currentAction = null;
      actions.some(action => {
        if (!isActionCompleted(action)) {
          currentAction = action;
          return true;
        }
      });
      return currentAction;
    },
    [isActionCompleted],
  );

  return (
    <Popout
      className={classNames({
        [classes.button]: true,
      })}
      open={props.isMenuOpen}
      content={
        <Grid className={classes.test} container direction="column" spacing={2}>
          <Grid item>
            <Text>Getting Started Guide</Text>
          </Grid>

          <Grid item>
            <Text variant="subtitle3">Required</Text>
            <List className={classes.profileList}>
              {onboardingActions.requiredOnboardingActions.map(
                (onBoardingAction, _index) => {
                  return (
                    <ListItem
                      key={onBoardingAction}
                      classes={{gutters: classes.itemGutters}}>
                      <Text
                        className={classes.onBoardingAction}
                        variant="subtitle2"
                        weight="regular">
                        {onBoardingAction}
                      </Text>
                      <ListItemSecondaryAction>
                        {!isActionCompleted(onBoardingAction) ? (
                          getCurrentAction(
                            onboardingActions.requiredOnboardingActions,
                          ) === onBoardingAction ? (
                            <Button color="primary">Get Started</Button>
                          ) : (
                            <IconButton className={classes.icon} edge="end">
                              <LockIcon />
                            </IconButton>
                          )
                        ) : (
                          <Button className={classes.completedAction}>
                            Completed
                          </Button>
                        )}
                      </ListItemSecondaryAction>
                    </ListItem>
                  );
                },
              )}
            </List>
          </Grid>
          <Grid item>
            <Text variant="subtitle3">Recommended</Text>
            <List className={classes.profileList}>
              {onboardingActions.recommendedOnboardingActions.map(
                onBoardingAction => {
                  return (
                    <ListItem
                      key={onBoardingAction}
                      classes={{gutters: classes.itemGutters}}>
                      <Text
                        className={classes.onBoardingAction}
                        variant="subtitle2"
                        weight="regular">
                        {onBoardingAction}
                      </Text>
                      <ListItemSecondaryAction>
                        {!isActionCompleted(onBoardingAction) ? (
                          getCurrentAction(
                            onboardingActions.recommendedOnboardingActions,
                          ) === onBoardingAction &&
                          !getCurrentAction(
                            onboardingActions.requiredOnboardingActions,
                          ) ? (
                            <Button color="primary">Get Started</Button>
                          ) : (
                            <IconButton className={classes.icon} edge="end">
                              <LockIcon />
                            </IconButton>
                          )
                        ) : (
                          <Button className={classes.completedAction}>
                            Completed
                          </Button>
                        )}
                      </ListItemSecondaryAction>
                    </ListItem>
                  );
                },
              )}
            </List>
          </Grid>
        </Grid>
      }
      onOpen={() => props.setMenuOpen(true)}
      onClose={() => props.setMenuOpen(false)}>
      <PlaylistAddCheckIcon
        data-testid="profileButton"
        className={classes.icon}
      />
      {props.expanded && (
        <Text className={classes.label} variant="body3">
          CheckList
        </Text>
      )}
    </Popout>
  );
};

export default OnBoardingChecklist;
