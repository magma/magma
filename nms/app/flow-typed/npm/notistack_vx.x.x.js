// flow-typed signature: ad0c56ce2cda382fa9b8663f831feb7f
// flow-typed version: <<STUB>>/notistack_v0.8.x

declare module 'notistack' {
  import type {ComponentType, ElementConfig, Node} from 'react';

  declare type SnackBarKey = string | number;
  declare type MutualProps = {|
    children?: Node,
    preventDuplicate?: boolean,
    action?: Node | ((key: SnackBarKey) => Node),
    anchorOrigin?: {
      horizontal: 'left' | 'center' | 'right',
      vertical: 'top' | 'bottom',
    },
    autoHideDuration?: number,
    disableWindowBlurListener?: boolean,
    onClose?: () => void,
    onEnter?: () => void,
    onEntered?: () => void,
    onEntering?: () => void,
    onExit?: () => void,
    onExited?: () => void,
    onExiting?: () => void,
    resumeHideDuration?: number,
    TransitionComponent?: Node,
    transitionDuration?: number | {enter?: number, exit?: number},
  |};

  declare type SnackbarProviderProps = {|
    classes?: Object,
    maxSnack?: number,
    iconVariant?: {success?: Node, warning?: Node, error?: Node, info?: Node},
    hideIconVariant?: boolean,
    dense?: boolean,
    ...MutualProps,
  |};

  declare export type Variants =
    | 'default'
    | 'success'
    | 'error'
    | 'info'
    | 'warning';

  declare export type EnqueueSnackbarOptions = {|
    key?: string | number,
    variant?: Variants,
    persist?: boolean,
    children?: SnackBarKey => Node,
    ...MutualProps,
  |};

  declare export type WithSnackbarProps = {|
    enqueueSnackbar: (
      message: string | Node,
      options?: EnqueueSnackbarOptions,
    ) => ?SnackBarKey,
    closeSnackbar: (key?: SnackBarKey) => void,
  |};

  declare export var withSnackbar: <WrappedComponent: ComponentType<*>>(
    Component: WrappedComponent,
  ) => ComponentType<$Diff<ElementConfig<WrappedComponent>, WithSnackbarProps>>;

  declare export var SnackbarProvider: ComponentType<SnackbarProviderProps>;
  declare export var useSnackbar: () => WithSnackbarProps;
}
