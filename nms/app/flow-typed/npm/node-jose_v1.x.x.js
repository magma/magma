declare module "node-jose"{

    declare type KeyQuery = {
        kty?: string,
        kid?: string,

    };

    declare type KeyStore = {
        get: (KeyQuery) => ?Key
    };

    declare type Key = {

    }

    declare type JWSResult = {
        header: {},
        payload: Buffer,
        signature: Buffer,
        key: string
    };

    declare interface JWSVerify {
        verify: (data: string) => Promise<JWSResult>
    }

    declare class JWS {
        createVerify(keystore: KeyStore) : JWSVerify
    }

    declare class NodeJose {
        JWS: JWS;
    }

    declare module.exports : NodeJose;
}
