// @flow strict-local

type ToggleButton = {
  children?: React$Node,
  classes?: Object,
  disabled?: boolean,
  disableFocusRipple?: boolean,
  disableRipple?: boolean,
  selected?: boolean,
  value: any,
};

declare module "@material-ui/lab/ToggleButton" {
  declare module.exports: React$ComponentType<ToggleButton>;
}
declare module "@material-ui/lab/ToggleButtonGroup" {
  declare module.exports: React$ComponentType<{
    children?: React$Node,
    classes?: Object,
    exclusive?: boolean,
    onChange?: Function,
    selected?: boolean | 'auto',
    value?: any,
  }>;
}


type Alert = {
  severity: 'error' | 'info' | 'success' | 'warning',
  variant?: 'filled' | 'outlined' | 'standard',
  children?: React$Node,
  classes?: Object,
}
type AlertTitle = {
  children?: React$Node,
  classes?: Object,
}
declare module "@material-ui/lab/Alert" {
  declare module.exports: React$ComponentType<Alert>;
}
declare module "@material-ui/lab/AlertTitle" {
  declare module.exports: React$ComponentType<AlertTitle>;
}
  
declare module "@material-ui/lab" {
  declare module.exports: {
      Autocomplete: $Exports<"@material-ui/lab/Autocomplete">,
      ToggleButton: $Exports<"@material-ui/lab/ToggleButton">,
      ToggleButtonGroup: $Exports<"@material-ui/lab/ToggleButtonGroup">,
      Alert: $Exports<"@material-ui/lab/Alert">,
      AlertTitle: $Exports<"@material-ui/lab/AlertTitle">,

  };
}
  