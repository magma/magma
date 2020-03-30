package com.netflix.conductor.contribs.tenancy.server.jackson;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.inject.AbstractModule;

public class JacksonModule extends AbstractModule {

    @Override
    protected void configure() {
        bind(ObjectMapper.class).toProvider(JsonMapperProvider.class);
    }
}
