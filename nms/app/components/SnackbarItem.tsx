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

import * as React from 'react';
import CheckCircleIcon from '@material-ui/icons/CheckCircle';
import CloseIcon from '@material-ui/icons/Close';
import ErrorIcon from '@material-ui/icons/Error';
import InfoIcon from '@material-ui/icons/Info';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../theme/design-system/Text';
import WarningIcon from '@material-ui/icons/Warning';
import classNames from 'classnames';
import {colors} from '../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useSnackbar} from 'notistack';
import type {VariantType} from 'notistack';

const useStyles = makeStyles(() => ({
  root: {
    backgroundColor: colors.primary.white,
    boxShadow: '0 0 0 1px #ccd0d5, 0 4px 8px 1px rgba(0,0,0,0.15)',
    borderRadius: '2px',
    display: 'flex',
    width: '420px',
  },
  bar: {
    borderLeft: '6px solid',
    borderRadius: '2px 0px 0px 2px',
  },
  errorBar: {
    borderColor: colors.state.error,
  },
  successBar: {
    borderColor: colors.state.positive,
  },
  warningBar: {
    borderColor: colors.state.warning,
  },
  defaultBar: {
    borderColor: colors.primary.comet,
  },
  infoBar: {
    borderColor: colors.secondary.dodgerBlue,
  },
  content: {
    marginLeft: '6px',
    display: 'flex',
    flexDirection: 'row',
    padding: '12px 12px 12px 0px',
    alignItems: 'center',
    flexGrow: 1,
  },
  message: {
    marginLeft: '6px',
    fontSize: '13px',
    lineHeight: '17px',
    flexGrow: 1,
  },
  icon: {
    '&&': {fontSize: '20px'},
    marginRight: '12px',
  },
  errorIcon: {
    '&&': {fill: colors.state.error},
  },
  successIcon: {
    '&&': {fill: colors.state.positive},
  },
  warningIcon: {
    '&&': {fill: colors.state.warning},
  },
  defaultIcon: {
    '&&': {fill: colors.primary.comet},
  },
  infoIcon: {
    '&&': {fill: colors.secondary.dodgerBlue},
  },
  closeButton: {
    marginLeft: '16px',
    color: colors.primary.gullGray,
    '&:hover': {
      color: colors.primary.comet,
    },
    cursor: 'pointer',
  },
}));

type Props = {
  id: number | string;
  message: string;
  variant: VariantType;
};

const IconVariants: Record<VariantType, React.ComponentType<any>> = {
  error: ErrorIcon,
  success: CheckCircleIcon,
  warning: WarningIcon,
  default: InfoIcon,
  info: InfoIcon,
};

const SnackbarItem = React.forwardRef<HTMLDivElement, Props>(
  (props, fwdRef) => {
    const {id, message, variant} = props;
    const classes = useStyles();
    const {closeSnackbar} = useSnackbar();
    const Icon = IconVariants[variant];
    return (
      <div className={classes.root} ref={fwdRef}>
        <div className={classNames(classes.bar, classes[`${variant}Bar`])} />
        <div className={classes.content}>
          <Icon
            className={classNames(classes.icon, classes[`${variant}Icon`])}
          />
          <Text className={classes.message}>{message}</Text>
          <CloseIcon
            className={classes.closeButton}
            onClick={() => closeSnackbar(id)}
          />
        </div>
      </div>
    );
  },
);

export default SnackbarItem;
