// flow-typed signature: 7cb18569665ce92a9a8156df681f8f2e
// flow-typed version: c6154227d1/shortid_v2.2.x/flow_>=v0.104.x

declare module 'shortid' {
  declare type ShortIdModule = {|
    (): string,
    generate(): string,
    seed(seed: number): ShortIdModule,
    worker(workerId: number): ShortIdModule,
    characters(characters: string): string,
    decode(id: string): {
      version: number,
      worker: number,
      ...
    },
    isValid(id: mixed): boolean,
  |};
  declare module.exports: ShortIdModule;
};
