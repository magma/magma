// flow-typed signature: 65fe3d22f0866949d449eb0fd198b273
// flow-typed version: 6f07eebf2a/lodash-es_v4.x.x/flow_>=v0.63.x

declare module "bcryptjs" {
  declare function genSaltSync(rounds?: number): string;
  declare function genSalt(rounds: number): Promise<string>;
  declare function genSalt(): Promise<string>;
  declare function genSalt(callback: (err: Error, salt: string) => void): void;
  declare function genSalt(
    rounds: number,
    callback: (err: Error, salt: string) => void
  ): void;
  declare function hashSync(data: string, salt: string): string;
  declare function hashSync(data: string, rounds: number): string;
  declare function hash(
    data: string,
    saltOrRounds: string | number
  ): Promise<string>;
  declare function hash(
    data: string,
    rounds: number,
    callback: (err: Error, encrypted: string) => void
  ): void;
  declare function hash(
    data: string,
    salt: string,
    callback: (err: Error, encrypted: string) => void
  ): void;
  declare function compareSync(data: string, encrypted: string): boolean;
  declare function compare(data: string, encrypted: string): Promise<boolean>;
  declare function compare(
    data: string,
    encrypted: string,
    callback: (err: Error, same: boolean) => void
  ): void;
  declare function getRounds(encrypted: string): number;
}
