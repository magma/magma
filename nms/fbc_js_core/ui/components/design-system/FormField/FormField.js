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
import Text from '../../../../../app/theme/design-system/Text';
import classNames from 'classnames';
import grey from '@material-ui/core/colors/grey';
import nullthrows from 'nullthrows';

import {colors} from '../../../../../app/theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext, useMemo} from 'react';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
  },
  disabled: {
    '& $bottomText': {
      color: colors.primary.gullGray,
    },
  },
  hasError: {
    '& $bottomText': {
      color: colors.state.error,
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
