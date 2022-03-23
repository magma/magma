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
 * @flow
 * @format
 */

import * as React from 'react';
import FormAlertsContext from '../Form/FormAlertsContext';
import FormElementContext from '../Form/FormElementContext';
import Grid from '@material-ui/core/Grid';
import ListItem from '@material-ui/core/ListItem';
import Text from '../Text';
import Typography from '@material-ui/core/Typography';
import classNames from 'classnames';
import grey from '@material-ui/core/colors/grey';
import nullthrows from 'nullthrows';
import symphony from '../../../../../fbc_js_core/ui/theme/symphony';

import {makeStyles} from '@material-ui/styles';
import {useContext, useMemo} from 'react';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
  },
  disabled: {
    '& $bottomText': {
      color: symphony.palette.disabled,
    },
  },
  hasError: {
    '& $bottomText': {
      color: symphony.palette.R600,
    },
  },
  labelContainer: {
    marginBottom: '6px',
  },
  bottomText: {
    marginTop: '4px',
    lineHeight: '16px',
  },
  spacer: {
    marginTop: '4px',
    height: '16px',
  },
  subheading: {
    fontWeight: '400',
  },
  optionalLabel: {
    color: grey.A700,
    fontStyle: 'italic',
    fontWeight: '400',
    marginLeft: '8px',
  },
  label: {
    fontSize: '16px',
  },
  children: {
    padding: '8px 0',
  },
}));

export type FormFieldProps = $ReadOnly<{|
  className?: string,
  label?: string,
  helpText?: string,
  children: React.Node,
  disabled?: ?boolean,
  hasError?: boolean,
  required?: boolean,
  errorText?: ?string,
  hasSpacer?: boolean,
  validation?: {
    id: string,
    value: string | number,
  },
  ignorePermissions?: ?boolean,
|}>;

const FormField = (props: FormFieldProps) => {
  const {
    children,
    label,
    helpText,
    disabled: disabledProp,
    className,
    hasError: hasErrorProp,
    errorText: errorTextProp,
    hasSpacer,
    required = false,
    validation,
    ignorePermissions,
  } = props;
  const classes = useStyles();

  const validationContext = useContext(FormAlertsContext);
  const disabled = useMemo(
    () =>
      disabledProp ||
      (validationContext.missingPermissions.detected &&
        ignorePermissions != true) ||
      validationContext.editLock.detected,
    [
      disabledProp,
      ignorePermissions,
      validationContext.editLock.detected,
      validationContext.missingPermissions.detected,
    ],
  );

  const requireFieldError =
    validation == null
      ? ''
      : validationContext.error.check({
          fieldId: validation.id,
          fieldDisplayName: label ?? validation.id,
          value: validation.value,
          required: required,
        });
  const errorText = errorTextProp ?? requireFieldError;
  const hasError = hasErrorProp || !!requireFieldError;
  return (
    <FormElementContext.Provider value={{disabled, hasError}}>
      <div
        className={classNames(
          classes.root,
          {[classes.disabled]: disabled},
          {[classes.hasError]: hasError},
          className,
        )}>
        {label && (
          <Text variant="body2" className={classes.labelContainer}>
            {label}
            {required && ' *'}
          </Text>
        )}
        {children}
        {(helpText || (hasError && errorText)) && (
          <Text className={classes.bottomText} variant="caption">
            {nullthrows((hasError && errorText) || helpText)}
          </Text>
        )}
        {!helpText && !hasError && hasSpacer && (
          <div className={classes.spacer} />
        )}
      </div>
    </FormElementContext.Provider>
  );
};

FormField.defaultProps = {
  disabled: false,
  hasError: false,
  required: false,
};

export default FormField;

type AltFormFieldProps = {
  // Label of the form field
  label: string,
  // Content of the component (Eg, Input, OutlinedInpir, Switch)
  children?: any,
  // If true, compact vertical padding designed for keyboard and mouse input is used
  dense?: boolean,
  // Tooltio of the field
  tooltip?: string,
  // SubLabel of the form field
  subLabel?: string,
  // If true, adds a optional caption to the form field
  isOptional?: boolean,
  // If true, the left and right padding is removed.
  disableGutters?: boolean,
};

export function AltFormField(props: AltFormFieldProps) {
  const classes = useStyles();
  return (
    <ListItem dense={props.dense} disableGutters={props.disableGutters}>
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
