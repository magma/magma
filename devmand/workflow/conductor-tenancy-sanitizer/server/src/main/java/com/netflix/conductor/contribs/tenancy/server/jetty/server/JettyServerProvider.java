package com.netflix.conductor.contribs.tenancy.server.jetty.server;

import java.util.Optional;
import javax.inject.Provider;

public class JettyServerProvider implements Provider<Optional<JettyServer>> {

    @Override
    public Optional<JettyServer> get() {
        return
                Optional.of(
                        new JettyServer(
                                JettyServerConfiguration.PORT_DEFAULT_VALUE,
                                JettyServerConfiguration.JOIN_DEFAULT_VALUE
                        ));
    }
}
