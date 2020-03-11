// flow-typed signature: 904ef6fb904ecaca17975b9bd02e8220
// flow-typed version: <<STUB>>/beaver-logger_v4.0.12

declare module 'beaver-logger' {
  declare type LogLevel = 'debug' | 'info' | 'warn' | 'error';
  declare type LogLevelObj = {
    DEBUG: 'debug',
    INFO: 'info',
    WARN: 'warn',
    ERROR: 'error',
  };
  declare var LOG_LEVEL: LogLevelObj;

  declare type Transport = ({
    url: string,
    method: string,
    headers: Payload,
    json: {
      events: Array<{
        event: {...},
        level: LogLevel,
        payload: {
          timestamp: number,
          data: {...},
          user: {|
            tenant: string,
            email: string,
            isSuperUser: boolean,
            isReadOnlyUser: boolean,
          |},
        },
      }>,
    },
  }) => Promise<void>;

  declare type LoggerOptions = {|
    url: string,
    prefix?: string,
    logLevel?: LogLevel,
    flushInterval?: number,
    transport: Transport,
  |};

  declare type ClientPayload = {...};
  declare type Log = (name: string, payload?: ClientPayload) => LoggerType;
  declare type Track = (payload: ClientPayload) => LoggerType;
  declare type Builder = (Payload) => ClientPayload;
  declare type AddBuilder = (Builder) => LoggerType;

  declare type LoggerType = {|
    debug: Log,
    info: Log,
    warn: Log,
    error: Log,
    track: Track,
    addPayloadBuilder: AddBuilder,
    addMetaBuilder: AddBuilder,
    addTrackingBuilder: AddBuilder,
    addHeaderBuilder: AddBuilder,
  |};

  declare class BeaverLogger {
    Logger: (options: LoggerOptions) => LoggerType;
    LOG_LEVEL: LogLevelObj;
  }

  declare module.exports: BeaverLogger;
}
