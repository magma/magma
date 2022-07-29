/*
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

import Card from '@mui/material/Card';
import CardHeader from '@mui/material/CardHeader';
import Collapse from '@mui/material/Collapse';
import DeviceStatusCircle from '../theme/design-system/DeviceStatusCircle';
import Divider from '@mui/material/Divider';
import ExpandLess from '@mui/icons-material/ExpandLess';
import ExpandMore from '@mui/icons-material/ExpandMore';
import Grid from '@mui/material/Grid';
import IconButton from '@mui/material/IconButton';
import Input from '@mui/material/Input';
import InputAdornment from '@mui/material/InputAdornment';
import Link from '@mui/material/Link';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import React from 'react';
import SvgIcon from '@mui/material/SvgIcon';
import Visibility from '@mui/icons-material/Visibility';
import VisibilityOff from '@mui/icons-material/VisibilityOff';
import {Theme} from '@mui/material/styles';
import {colors} from '../theme/default';
import {makeStyles} from '@mui/styles';

const useStyles = makeStyles<Theme, ConfigureStyleParameters>(theme => ({
  dataHeaderBlock: {
    display: 'flex',
    alignItems: 'center',
    padding: 0,
  },
  dataHeaderContent: {
    display: 'flex',
    alignItems: 'center',
  },
  dataHeaderIcon: {
    fill: colors.primary.comet,
    marginRight: theme.spacing(1),
  },
  dataBlock: {
    boxShadow: `0 0 0 1px ${colors.primary.concrete}`,
  },
  dataLabel: {
    color: colors.primary.comet,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
  },
  dataValue: {
    color: colors.primary.brightGray,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    width: props =>
      props.hasStatus
        ? 'calc(100% - 16px)'
        : props.hasIcon
        ? 'calc(100% - 32px)'
        : '100%',
  },
  dataObscuredValue: {
    color: colors.primary.brightGray,
    width: '100%',

    '& input': {
      whiteSpace: 'nowrap',
      overflow: 'hidden',
      textOverflow: 'ellipsis',
    },
  },
  dataBox: {
    width: '100%',
    padding: props => (props.collapsed ? '0' : undefined),

    '& > div': {
      width: '100%',
    },
  },
  dataIcon: {
    display: 'flex',
    alignItems: 'center',

    '& svg': {
      fill: colors.primary.comet,
      marginRight: theme.spacing(1),
    },
  },
  list: {
    padding: 0,
  },
}));

type ConfigureStyleParameters = {
  collapsed?: boolean;
  hasIcon?: boolean;
  hasStatus?: boolean;
};

// Status Indicator displays a small text with an DeviceStatusCircle icon
// disabled indicates if the status color is to be grayed out
// up/down indicates if we have to display status to be in green or in red
function StatusIndicator(disabled: boolean, up: boolean, val: string) {
  const props = {hasStatus: true};
  const classes = useStyles(props);
  return (
    <Grid container alignItems="center">
      <Grid item>
        <DeviceStatusCircle isGrey={disabled} isActive={up} />
      </Grid>
      <Grid item className={classes.dataValue}>
        {val}
      </Grid>
    </Grid>
  );
}

// Data Icon adds an icon to the left of the value
function DataIcon(icon: typeof SvgIcon, val: string) {
  const props = {hasIcon: true};
  const classes = useStyles(props);
  const Icon = icon;
  return (
    <Grid container alignItems="center">
      <Grid item className={classes.dataIcon}>
        <Icon />
      </Grid>
      <Grid item className={classes.dataValue}>
        {val}
      </Grid>
    </Grid>
  );
}

// Data Obscure makes the field into a password type filed with a visibility toggle for more sensitive fields.
function DataObscure(
  value: number | string,
  category: string | null | undefined,
) {
  const [showPassword, setShowPassword] = React.useState(false);
  return (
    <Input
      type={showPassword ? 'text' : 'password'}
      fullWidth={true}
      value={value}
      disableUnderline={true}
      readOnly={true}
      data-testid={`${category ?? value} obscure`}
      endAdornment={
        <InputAdornment position="end">
          <IconButton
            aria-label="toggle password visibility"
            onClick={() => setShowPassword(!showPassword)}
            onMouseDown={event => event.preventDefault()}
            size="large">
            {showPassword ? <Visibility /> : <VisibilityOff />}
          </IconButton>
        </InputAdornment>
      }
    />
  );
}

function DataCollapse(data: Data) {
  const props = {collapsed: true};
  const classes = useStyles(props);
  const [open, setOpen] = React.useState(true);
  const dataEntryValue = `${data.value}${data.unit ?? ''}`;
  return (
    <List
      key={`${data.category ?? data.value}Collapse`}
      className={classes.list}>
      <ListItem button onClick={() => setOpen(!open)}>
        <CardHeader
          data-testid={data.category}
          title={data.category}
          className={classes.dataBox}
          subheader={
            data.statusCircle === true
              ? StatusIndicator(
                  data.statusInactive || false,
                  data.status || false,
                  dataEntryValue,
                )
              : data.icon
              ? DataIcon(data.icon, dataEntryValue)
              : data.obscure === true
              ? DataObscure(data.value, data.category)
              : dataEntryValue
          }
          titleTypographyProps={{
            variant: 'caption',
            className: classes.dataLabel,
            title: data.category,
          }}
          subheaderTypographyProps={{
            variant: 'body1',
            className: classes.dataValue,
            title: data.tooltip ?? dataEntryValue,
          }}
        />
        {open ? <ExpandLess /> : <ExpandMore />}
      </ListItem>
      <Divider />
      <Collapse key={data.value} in={open} timeout="auto" unmountOnExit>
        {data.collapse ?? <></>}
      </Collapse>
    </List>
  );
}

function DataLink(data: string) {
  return (
    <Link href={data} underline="hover">
      {data}
    </Link>
  );
}

type Data = {
  icon?: typeof SvgIcon;
  category?: string;
  value: number | string;
  obscure?: boolean;
  collapse?: React.ReactNode | boolean;
  unit?: string;
  statusCircle?: boolean;
  statusInactive?: boolean;
  status?: boolean;
  tooltip?: string;
  isLink?: boolean;
};

export type DataRows = Array<Data>;

type Props = {
  data: Array<DataRows>;
  testID?: string;
};

export default function DataGrid(props: Props) {
  const classes = useStyles({});
  const dataGrid = props.data.map((row, i) => (
    <Grid key={i} container direction="row">
      {row.map((data, j) => {
        const dataEntryValue = `${data.value}${data.unit ?? ''}`;

        return (
          <React.Fragment key={`data-${i}-${j}`}>
            <Grid
              item
              container
              alignItems="center"
              xs={12}
              md
              key={`data-${i}-${j}`}
              zeroMinWidth
              className={classes.dataBlock}>
              <Grid item xs={12}>
                {data.collapse !== undefined && data.collapse !== false ? (
                  DataCollapse(data)
                ) : (
                  <CardHeader
                    data-testid={data.category}
                    className={classes.dataBox}
                    title={data.category}
                    titleTypographyProps={{
                      variant: 'caption',
                      className: classes.dataLabel,
                      title: data.category,
                    }}
                    subheaderTypographyProps={{
                      variant: 'body1',
                      className:
                        data.obscure === true
                          ? classes.dataObscuredValue
                          : classes.dataValue,
                      title: data.tooltip ?? dataEntryValue,
                    }}
                    subheader={
                      data.statusCircle === true
                        ? StatusIndicator(
                            data.statusInactive || false,
                            data.status || false,
                            dataEntryValue,
                          )
                        : data.icon
                        ? DataIcon(data.icon, dataEntryValue)
                        : data.obscure === true
                        ? DataObscure(data.value, data.category)
                        : data.isLink === true
                        ? DataLink(dataEntryValue)
                        : dataEntryValue
                    }
                  />
                )}
              </Grid>
            </Grid>
          </React.Fragment>
        );
      })}
    </Grid>
  ));
  return (
    <Card elevation={0}>
      <Grid
        container
        alignItems="center"
        justifyContent="center"
        data-testid={props.testID ?? null}>
        {dataGrid}
      </Grid>
    </Card>
  );
}
