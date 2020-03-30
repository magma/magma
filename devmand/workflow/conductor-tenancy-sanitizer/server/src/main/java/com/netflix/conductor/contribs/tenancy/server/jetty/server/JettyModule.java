package com.netflix.conductor.contribs.tenancy.server.jetty.server;

import com.google.inject.AbstractModule;

public class JettyModule extends AbstractModule {
    @Override
    protected void configure() {
        bind(JettyServerProvider.class);
    }
}
