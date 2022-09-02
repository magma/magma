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

import Box from '@mui/material/Box';
import Divider from '@mui/material/Divider';
import Grid from '@mui/material/Grid';
import HelpIcon from '@mui/icons-material/Help';
import IconButton from '@mui/material/IconButton';
import InputAdornment from '@mui/material/InputAdornment';
import LinearProgress from '@mui/material/LinearProgress';
import ListItem from '@mui/material/ListItem';
import OutlinedInput from '@mui/material/OutlinedInput';
import React from 'react';
import Text from '../theme/design-system/Text';
import Tooltip from '@mui/material/Tooltip';
import Typography from '@mui/material/Typography';
import Visibility from '@mui/icons-material/Visibility';
import VisibilityOff from '@mui/icons-material/VisibilityOff';
import {Theme} from '@mui/material/styles';
import {colors} from '../theme/default';
import {grey} from '@mui/material/colors';
import {makeStyles} from '@mui/styles';
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
  label: {
    fontSize: '14px',
    fontFamily: 'Inter',
    fontStyle: 'normal',
    fontWeight: 500,
    lineHeight: '20px',
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
    margin: `${theme.spacing(3)} 0 ${theme.spacing(2)}`,
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
  className?: string;
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
    <ListItem
      dense={props.dense}
      disableGutters={props.disableGutters}
      className={props.className}>
      <Grid container>
        <Grid item xs={12} className={classes.label}>
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
  className?: string;
  'data-testid'?: string;
  error?: boolean;
  onChange: (onChange: string) => void;
  placeholder?: string;
  value: string;
  required?: boolean;
  fullWidth?: boolean;
  autoComplete?: string;
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
            edge="end"
            size="large">
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
