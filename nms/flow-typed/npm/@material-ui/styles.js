declare module "@material-ui/styles/ThemeProvider" {
  import type {ComponentType} from "react";
  declare module.exports: ComponentType<Object>;
}

declare module "@material-ui/styles" {
  import type {ComponentType} from "react";
  declare module.exports: {
    makeStyles: (color: Object) => (props: any) => any,
    StylesProvider: ComponentType<{generateClassName:()=>string}>,
    createGenerateClassName:() => () => string
  };
}
