/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * @format
 * @flow strict-local
 */

type KeyboardDatePickerProps = {
  disableToolbar?: boolean,
  inputVariant?: string,
  format?: string,
  margin?: string,
  value?: string | Date,
  onChange?: Object => any,
  KeyboardButtonProps?: Object,
};

type DatePickerProps = {
  autoOk?: boolean,
  disableFuture?: boolean,
  variant?: string,
  inputVariant?: string,
  inputProps?: Object,
  maxDate?: ParsableDate,
  format?: string,
  value?: moment$Moment,
  onChange?: Object => any,
};

declare module '@material-ui/pickers/KeyboardDatePicker' {
  declare module.exports: React$ComponentType<KeyboardDatePickerProps>;
}

declare module '@material-ui/pickers/DateTimePicker' {
  declare module.exports: React$ComponentType<DatePickerProps>;
}

declare module '@material-ui/pickers/MuiPickersUtilsProvider' {
  declare module.exports: React$ComponentType<KeyboardDatePickerProps>;
}

declare module '@material-ui/pickers' {
  declare module.exports: {
    KeyboardDatePicker: $Exports<'@material-ui/pickers/KeyboardDatePicker'>,
    KeyboardTimePicker: $Exports<'@material-ui/pickers/KeyboardTimePicker'>,
    DateTimePicker: $Exports<'@material-ui/pickers/DateTimePicker'>,
    MuiPickersUtilsProvider: $Exports<
      '@material-ui/pickers/MuiPickersUtilsProvider',
    >,
  };
}
