// flow-typed signature: 2f5b0f18c5b5b31b01f63163429685c7
// flow-typed version: 5edd39ab2e/@storybook/addon-actions_v3.x.x/flow_>=v0.25.x

declare module '@storybook/addon-actions' {
  declare type Action = (name: string) => (...args: Array<any>) => void;
  declare type DecorateFn = (args: Array<any>) => Array<any>;

  declare module.exports: {
    action: Action,
    decorateAction(args: Array<DecorateFn>): Action;
  };
}
