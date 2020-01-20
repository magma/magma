declare module "@material-ui/styles/ThemeProvider" {
  import type {ComponentType} from "react";
  declare module.exports: ComponentType<Object>;
}

declare module "@material-ui/styles" {
  import type {ComponentType, Node} from "react";
  import type {Theme} from '@material-ui/core';
  declare class ServerStyleSheets {
    collect:(Node) => Node,
  }
  declare module.exports: {
    makeStyles: (color: Object) => (props: any) => any,
    StylesProvider: ComponentType<{generateClassName?:()=>string}>,
    ServerStyleSheets: Class<ServerStyleSheets>,
    createGenerateClassName:() => () => string,
    useTheme:() => Theme,
  };
}
