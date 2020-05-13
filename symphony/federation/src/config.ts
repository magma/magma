import config from 'config'
import { ListenOptions } from 'net'
import { GatewayConfig } from '@apollo/gateway'
import { ApolloServerExpressConfig } from 'apollo-server-express'

export class Config {
    constructor(private _config: Omit<config.IConfig, 'util'> = config) {}

    get gateway(): GatewayConfig | undefined {
        return this.get('gateway')
    }

    get server(): ApolloServerExpressConfig | undefined {
        return this.get('server')
    }

    get listen(): Array<ListenOptions> {
        const opt = this.get<ListenOptions>('listen')
        if (opt !== undefined) {
            return [opt]
        }
        return []
    }

    private get<T>(setting: string): T | undefined {
        if (this._config.has(setting)) {
            return this._config.get(setting)
        }
        return undefined
    }
}

const instance = new Config()
export default instance
