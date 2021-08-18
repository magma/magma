declare module '@material-ui/styles/ThemeProvider' {
  import type {ComponentType} from 'react';
  declare module.exports: ComponentType<Object>;
}

declare module '@material-ui/styles' {
  import type {ComponentType, Node} from 'react';
  import type {Theme} from '@material-ui/core';
  declare class ServerStyleSheets {
    collect: Node => Node;
  }

  declare type Style<Props, Classes> = $Shape<{
    [Classes]: {...} | (Props => {...}),
  }>;

  declare type StyleHookFn<_Props, Stl> = (
    props?: _Props,
  ) => $ObjMap<Stl, () => string>;

  declare module.exports: {
    makeStyles: <Props, Stl: Style<Props, string>>(
      Theme => Stl,
    ) => StyleHookFn<Props, Stl>,
    StylesProvider: ComponentType<{generateClassName?: () => string}>,
    ServerStyleSheets: Class<ServerStyleSheets>,
    createGenerateClassName: () => () => string,
    useTheme: () => Theme,
  };
}
