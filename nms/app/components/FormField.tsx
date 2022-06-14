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

import Box from '@material-ui/core/Box';
import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import HelpIcon from '@material-ui/icons/Help';
import IconButton from '@material-ui/core/IconButton';
import InputAdornment from '@material-ui/core/InputAdornment';
import LinearProgress from '@material-ui/core/LinearProgress';
import ListItem from '@material-ui/core/ListItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
import Text from '../theme/design-system/Text';
import Tooltip from '@material-ui/core/Tooltip';
import Typography from '@material-ui/core/Typography';
import Visibility from '@material-ui/icons/Visibility';
import VisibilityOff from '@material-ui/icons/VisibilityOff';
import grey from '@material-ui/core/colors/grey';
import {Theme} from '@material-ui/core/styles';
import {colors} from '../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useState} from 'react';

const useStyles = makeStyles<Theme>(theme => ({
  children: {
    padding: '8px 0',
  },
  root: {
    display: 'flex',
    marginBottom: '5px',
    alignItems: 'center',
  },
  heading: {
    flexBasis: '33.33%',
    marginRight: '15px',
    textAlign: 'right',
  },
  subheading: {
    fontWeight: 400,
  },
  optionalLabel: {
    color: grey.A700,
    fontStyle: 'italic',
    fontWeight: 400,
    marginLeft: '8px',
  },
  secondaryHeading: {
    flexBasis: '66.66%',
  },
  icon: {
    marginLeft: '5px',
    paddingTop: '4px',
    verticalAlign: 'bottom',
    width: '15px',
  },
  formDivider: {
    margin: `${theme.spacing(3)}px 0 ${theme.spacing(2)}px`,
    backgroundColor: colors.primary.gullGray,
    opacity: 0.4,
    height: '1px',
  },
}));

type Props = {
  // Label of the form field
  label: string;
  // Content of the component (Eg, Input, OutlinedInpir, Switch)
  children?: any;
  // If true, compact vertical padding designed for keyboard and mouse input is used
  dense?: boolean;
  // Tooltio of the field
  tooltip?: string;
  // SubLabel of the form field
  subLabel?: string;
  // If true, adds a optional caption to the form field
  isOptional?: boolean;
  // If true, the left and right padding is removed.
  disableGutters?: boolean;
};

export default function FormField(props: Props) {
  const classes = useStyles();
  const {tooltip} = props;
  return (
    <div className={classes.root}>
      <Text className={classes.heading} variant="body2">
        {props.label}
        {tooltip && (
          <Tooltip title={tooltip} placement="bottom-start">
            <HelpIcon className={classes.icon} />
          </Tooltip>
        )}
      </Text>
      <Typography
        className={classes.secondaryHeading}
        component="div"
        variant="body2">
        {props.children}
      </Typography>
    </div>
  );
}

export function AltFormField(props: Props) {
  const classes = useStyles();
  return (
    <ListItem dense={props.dense} disableGutters={props.disableGutters}>
      <Grid container>
        <Grid item xs={12}>
          {props.label}
          {props.isOptional && (
            <Typography
              className={classes.optionalLabel}
              variant="caption"
              gutterBottom>
              {'optional'}
            </Typography>
          )}
        </Grid>
        {props.subLabel && (
          <Grid item xs={12}>
            <Typography
              className={classes.subheading}
              variant="caption"
              display="block"
              gutterBottom>
              {props.subLabel}
            </Typography>
          </Grid>
        )}
        <Grid item xs={12} className={classes.children}>
          {props.children}
        </Grid>
      </Grid>
    </ListItem>
  );
}

export function AltFormFieldSubheading(props: Props) {
  const classes = useStyles();
  return (
    <Grid container>
      <Grid item xs={12}>
        <Typography
          className={classes.subheading}
          variant="caption"
          display="block"
          gutterBottom>
          {props.label}
        </Typography>
      </Grid>
      <Grid item xs={12}>
        {props.children}
      </Grid>
    </Grid>
  );
}

type PasswordProps = {
  value: string;
  onChange: (onChange: string) => void;
  placeholder?: string;
};

export function FormDivider() {
  const classes = useStyles();
  return <Divider className={classes.formDivider} />;
}

export function PasswordInput(props: PasswordProps) {
  const [showPassword, setShowPassword] = useState(false);
  return (
    <OutlinedInput
      {...props}
      type={showPassword ? 'text' : 'password'}
      value={props.value}
      onChange={e => props.onChange(e.target.value)}
      endAdornment={
        <InputAdornment position="end">
          <IconButton
            aria-label="toggle password visibility"
            onClick={() => setShowPassword(true)}
            onMouseDown={() => setShowPassword(false)}
            edge="end">
            {showPassword ? <Visibility /> : <VisibilityOff />}
          </IconButton>
        </InputAdornment>
      }
    />
  );
}

type ProgressProps = {
  value: number;
  text?: string;
};

export function LinearProgressWithLabel(props: ProgressProps) {
  return (
    <Box display="flex" alignItems="center">
      <Box width="100%" mr={1}>
        <LinearProgress variant="determinate" value={props.value} />
      </Box>
      <Box minWidth={35}>
        <Text>{props.text ?? `${Math.round(props.value)}%`}</Text>
      </Box>
    </Box>
  );
}
